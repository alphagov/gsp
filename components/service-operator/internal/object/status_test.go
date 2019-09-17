package object_test

import (
	"github.com/alphagov/gsp/components/service-operator/internal/object"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Status", func() {

	type Object struct {
		metav1.TypeMeta   `json:",inline"`
		metav1.ObjectMeta `json:"metadata,omitempty"`
		object.Status     `json:"status"`
	}

	var o *Object

	BeforeEach(func() {
		o = &Object{
			Status: object.Status{
				State: object.ReadyState,
			},
		}
	})

	It("should allow fetching state", func() {
		state := o.GetState()
		Expect(state).To(Equal(object.ReadyState))
	})

	It("should allow setting state", func() {
		o.SetState(object.ReconcilingState)
		Expect(o.Status.State).To(Equal(object.ReconcilingState))
	})

	It("should allow fetching status", func() {
		status := o.GetStatus()
		Expect(status.State).To(Equal(object.ReadyState))
	})

	It("should allow setting status", func() {
		status := object.Status{
			State: object.ReconcilingState,
		}
		o.SetStatus(status)
		Expect(o.Status.State).To(Equal(object.ReconcilingState))
	})

	It("should return zero state if object is nil", func() {
		var none *object.Status
		Expect(none.GetState()).To(BeZero())
	})

	It("should return zero status if object is nil", func() {
		var none *object.Status
		Expect(none.GetStatus()).To(BeZero())
	})

})
