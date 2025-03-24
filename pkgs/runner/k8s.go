package runner

import (
	"context"
	"fmt"

	spec "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type RunnerPod struct {
	cs  *kubernetes.Clientset
	pod *spec.Pod
}

func NewK8sRunner(cs *kubernetes.Clientset, opts RunnerOpts) *RunnerPod {
	return &RunnerPod{
		cs: cs,
		pod: &spec.Pod{
			Spec: spec.PodSpec{
				Containers: []spec.Container{
					{
						Name:    fmt.Sprintf("faas-%s", opts.Id),
						Image:   opts.Image,
						Command: opts.Cmd,
						VolumeMounts: []spec.VolumeMount{
							{
								Name:      "faas-artifacts",
								MountPath: "/artifacts",
							},
						},
					},
				},
				Volumes: []spec.Volume{
					{
						Name: "faas-artifacts",
						VolumeSource: spec.VolumeSource{
							PersistentVolumeClaim: &spec.PersistentVolumeClaimVolumeSource{
								ClaimName: "faas-artifacts",
							},
						},
					},
				},
			},
		},
	}
}

func (r *RunnerPod) Run() error {
	_, err := r.cs.CoreV1().Pods("").Create(context.Background(), r.pod, v1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (r *RunnerPod) TailLogs() (chan []byte, error) {
	var outch chan []byte

	plr, err := r.cs.CoreV1().Pods("").GetLogs(r.pod.Name, &spec.PodLogOptions{Follow: true}).Stream(context.Background())
	if err != nil {
		return nil, err
	}

	outch = make(chan []byte)
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := plr.Read(buf)
			if err != nil {
				close(outch)
				return
			}
			outch <- buf[:n]
		}
	}()

	return outch, nil
}
