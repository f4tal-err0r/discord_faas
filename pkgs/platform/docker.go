package platform

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	docker "github.com/docker/docker/client"
)

type Docker struct {
	*docker.Client
}

func NewDockerClient() (*Docker, error) {
	opts := []docker.Opt{
		docker.FromEnv,
		docker.WithAPIVersionNegotiation(),
	}

	client, err := docker.NewClientWithOpts(opts...)
	if err != nil {
		return nil, err
	}

	return &Docker{client}, nil
}

func (d *Docker) BuildImage(dockerfile *os.File, img *Image) error {
	_, err := d.ImageBuild(context.Background(), dockerfile, types.ImageBuildOptions{
		Tags: []string{img.Name},
		Labels: map[string]string{
			"guildid":   img.Meta.GuildID,
			"ownerid":   img.Meta.OwnerID,
			"userid":    img.Meta.UserID,
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// func (d *Docker) Exec(hash string) error {
// 	log, err := d.ContainerCreate(context.Background(), &types.ContainerCreateConfig{
// 		Config: &types.ContainerConfig{
// 			Image: hash,
// 		},
// 	})

// 	if err != nil {
// 		return err
// 	}

// 	return d.ContainerExecStart(context.Background(), log.ID, types.ExecStartCheck{
// 		Detach: true,
// 	})
// }

func (d *Docker) ListImages() ([]*Image, error) {
	var images []*Image
	dockerimgs, err := d.ImageList(context.Background(), image.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	for _, v := range dockerimgs {
		images = append(images, &Image{
			Name:      v.RepoTags[0],
			Runtime:   v.Labels["runtime"],
			Hash:      v.ID,
			Timestamp: time.Unix(v.Created, 0),
			Meta: &Labels{
				GuildID: v.Labels["guildid"],
				OwnerID: v.Labels["ownerid"],
				UserID:  v.Labels["userid"],
			},
		})
	}

	return images, nil
}

func (d *Docker) RemoveImage(hash string) ([]image.DeleteResponse, error) {
	return d.ImageRemove(context.Background(), hash, image.RemoveOptions{
		Force:         true,
		PruneChildren: true,
	})
}
