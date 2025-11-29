package kubernetes

import (
	"fmt"
	"net/url"

	"github.com/f4tal-err0r/discord_faas/pkgs/platform"
	v1 "k8s.io/api/core/v1"
)

type Runner struct {
	cs      *Client
	podspec *v1.PodSpec
}

func (c *Client) NewRunner(config platform.RunnerConfig) (*Runner, error) {
	runner := &Runner{
		cs: c,
	}

	runner.podspec = &v1.PodSpec{
		Containers: []v1.Container{
			{
				Name:       config.Name,
				Image:      config.Image,
				Command:    config.Command,
				WorkingDir: config.StoragePath,
			},
		},
		RestartPolicy: v1.RestartPolicyNever,
	}

	sp, err := url.Parse(config.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("invalid storage path: %v", err)
	}

	switch sp.Scheme {
	case "local":
		runner.podspec = appendLocalVolume(runner.podspec, "/app/data/")
	case "s3":
		//TODO
	default:
		return nil, fmt.Errorf("unsupported storage path scheme: %s", sp.Scheme)
	}

	return runner, nil
}
