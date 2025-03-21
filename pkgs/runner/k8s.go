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

func NewK8sRunnerSpec(cs *kubernetes.Clientset, img string, hash string, cmd []string) *RunnerPod {
	return &RunnerPod{
		pod: &spec.Pod{
			Spec: spec.PodSpec{
				Containers: []spec.Container{
					{
						Name:    fmt.Sprintf("faas-%s", hash),
						Image:   img,
						Command: cmd,
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

func (r *RunnerPod) TailLogs(cs *kubernetes.Clientset) (chan []byte, error) {
	var outch chan []byte

	plr, err := cs.CoreV1().Pods("").GetLogs(r.pod.Name, &spec.PodLogOptions{Follow: true}).Stream(context.Background())
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
