package deploy

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/f4tal-err0r/discord_faas/pkgs/config"
	"github.com/f4tal-err0r/discord_faas/pkgs/discord"
	"github.com/f4tal-err0r/discord_faas/pkgs/runner"
	"k8s.io/client-go/kubernetes"
)

type Handler struct {
	dbot    *discord.Client
	cfg     *config.Config
	cs      *kubernetes.Clientset
	storage Storage
}

type Storage interface {
	AddSrcArtifact(ctx context.Context, name string, data io.Reader, size int64) error
	GetSrcPath(ctx context.Context, name string) (string, error)
	GetPresignedUrl(ctx context.Context, bucket string, cmdid string) (*url.URL, error)
	DeleteSrcArtifact(ctx context.Context, name string) error
}

func NewHandler(cfg *config.Config, dbot *discord.Client, cs *kubernetes.Clientset, storage Storage) *Handler {
	return &Handler{
		dbot:    dbot,
		cfg:     cfg,
		cs:      cs,
		storage: storage,
	}
}

func (h *Handler) Builder(cmdid string) error {
	uri, err := h.storage.GetSrcPath(context.Background(), cmdid)
	if err != nil {
		return err
	}
	ropts := runner.RunnerOpts{
		Id:    fmt.Sprintf("build-%s", cmdid),
		Image: "gcr.io/kaniko-project/executor:latest",
		Cmd: []string{
			"/kaniko/executor",
			fmt.Sprintf("--context=%s.tar.gz", uri),
			"--no-push",
			"--build-arg=S3_UPLOAD_URL=\"%s\"",
		},
	}
	r := runner.NewK8sRunner(h.cs)

	uploadUrl, err := h.storage.GetPresignedUrl(context.Background(), "faas-artifacts", cmdid)
	if err != nil {
		return err
	}

	return r.CreateRunner(ropts, uploadUrl.String())
}
