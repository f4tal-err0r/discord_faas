package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	v1 "github.com/f4tal-err0r/discord_faas/api/v1"
	"github.com/olekukonko/tablewriter"
)

func ListFunctions() {
	c, err := GetCurrentContext()
	if err != nil {
		log.Fatalf("Unable to get current context: %v", err)
	}

	//set auth header
	funcurl, err := url.JoinPath(c.ServerUrl, "/api/functions", c.GuildId)
	if err != nil {
		log.Fatalf("Unable to create function list URL: %v", err)
	}
	req, err := http.NewRequest(http.MethodGet, funcurl, nil)
	if err != nil {
		log.Fatalf("Unable to create request: %v", err)
	}

	req.Header.Add("Authorization", c.JwToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error executing request: %v, %+v", err, resp)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unable to get functions: %s", resp.Status)
	}

	var funcList v1.ListFunctions

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Unable to read response body: %v", err)
	}

	err = json.Unmarshal(body, &funcList)
	if err != nil {
		log.Fatalf("Unable to unmarshal response: %v", err)
	}

	header := []string{"Name", "Description", "Runtime", "Version"}
	tw := tablewriter.NewWriter(os.Stdout)
	tw.Header(header)
	fmt.Println("Available Functions:")
	for _, f := range funcList.Functions {
		tw.Append([]string{f.Name, f.Description, f.Runtime, f.Version})
	}
	tw.Render()
}
