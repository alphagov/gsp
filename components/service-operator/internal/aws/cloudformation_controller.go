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

	"github.com/alphagov/gsp/components/service-operator/internal"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	goformation "github.com/awslabs/goformation/cloudformation"
)

type CloudFormationTemplate interface {
	Template(string) *goformation.Template
	Parameters() ([]*cloudformation.Parameter, error)
	ResourceType() string
}

// CloudFormationController reconciles an AWS object
type CloudFormationController struct {
	ClusterName string
}

var (
	nonUpdatable = []string{
		cloudformation.StackStatusCreateInProgress,
		cloudformation.StackStatusRollbackInProgress,
		cloudformation.StackStatusDeleteInProgress,
		cloudformation.StackStatusUpdateInProgress,
		cloudformation.StackStatusUpdateCompleteCleanupInProgress,
		cloudformation.StackStatusUpdateRollbackInProgress,
		cloudformation.StackStatusUpdateRollbackCompleteCleanupInProgress,
		cloudformation.StackStatusReviewInProgress,
		cloudformation.StackStatusDeleteComplete,
	}
)

func (r *CloudFormationController) Reconcile(log logr.Logger, ctx context.Context, req ctrl.Request, cloudFormationTemplate CloudFormationTemplate, deleting bool) (internal.Action, StackData, error) {
	stackName := fmt.Sprintf("%s-%s-%s-%s-%s", r.ClusterName, "gsp-service-operator", cloudFormationTemplate.ResourceType(), req.NamespacedName.Namespace, req.NamespacedName.Name)
	// secretName := coalesceString(postgres.Spec.Secret, postgres.Name)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// metadata := ec2metadata.New(sess)
	// awsRegion, err := metadata.Region()
	// if err != nil {
	// 	return ctrl.Result{}, fmt.Errorf("unable to get AWS region from metdata: %s", err)
	// }
	awsRegion := "eu-west-2"
	sess.Config.Region = aws.String(awsRegion)

	svc := cloudformation.New(sess, aws.NewConfig())

	stackData, stackExists := r.getCloudFormationStackStatus(svc, stackName)

	if deleting {
		// The resource needs deleting
		if !stackExists || stackData.Status == cloudformation.StackStatusDeleteComplete {
			return internal.Delete, stackData, nil
		}

		if stackData.Status == cloudformation.StackStatusDeleteInProgress {
			return internal.Retry, stackData, nil
		}

		return internal.Retry, stackData, r.deleteCloudFormationStack(svc, stackName, log)
	}

	yaml, err := cloudFormationTemplate.Template(stackName).YAML()
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

func (r *CloudFormationController) getCloudFormationStackStatus(svc *cloudformation.CloudFormation, stackName string) (StackData, stackExists) {
	describeOutput, err := svc.DescribeStacks(&cloudformation.DescribeStacksInput{StackName: aws.String(stackName)})
	if err != nil {
		return StackData{}, false
	}
	return StackData{
		ID:     *describeOutput.Stacks[0].StackId,
		Name:   stackName,
		Status: *describeOutput.Stacks[0].StackStatus,
		Reason: "NoReasonGiven",

		Outputs: describeOutput.Stacks[0].Outputs,
	}, true
}

func (r *CloudFormationController) createCloudFormationStack(
	yaml []byte,
	cloudFormationTemplate CloudFormationTemplate,
	svc *cloudformation.CloudFormation,
	stackName string,
	log logr.Logger) error {
	log.V(1).Info("creating stack...", "stackName", stackName)

	params, err := cloudFormationTemplate.Parameters()
	if err != nil {
		return fmt.Errorf("error creating parameters: %s", err)
	}

	_, err = svc.CreateStack(&cloudformation.CreateStackInput{
		TemplateBody: aws.String(string(yaml)),
		StackName:    aws.String(stackName),
		Parameters:   params,
	})
	if err != nil {
		return fmt.Errorf("error creating stack: %s", err)
	}

	// TODO: create Secret

	return nil
}

func (r *CloudFormationController) updateCloudFormationStack(
	yaml []byte,
	cloudFormationTemplate CloudFormationTemplate,
	svc *cloudformation.CloudFormation,
	stackName string,
	log logr.Logger) error {

	log.V(1).Info("updating stack...", "stackName", stackName)

	params, err := cloudFormationTemplate.Parameters()
	if err != nil {
		return fmt.Errorf("error creating parameters: %s", err)
	}

	_, err = svc.UpdateStack(&cloudformation.UpdateStackInput{
		TemplateBody: aws.String(string(yaml)),
		StackName:    aws.String(stackName),
		Parameters:   params,
	})
	if err != nil {
		return fmt.Errorf("error updating stack: %s", err)
	}

	return nil
}

func (r *CloudFormationController) deleteCloudFormationStack(svc *cloudformation.CloudFormation, stackName string, log logr.Logger) error {
	log.V(1).Info("deleting stack...", "stackName", stackName)
	_, err := svc.DeleteStack(&cloudformation.DeleteStackInput{StackName: aws.String(stackName)})
	// TODO: delete Secret
	return err
}
