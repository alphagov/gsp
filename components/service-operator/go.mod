module github.com/alphagov/gsp/components/service-operator

go 1.13

require (
	github.com/aws/aws-sdk-go v1.30.4
	github.com/awslabs/goformation/v4 v4.7.0
	github.com/go-logr/logr v0.1.0
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/maxbrunsfeld/counterfeiter/v6 v6.2.2
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.9.0
	github.com/sanathkr/yaml v0.0.0-20170819201035-0056894fa522
	istio.io/istio v0.0.0-20190925083542-b158283f0728
	k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	sigs.k8s.io/controller-runtime v0.2.0-beta.4
)

// avoid jsonpatch v2.0.0 which was yanked and republished and so has
// two different hashes floating around
// https://github.com/gomodules/jsonpatch/issues/21
replace gomodules.xyz/jsonpatch/v2 => gomodules.xyz/jsonpatch/v2 v2.0.1
