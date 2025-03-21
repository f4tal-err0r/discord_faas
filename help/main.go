package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"faas.dev/function"
	pb "faas.dev/proto"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func funcWrapper(content *pb.DiscordContent) (*pb.DiscordResp, error) {
	log.Printf("%+v\n", content.Command)
	resp, err := function.Handler(content)
	if err != nil {
		return nil, fmt.Errorf("error executing function: %v", err)
	}
	return resp, nil
}

func main() {
	if len(os.Args) < 1 {
		log.Fatal("ERR: Address not provided")
		return
	}

	ctx := context.Background()

	conn, err := grpc.NewClient(fmt.Sprintf(":%s", os.Args[1]), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	client := pb.NewProcessorServiceClient(conn)
	stream, err := client.RecvContent(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	for {
		content, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}
		r, err := funcWrapper(content)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := client.SendResp(ctx, r); err != nil {
			log.Fatal(err)
		}
	}

}
