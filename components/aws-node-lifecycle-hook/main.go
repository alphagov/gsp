package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alphagov/gsp/components/aws-node-lifecycle-hook/pkg/awsclient"
	"github.com/alphagov/gsp/components/aws-node-lifecycle-hook/pkg/k8sclient"
	"github.com/alphagov/gsp/components/aws-node-lifecycle-hook/pkg/k8sdrainer"
	"github.com/alphagov/gsp/components/aws-node-lifecycle-hook/pkg/lifecycle"
	"github.com/aws/aws-lambda-go/lambda"
)

// Start configures a lifecycle handler and registers it with the lambda handler.
func Start() error {
	awsClient, err := awsclient.New()
	if err != nil {
		return fmt.Errorf("failed to configure aws client: %s", err)
	}
	clusterName := os.Getenv("CLUSTER_NAME")
	if clusterName == "" {
		return fmt.Errorf("CLUSTER_NAME environment variable is required")
	}
	k8sClient, err := k8sclient.New(clusterName)
	if err != nil {
		return err
	}
	h := lifecycle.Handler{
		AWSClient:        awsClient,
		KubernetesClient: k8sClient,
		Drainer:          k8sdrainer.DefaultDrainer,
	}
	lambda.Start(h.HandleEvent)
	return nil
}

func main() {
	if err := Start(); err != nil {
		log.Fatalf("failed to startup: %s", err)
	}
}
