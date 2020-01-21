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
	var awsClient *sdkfakes.FakeClient
	var stack *cloudformationfakes.FakeStack
	var ctx context.Context
	var ctxCancel func()
	var ctxTimeout = time.Millisecond * 100

	BeforeEach(func() {
		ctx, ctxCancel = context.WithTimeout(context.Background(), ctxTimeout)
		awsClient = &sdkfakes.FakeClient{}
		client = &cloudformation.Client{
			Client:          awsClient,
			PollingInterval: time.Millisecond * 25,
		}
		stack = &cloudformationfakes.FakeStack{}
	})

	AfterEach(func() {
		ctxCancel()
	})

	Context("state transitions", func() {

		var state string

		JustBeforeEach(func() {
			// always return a dummy template and no error
			stack.GetStackTemplateReturns(&cloudformation.Template{}, nil)
			// return state as set in subtests
			awsClient.DescribeStacksWithContextReturns(&cloudformation.DescribeStacksOutput{
				Stacks: []*cloudformation.State{
					{
						StackStatus: aws.String(state),
					},
				},
			}, nil)
		})

		Context("when state is unprovisioned", func() {
			// There is no previous state, the stack does not exist
			BeforeEach(func() {
				awsClient.DescribeStacksWithContextReturnsOnCall(0, nil, sdkfakes.ResourceNotFoundException)
				state = cloudformation.CreateInProgress // used after initial call
			})
			It("should create on Apply()", func() {
				_, _ = client.Apply(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(1))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(0))
			})
			It("should be a no-op on Destroy()", func() {
				_ = client.Destroy(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(0))
			})
		})

		Context("when state is CREATE_IN_PROGRESS", func() {
			// Ongoing creation of one or more stacks.
			BeforeEach(func() {
				state = cloudformation.CreateInProgress
			})
			It("should be no-op on Apply()", func() {
				_, _ = client.Apply(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(0))
			})
			It("should delete on Destroy()", func() {
				_ = client.Destroy(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(1))
			})
		})

		Context("when state is REVIEW_IN_PROGRESS", func() {
			// Ongoing creation of one or more stacks with an expected
			// StackId but without any templates or resources.
			BeforeEach(func() {
				state = cloudformation.ReviewInProgress
			})
			It("should be no-op on Apply()", func() {
				_, _ = client.Apply(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(0))
			})
			It("should delete on Destroy()", func() {
				_ = client.Destroy(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(1))
			})
		})

		Context("when state is CREATE_COMPLETE", func() {
			// Successful creation of one or more stacks.
			BeforeEach(func() {
				state = cloudformation.CreateComplete
			})
			It("should update on Apply()", func() {
				_, _ = client.Apply(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(1))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(0))
			})
			It("should delete on Destroy()", func() {
				_ = client.Destroy(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(1))
			})
		})

		Context("when state is CREATE_FAILED", func() {
			// Unsuccessful creation of one or more stacks.
			// (rollback will be initiated)
			BeforeEach(func() {
				state = cloudformation.CreateFailed
			})
			It("should be no-op on Apply()", func() {
				_, _ = client.Apply(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(0))
			})
			It("should delete on Destroy()", func() {
				_ = client.Destroy(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(1))
			})
		})

		Context("when state is ROLLBACK_IN_PROGRESS", func() {
			// Ongoing removal of one or more stacks after a failed stack
			// creation or after an explicitly cancelled stack creation.
			BeforeEach(func() {
				state = cloudformation.RollbackInProgress
			})
			It("should be no-op on Apply()", func() {
				_, _ = client.Apply(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(0))
			})
			It("should delete on Destroy()", func() {
				_ = client.Destroy(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(1))
			})
		})

		Context("when state is UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS", func() {
			// Ongoing removal of old resources for one or more stacks
			// after a successful stack update. For stack updates that
			// require resources to be replaced, AWS CloudFormation creates
			// the new resources first and then deletes the old resources
			// to help reduce any interruptions with your stack. In this
			// state, the stack has been updated and is usable, but AWS
			// CloudFormation is still deleting the old resources.
			BeforeEach(func() {
				state = cloudformation.UpdateRollbackCompleteCleanupInProgress
			})
			It("should be no-op on Apply()", func() {
				_, _ = client.Apply(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(0))
			})
			It("should delete on Destroy()", func() {
				_ = client.Destroy(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(1))
			})
		})

		Context("when state is UPDATE_IN_PROGRESS", func() {
			// Ongoing update of one or more stacks.
			BeforeEach(func() {
				state = cloudformation.UpdateInProgress
			})
			It("should be no-op on Apply()", func() {
				_, _ = client.Apply(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(0))
			})
			It("should delete on Destroy()", func() {
				_ = client.Destroy(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(1))
			})
		})

		Context("when state is UPDATE_ROLLBACK_IN_PROGRESS", func() {
			// Ongoing return of one or more stacks to the previous working
			// state after failed stack update.
			BeforeEach(func() {
				state = cloudformation.UpdateRollbackInProgress
			})
			It("should be no-op on Apply()", func() {
				_, _ = client.Apply(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(0))
			})
			It("should delete on Destroy()", func() {
				_ = client.Destroy(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(1))
			})
		})

		Context("when state is UPDATE_COMPLETE", func() {
			// Successful update of one or more stacks.
			BeforeEach(func() {
				state = cloudformation.UpdateComplete
			})
			It("should update on Apply()", func() {
				_, _ = client.Apply(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(1))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(0))
			})
			It("should delete on Destroy()", func() {
				_ = client.Destroy(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(1))
			})
		})

		Context("when state is UPDATE_ROLLBACK_COMPLETE", func() {
			// Successful update of one or more stacks.
			BeforeEach(func() {
				state = cloudformation.UpdateRollbackComplete
			})
			It("should update on Apply()", func() {
				_, _ = client.Apply(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(1))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(0))
			})
			It("should delete on Destroy()", func() {
				_ = client.Destroy(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(1))
			})
		})

		Context("when state is DELETE_IN_PROGRESS", func() {
			// Ongoing removal of one or more stacks.
			BeforeEach(func() {
				state = cloudformation.DeleteInProgress
			})
			It("should be no-op on Apply()", func() {
				_, _ = client.Apply(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(0))
			})
			It("should be no-op on Destroy()", func() {
				_ = client.Destroy(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(0))
			})
		})

		Context("when state is DELETE_FAILED", func() {
			// Unsuccessful deletion of one or more stacks. Because
			// the delete failed, you might have some resources
			// that are still running; however, you cannot work
			// with or update the stack. Delete the stack again or
			// view the stack events to see any associated error
			// messages.
			BeforeEach(func() {
				state = cloudformation.DeleteFailed
			})
			It("should trigger delete on Apply()", func() {
				Skip("we should delete and then recreate on apply, but this is not implemented yet")
				_, _ = client.Apply(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(1))
			})
			It("should retry delete on Destroy()", func() {
				_ = client.Destroy(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(1))
			})
		})

		Context("when state is DELETE_COMPLETE", func() {
			// Successful deletion of one or more stacks. Deleted
			// stacks are retained and viewable for 90 days.
			BeforeEach(func() {
				state = cloudformation.DeleteComplete
			})
			It("should be no-op on Apply()", func() {
				_, _ = client.Apply(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(0))
			})
			It("should be no-op on Destroy()", func() {
				_ = client.Destroy(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(0))
			})
		})

		Context("when state is ROLLBACK_COMPLETE", func() {
			// This state is a failed creation, and the only valid
			// transition is delete
			BeforeEach(func() {
				state = cloudformation.RollbackComplete
			})
			It("should delete on Apply()", func() {
				Skip("we should delete and then recreate on apply, but this is not implemented yet")
				_, _ = client.Apply(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(1))
			})
			It("should delete on Destroy()", func() {
				_ = client.Destroy(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(1))
			})
		})

		Context("when state is UPDATE_ROLLBACK_FAILED", func() {
			// Unsuccessful return of one or more stacks to a previous
			// working state after a failed stack update. When in this
			// state, you can delete the stack or continue rollback. You
			// might need to fix errors before your stack can return to a
			// working state. Or, you can contact customer support to
			// restore the stack to a usable state.
			BeforeEach(func() {
				state = cloudformation.UpdateRollbackFailed
			})
			It("should continue-rollback on Apply()", func() {
				// TODO: what is "continue rollback" and how do we do that?
				// Currently this is a no-op
				_, _ = client.Apply(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(0))
			})
			It("should delete on Destroy()", func() {
				_ = client.Destroy(ctx, stack)
				Expect(awsClient.DescribeStackEventsWithContextCallCount()).To(BeNumerically(">", 0))
				Expect(awsClient.CreateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.UpdateStackWithContextCallCount()).To(Equal(0))
				Expect(awsClient.DeleteStackWithContextCallCount()).To(Equal(1))
			})
		})

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
			awsClient.DescribeStacksWithContextReturns(&cloudformation.DescribeStacksOutput{
				Stacks: []*cloudformation.State{state},
			}, nil)
			awsClient.GetSecretValueWithContextReturns(&secretsmanager.GetSecretValueOutput{
				SecretString: aws.String(`{"username": "admin", "password": "abc123"}`),
			}, nil)
			awsClient.UpdateStackWithContextReturns(nil, sdkfakes.NoUpdateRequiredException)
		})

		It("should attempt to resolve template outputs that are references to secrets manager", func() {
			outputs, err := client.Outputs(ctx, stack)
			Expect(err).ToNot(HaveOccurred())
			Expect(outputs).To(HaveLen(1))
			Expect(outputs).To(HaveKeyWithValue("a-secret-output", "abc123"))
		})

	})

})
