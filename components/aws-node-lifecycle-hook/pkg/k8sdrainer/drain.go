package k8sdrainer

import (
	"os"
	"time"

	"github.com/alphagov/gsp/components/aws-node-lifecycle-hook/pkg/k8sclient"
	v1 "k8s.io/api/core/v1"
	"k8s.io/kubectl/pkg/drain"
)

const (
	CORDON = true
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fakek8sdrainer/fake_drainer.go . Drainer

type Drainer interface {
	Cordon(k8sclient.Client, *v1.Node) error
	Drain(k8sclient.Client, *v1.Node) error
}

var DefaultDrainer Drainer = &DrainHandler{}

type DrainHandler struct{}

func (d *DrainHandler) Drain(c k8sclient.Client, node *v1.Node) error {
	cfg := &drain.Helper{
		Client:              c,
		Force:               true,
		GracePeriodSeconds:  120,
		IgnoreAllDaemonSets: true,
		Timeout:             time.Minute * 9,
		DeleteLocalData:     true,
		Out:                 os.Stdout,
		ErrOut:              os.Stderr,
	}
	return drain.RunNodeDrain(cfg, node.ObjectMeta.Name)
}

func (d *DrainHandler) Cordon(c k8sclient.Client, node *v1.Node) error {
	cfg := &drain.Helper{
		Client:  c,
		Timeout: time.Minute * 9,
		Out:     os.Stdout,
		ErrOut:  os.Stderr,
	}
	return drain.RunCordonOrUncordon(cfg, node, CORDON)
}
