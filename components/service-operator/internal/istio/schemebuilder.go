package istio

import (
	"k8s.io/apimachinery/pkg/runtime"
	k8sschema "k8s.io/apimachinery/pkg/runtime/schema"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	istioschemas "istio.io/istio/pkg/config/schemas"
	istiocrd "istio.io/istio/pilot/pkg/config/kube/crd"
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
