package main

import (
	"fmt"
	"io"
	"net/http"

	"func/function"

	"google.golang.org/protobuf/proto"
)

func funcWrapper(w http.ResponseWriter, req *http.Request) {

	content, err := unmarshalRequest(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(content.Command)
	resp, err := function.Handler(content)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Message)
	w.Write([]byte(resp.Message))
}

func unmarshalRequest(r *http.Request) (*function.DiscordContent, error) {
	// Parse the request body into a DiscordContent struct
	var content function.DiscordContent

	//body to bytes
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = proto.Unmarshal(data, &content)
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func main() {
	//create a http router on 8080
	http.HandleFunc("/", funcWrapper)
	fmt.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
