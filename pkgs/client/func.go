package client

import (
	"fmt"
	"os"

	"github.com/f4tal-err0r/discord_faas/proto"
	"gopkg.in/yaml.v3"
)

func DeployFunc(fp string) error {
	data, err := os.ReadFile(fp + "/dfaas.yaml")
	if err != nil {
		return fmt.Errorf("unable to open dfaas.yaml: %s", err)
	}

	// parse yaml

	var BuildReq proto.BuildFunc

	err = yaml.Unmarshal(data, &BuildReq)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", &BuildReq)

	return nil
}
