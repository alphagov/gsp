package istio

import (
	istiocrd "istio.io/istio/pilot/pkg/config/kube/crd"
	istioschemas "istio.io/istio/pkg/config/schemas"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sschema "k8s.io/apimachinery/pkg/runtime/schema"
)

func AddToScheme(scheme *runtime.Scheme) error {
	istioSchemeBuilder := runtime.NewSchemeBuilder(
		func(scheme *runtime.Scheme) error {
			gv := k8sschema.GroupVersion{Group: "networking.istio.io", Version: "v1alpha3"}
			st := istiocrd.KnownTypes[istioschemas.ServiceEntry.Type]
			scheme.AddKnownTypes(gv, st.Object, st.Collection)
			meta_v1.AddToGroupVersion(scheme, gv)
			return nil
		})
	return istioSchemeBuilder.AddToScheme(scheme)
}
