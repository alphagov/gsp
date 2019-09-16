package cloudformation_test

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/cloudformation/cloudformationfakes"
	"github.com/alphagov/gsp/components/service-operator/internal/aws/sdk/sdkfakes"
)

var _ = Describe("Client", func() {

	var client *cloudformation.Client
	var sdk *sdkfakes.FakeClient
	var stack *cloudformationfakes.FakeStack
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
		sdk = &sdkfakes.FakeClient{}
		client = &cloudformation.Client{
			Client:          sdk,
			PollingInterval: time.Millisecond * 100,
		}
		stack = &cloudformationfakes.FakeStack{}
	})

	Context("when outputs contain secret manager references", func() {

		BeforeEach(func() {
			state := &cloudformation.State{
				StackId:           aws.String(fmt.Sprintf("stack-%d", rand.Intn(10000))),
				StackName:         aws.String("test-stack"),
				StackStatus:       aws.String(cloudformation.CreateComplete),
				StackStatusReason: aws.String("faked"),
				Outputs: []*cloudformation.Output{
					{
						OutputKey:   aws.String("a-secret-output"),
						OutputValue: aws.String(`{{resolve:secretsmanager:arn:xxxxxxx:SecretString:password}}`),
					},
				},
			}
			sdk.DescribeStacksWithContextReturns(&cloudformation.DescribeStacksOutput{
				Stacks: []*cloudformation.State{state},
			}, nil)
			sdk.GetSecretValueWithContextReturns(&secretsmanager.GetSecretValueOutput{
				SecretString: aws.String(`{"username": "admin", "password": "abc123"}`),
			}, nil)
			sdk.UpdateStackWithContextReturns(nil, sdkfakes.NoUpdateRequiredException)
		})

		It("should attempt to resolve template outputs that are references to secrets manager", func() {
			outputs, err := client.Outputs(ctx, stack)
			Expect(err).ToNot(HaveOccurred())
			Expect(outputs).To(HaveLen(1))
			Expect(outputs).To(HaveKeyWithValue("a-secret-output", "abc123"))
		})

	})

})
