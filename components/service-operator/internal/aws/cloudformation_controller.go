/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/alphagov/gsp/components/service-operator/internal"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awscloudformation "github.com/aws/aws-sdk-go/service/cloudformation"
	goformation "github.com/awslabs/goformation/cloudformation"
	"github.com/awslabs/goformation/cloudformation/resources"
)

type CloudFormationTemplate interface {
	Template(string, []resources.Tag) *goformation.Template
	CreateParameters() ([]*awscloudformation.Parameter, error)
	UpdateParameters() ([]*awscloudformation.Parameter, error)
	ResourceType() string
}

type CloudFormationReconciler interface {
	Reconcile(context.Context, logr.Logger, ctrl.Request, CloudFormationTemplate, bool) (internal.Action, StackData, error)
}

type CloudFormationController struct {
	ClusterName string
}

var (
	nonUpdatable = []string{
		awscloudformation.StackStatusCreateInProgress,
		awscloudformation.StackStatusRollbackInProgress,
		awscloudformation.StackStatusDeleteInProgress,
		awscloudformation.StackStatusUpdateInProgress,
		awscloudformation.StackStatusUpdateCompleteCleanupInProgress,
		awscloudformation.StackStatusUpdateRollbackInProgress,
		awscloudformation.StackStatusUpdateRollbackCompleteCleanupInProgress,
		awscloudformation.StackStatusReviewInProgress,
		awscloudformation.StackStatusDeleteComplete,
	}
)

func (r *CloudFormationController) Reconcile(ctx context.Context, log logr.Logger, req ctrl.Request, cloudFormationTemplate CloudFormationTemplate, deleting bool) (internal.Action, StackData, error) {
	stackName := fmt.Sprintf("%s-%s-%s-%s", r.ClusterName, cloudFormationTemplate.ResourceType(), req.Namespace, req.Name)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// FIXME: We could/should obtain this automatically... At the moment, we have no plans to run this outside of london.
	awsRegion := "eu-west-2"
	sess.Config.Region = aws.String(awsRegion)

	svc := awscloudformation.New(sess, aws.NewConfig())

	stackData, stackExists := r.getCloudFormationStackStatus(svc, stackName, log)

	if deleting {
		// The resource needs deleting
		if !stackExists || stackData.Status == awscloudformation.StackStatusDeleteComplete {
			return internal.Delete, stackData, nil
		}

		if stackData.Status == awscloudformation.StackStatusDeleteInProgress {
			return internal.Retry, stackData, nil
		}

		return internal.Retry, stackData, r.deleteCloudFormationStack(svc, stackName, log)
	}

	yaml, err := cloudFormationTemplate.Template(stackName, DefineTags(r.ClusterName, req.Name, req.Namespace, cloudFormationTemplate.ResourceType())).YAML()
	if err != nil {
		return internal.Retry, stackData, fmt.Errorf("error serialising template: %s", err)
	}

	if !stackExists { // create
		return internal.Create, stackData, r.createCloudFormationStack(yaml, cloudFormationTemplate, svc, stackName, log)
	} else if !internal.IsInList(stackData.Status, nonUpdatable...) { // update
		return internal.Update, stackData, r.updateCloudFormationStack(yaml, cloudFormationTemplate, svc, stackName, log)
	}

	return internal.Retry, stackData, nil
}

type stackExists bool

func (r *CloudFormationController) getCloudFormationStackStatus(svc *awscloudformation.CloudFormation, stackName string, log logr.Logger) (StackData, stackExists) {
	data := StackData{}
	describeOutput, err := svc.DescribeStacks(&awscloudformation.DescribeStacksInput{StackName: aws.String(stackName)})
	if err != nil {
		return data, false
	}
	data.ID = *describeOutput.Stacks[0].StackId
	data.Name = stackName
	data.Status = *describeOutput.Stacks[0].StackStatus
	data.Reason = "NoReasonGiven"
	data.Outputs = describeOutput.Stacks[0].Outputs

	eventsOutput, err := svc.DescribeStackEvents(&awscloudformation.DescribeStackEventsInput{StackName: aws.String(stackName)})
	if err != nil {
		log.Error(err, "unable to retreive stackEvents")
		return data, true
	}
	data.Events = eventsOutput.StackEvents

	return data, true
}

func (r *CloudFormationController) createCloudFormationStack(
	yaml []byte,
	cloudFormationTemplate CloudFormationTemplate,
	svc *awscloudformation.CloudFormation,
	stackName string,
	log logr.Logger) error {
	log.V(1).Info("creating stack...", "stackName", stackName)

	params, err := cloudFormationTemplate.CreateParameters()
	if err != nil {
		return fmt.Errorf("error creating parameters: %s", err)
	}

	_, err = svc.CreateStack(&awscloudformation.CreateStackInput{
		TemplateBody: aws.String(string(yaml)),
		StackName:    aws.String(stackName),
		Parameters:   params,
	})
	if err != nil {
		return fmt.Errorf("error creating stack: %s", err)
	}

	return nil
}

func (r *CloudFormationController) updateCloudFormationStack(
	yaml []byte,
	cloudFormationTemplate CloudFormationTemplate,
	svc *awscloudformation.CloudFormation,
	stackName string,
	log logr.Logger) error {
	log.V(1).Info("updating stack...", "stackName", stackName)

	params, err := cloudFormationTemplate.UpdateParameters()
	if err != nil {
		return fmt.Errorf("error creating parameters: %s", err)
	}

	_, err = svc.UpdateStack(&awscloudformation.UpdateStackInput{
		TemplateBody: aws.String(string(yaml)),
		StackName:    aws.String(stackName),
		Parameters:   params,
	})
	// FIXME: We want to just ignore it if there are no changes to make but AWS
	// don't strongly type errors so we use string comparison.
	if err != nil && !strings.Contains(err.Error(), "No updates are to be performed") {
		return fmt.Errorf("error updating stack: %s", err)
	}

	return nil
}

func (r *CloudFormationController) deleteCloudFormationStack(svc *awscloudformation.CloudFormation, stackName string, log logr.Logger) error {
	log.V(1).Info("deleting stack...", "stackName", stackName)
	_, err := svc.DeleteStack(&awscloudformation.DeleteStackInput{StackName: aws.String(stackName)})
	return err
}
