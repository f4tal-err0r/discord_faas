package client

import (
	"fmt"
	"os"

	"github.com/f4tal-err0r/discord_faas/proto"
	"gopkg.in/yaml.v3"
)

func DeployFunc(fp string) error {
<<<<<<< HEAD
=======

>>>>>>> fa4cfec (wip, working copy of app in k8s; working yaml marshal)
	data, err := os.ReadFile(fp)
	if err != nil {
		return err
	}
<<<<<<< HEAD
=======

	// parse yaml

>>>>>>> fa4cfec (wip, working copy of app in k8s; working yaml marshal)
	var BuildReq proto.BuildFunc

	err = yaml.Unmarshal(data, &BuildReq)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", &BuildReq)

	return nil
}
