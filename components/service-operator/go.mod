module github.com/alphagov/gsp/components/service-operator

go 1.12

require (
	github.com/aws/aws-sdk-go v1.22.2
	github.com/awslabs/goformation v0.0.0-20190320125420-ac0a17860cf1
	github.com/go-logr/logr v0.1.0
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/golang/mock v1.3.1
	github.com/maxbrunsfeld/counterfeiter/v6 v6.2.2 // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/sanathkr/yaml v0.0.0-20170819201035-0056894fa522
	go.uber.org/zap v1.10.0 // indirect
	k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	sigs.k8s.io/controller-runtime v0.2.0-beta.4
	sigs.k8s.io/controller-tools v0.2.0-beta.4 // indirect
)
