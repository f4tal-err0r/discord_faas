package kubernetes

import (
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	clientset *kubernetes.Clientset
	namespace string
}

func NewKubernetesPlatform() (*Client, error) {
	var client Client

	ns := os.Getenv("POD_NAMESPACE")
	if ns == "" {
		return nil, fmt.Errorf("POD_NAMESPACE environment variable not set")
	}
	client.namespace = ns

	// Creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("error creating in-cluster config: %v", err)
	}
	// Create the Kubernetes client
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating Kubernetes client: %v", err)
	}

	client.clientset = cs
	return &client, nil
}

func appendLocalVolume(podspec *v1.PodSpec, path string) *v1.PodSpec {
	volumeName := "artifacts-volume"
	podspec.Volumes = append(podspec.Volumes, v1.Volume{
		Name: volumeName,
		VolumeSource: v1.VolumeSource{
			PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
				ClaimName: "discord-faas-artifacts",
			},
		},
	})

	for i := range podspec.Containers {
		podspec.Containers[i].VolumeMounts = append(podspec.Containers[i].VolumeMounts, v1.VolumeMount{
			Name:      volumeName,
			MountPath: path,
		})
	}

	return podspec
}
