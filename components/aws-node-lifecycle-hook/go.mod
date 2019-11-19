module github.com/alphagov/gsp/components/aws-node-lifecycle-hook

go 1.12

require (
	github.com/aws/aws-lambda-go v1.13.2
	github.com/aws/aws-sdk-go v1.25.29
	github.com/gofrs/flock v0.7.1 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/kubernetes-sigs/aws-iam-authenticator v0.4.0
	github.com/maxbrunsfeld/counterfeiter/v6 v6.2.2
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.7.0
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/crypto v0.0.0-20190820162420-60c769a6c586 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	golang.org/x/tools v0.0.0-20190920225731-5eefd052ad72 // indirect
	k8s.io/api v0.0.0-20191107030003-665c8a257c1a
	k8s.io/apimachinery v0.0.0-20191107105744-2c7f8d2b0fd8
	k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
	k8s.io/kubectl v0.0.0-20190918164019-21692a0861df
	k8s.io/utils v0.0.0-20191030222137-2b95a09bc58d // indirect
)

// fix broken upstream
// https://github.com/dominikh/go-tools/issues/658
replace honnef.co/go/tools => github.com/dominikh/go-tools v0.0.0-20190102054323-c2f93a96b099
