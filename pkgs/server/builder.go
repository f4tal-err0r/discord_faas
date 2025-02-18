package server

import (
	"fmt"
	"os"
	"path/filepath"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createKanikoJob(buildImage BuildImage) *batchv1.Job {
	registryUrl := fmt.Sprintf("%s.%s.internal", os.Getenv("POD_NAME"), os.Getenv("POD_NAMESPACE"))

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kaniko-build-job",
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "kaniko",
							Image: "gcr.io/kaniko-project/executor:latest",
							Args: []string{
								"--dockerfile=Dockerfile",
								"--context=tar:///workspace/" + filepath.Join(buildImage.GuildID, buildImage.Name) + ".tar.gz",
								"--destination=" + filepath.Join(registryUrl, buildImage.GuildID, buildImage.Name+":"+buildImage.Version),
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "faas-workspace",
									MountPath: "/workspace",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "faas-workspace",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "faas-workspace-" + os.Getenv("POD_NAME"),
									ReadOnly:  true,
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
			BackoffLimit: int32Ptr(3),
		},
	}
}

func int32Ptr(i int32) *int32 {
	return &i
}
