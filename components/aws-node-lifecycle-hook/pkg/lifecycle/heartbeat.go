package lifecycle

import (
	"context"
	"log"
	"time"

	"github.com/alphagov/gsp/components/aws-node-lifecycle-hook/pkg/awsclient"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func heartbeat(ctx context.Context, asgClient awsclient.Client, asgEvent ASGLifecycleEventDetail, interval time.Duration) {
	log.Printf("starting heartbeat every %v for %s", interval, asgEvent.EC2InstanceId)
	if interval == 0 {
		interval = time.Second * 60
	}
	for {
		select {
		case <-ctx.Done():
			log.Printf("stopping heartbeat for %s", asgEvent.EC2InstanceId)
			return
		case <-time.After(interval):
			_, err := asgClient.RecordLifecycleActionHeartbeatWithContext(ctx, &autoscaling.RecordLifecycleActionHeartbeatInput{
				AutoScalingGroupName: &asgEvent.AutoScalingGroupName,
				InstanceId:           &asgEvent.EC2InstanceId,
				LifecycleHookName:    &asgEvent.LifecycleHookName,
				LifecycleActionToken: &asgEvent.LifecycleActionToken,
			})
			if err != nil {
				log.Printf("heartbeat failed for %s: %s", asgEvent.EC2InstanceId, err)
			}
		}
	}
}
