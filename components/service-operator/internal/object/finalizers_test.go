package object_test

import (
	"github.com/alphagov/gsp/components/service-operator/internal/object"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Finalizers", func() {

	type Object struct {
		metav1.TypeMeta   `json:",inline"`
		metav1.ObjectMeta `json:"metadata,omitempty"`
	}

	var o *Object
	var f string

	BeforeEach(func() {
		o = &Object{}
		f = "my-finalizer"
	})

	It("should add/remove finalizers", func() {
		Expect(o.GetFinalizers()).To(HaveLen(0))
		Expect(object.HasFinalizer(o, f)).To(BeFalse())
		object.SetFinalizer(o, f)
		Expect(o.GetFinalizers()).To(ContainElement("my-finalizer"))
		Expect(object.HasFinalizer(o, f)).To(BeTrue())
		object.RemoveFinalizer(o, "my-finalizer")
		Expect(o.GetFinalizers()).ToNot(ContainElement("my-finalizer"))
		Expect(object.HasFinalizer(o, f)).To(BeFalse())
	})

	It("should not duplicate finalizers", func() {
		object.SetFinalizer(o, f)
		object.SetFinalizer(o, f)
		object.SetFinalizer(o, f)
		Expect(o.GetFinalizers()).To(HaveLen(1))
	})

})
