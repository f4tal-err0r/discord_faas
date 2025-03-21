package client

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	pb "github.com/f4tal-err0r/discord_faas/proto"
	proto "google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

func DeployFunc(context *pb.ContextResp, fp string) error {

	metadata, err := marshalFaasYaml(fp)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	metapart, err := writer.CreateFormField("metadata")
	if err != nil {
		return err
	}
	_, err = metapart.Write(metadata)
	if err != nil {
		return err
	}

	functar, err := createTar(fp)
	if err != nil {
		return err
	}

	tarpart, err := writer.CreateFormFile("func", "func")
	if err != nil {
		return err
	}
	_, err = io.Copy(tarpart, functar)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", context.ServerURL+"/api/deploy", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to deploy function: %s", resp.Status)
	}

	// Parse the response body
	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Print(string(respbody))

	return nil
}

// createTar creates a gzip tarball of the function directory
func createTar(functionDir string) (*bytes.Buffer, error) {

	tarBuffer := new(bytes.Buffer)
	tw := tar.NewWriter(tarBuffer)
	defer tw.Close()

	err := filepath.Walk(functionDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(path)
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := io.Copy(tw, file); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	gz := gzip.NewWriter(tarBuffer)
	defer gz.Close()
	if _, err := gz.Write(tarBuffer.Bytes()); err != nil {
		return nil, err
	}
	return tarBuffer, nil
}

func marshalFaasYaml(fp string) ([]byte, error) {
	data, err := os.ReadFile(fp + "/dfaas.yaml")
	if err != nil {
		return nil, fmt.Errorf("unable to open dfaas.yaml: %s", err)
	}

	var BuildReq pb.BuildFunc

	err = yaml.Unmarshal(data, &BuildReq)
	if err != nil {
		return nil, err
	}

	if BuildReq.Name == "" {
		return nil, fmt.Errorf("no name found in dfaas.yaml")
	}

	if len(BuildReq.Commands) == 0 {
		return nil, fmt.Errorf("no commands found in dfaas.yaml")
	}

	if BuildReq.Runtime == "" {
		return nil, fmt.Errorf("no runtime found in dfaas.yaml")
	}

	metadata, err := proto.Marshal(&BuildReq)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}
