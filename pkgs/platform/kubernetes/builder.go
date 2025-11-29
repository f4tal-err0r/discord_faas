package kubernetes

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/f4tal-err0r/discord_faas/pkgs/platform"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Builder struct {
	cs  *Client
	job batchv1.Job
}

func (c *Client) NewBuilder(config platform.RunnerConfig) (*Builder, error) {
	builder := &Builder{
		cs: c,
	}

	ps, err := c.NewRunner(config)
	if err != nil {
		return nil, err
	}

	builder.job.Spec = batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			Spec: *ps.podspec,
		},
	}

	builder.job.ObjectMeta.Name = config.Name + "-" + createHash()
	builder.job.ObjectMeta.Labels = config.Labels

	return builder, nil
}

func (b *Builder) Start(ctx *context.Context) error {
	job, err := b.cs.clientset.BatchV1().Jobs(b.cs.namespace).Create(
		*ctx,
		&b.job,
		v1.CreateOptions{},
	)
	if err != nil {
		return err
	}

	b.job = *job

	for {
		j, err := b.cs.clientset.BatchV1().Jobs(b.cs.namespace).Get(
			*ctx,
			b.job.ObjectMeta.Name,
			v1.GetOptions{},
		)
		if err != nil {
			return err
		}

		// check if job started running
		if j.Status.Active > 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}

func createHash() string {
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
