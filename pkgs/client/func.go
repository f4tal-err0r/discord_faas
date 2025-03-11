package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	pb "github.com/f4tal-err0r/discord_faas/proto"
	proto "google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

func DeployFunc(fp string) error {
	data, err := os.ReadFile(fp + "/dfaas.yaml")
	if err != nil {
		return fmt.Errorf("unable to open dfaas.yaml: %s", err)
	}

	// parse yaml

	var BuildReq pb.BuildFunc

	err = yaml.Unmarshal(data, &BuildReq)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", &BuildReq)

	return nil
}

func TestCommand(cmd *pb.DiscordResp, conn string) (*pb.DiscordResp, error) {
	var response pb.DiscordResp
	client := &http.Client{Timeout: 10 * time.Second}

	marshaledCmd, err := proto.Marshal(cmd)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", conn, bytes.NewBuffer(marshaledCmd))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = proto.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
