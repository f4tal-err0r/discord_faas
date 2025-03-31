package runner

import (
	"context"
	"fmt"
	"os"

	batchv1 "k8s.io/api/batch/v1"
	spec "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type K8sRunners struct {
	cs   *kubernetes.Clientset
	spec *batchv1.Job
}

func NewK8sRunner(cs *kubernetes.Clientset) *K8sRunners {
	spec := &batchv1.Job{
		Spec: batchv1.JobSpec{
			Template: spec.PodTemplateSpec{
				Spec: spec.PodSpec{
					InitContainers: []spec.Container{
						{
							Name: "discord-faas",
							Env: []spec.EnvVar{
								{
									Name: "AWS_ACCESS_KEY_ID",
									ValueFrom: &spec.EnvVarSource{
										SecretKeyRef: &spec.SecretKeySelector{
											LocalObjectReference: spec.LocalObjectReference{
												Name: "faas-minio-root",
											},
											Key: "MINIO_ROOT_USER",
										},
									},
								},
								{
									Name: "AWS_SECRET_ACCESS_KEY",
									ValueFrom: &spec.EnvVarSource{
										SecretKeyRef: &spec.SecretKeySelector{
											LocalObjectReference: spec.LocalObjectReference{
												Name: "faas-minio-root",
											},
											Key: "MINIO_ROOT_PASSWORD",
										},
									},
								},
								{
									Name:  "S3_ENDPOINT",
									Value: "http://discord-faas:9000",
								},
								{
									Name:  "S3_FORCE_PATH_STYLE",
									Value: "true",
								},
							},
							VolumeMounts: []spec.VolumeMount{
								{
									Name:      "func",
									MountPath: "/artifacts",
								},
							},
						},
					},
					Containers: []spec.Container{
						{
							Name:  "exporter",
							Image: "curlimages/curl:latest",
							VolumeMounts: []spec.VolumeMount{
								{
									Name:      "func",
									MountPath: "/artifacts",
								},
							},
						},
					},
					Volumes: []spec.Volume{
						{
							Name: "func",
							VolumeSource: spec.VolumeSource{
								EmptyDir: &spec.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}

	rp := &K8sRunners{
		spec: spec,
		cs:   cs,
	}
	return rp
}

func (r *K8sRunners) CreateRunner(opts RunnerOpts, uploadUrl string) error {
	runner := r.spec.DeepCopy()

	runner.ObjectMeta.Name = fmt.Sprintf("dfaas-%s", opts.Id)
	runner.ObjectMeta.Namespace = os.Getenv("POD_NAMESPACE")

	runner.Spec.Template.ObjectMeta.Name = fmt.Sprintf("dfaas-%s", opts.Id)
	runner.Spec.Template.Spec.InitContainers[0].Image = opts.Image
	runner.Spec.Template.Spec.InitContainers[0].Command = opts.Cmd
	runner.Spec.Template.Spec.RestartPolicy = "Never"
	runner.Spec.Template.Spec.Containers[0].Env = append(runner.Spec.Template.Spec.Containers[0].Env,
		spec.EnvVar{
			Name:  "S3_UPLOAD_URL",
			Value: uploadUrl,
		},
		spec.EnvVar{
			Name:  "S3_FORCE_PATH_STYLE",
			Value: "true",
		},
	)
	runner.Spec.Template.Spec.Containers[0].Command = []string{
		"curl",
		"-X",
		"PUT",
		"-H",
		"Content-Type: application/octet-stream",
		"-T",
		"/artifacts/func.x",
		uploadUrl,
	}

	_, err := r.cs.BatchV1().Jobs(os.Getenv("POD_NAMESPACE")).Create(context.Background(), runner, v1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
