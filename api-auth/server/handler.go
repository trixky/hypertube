package server

import (
	"context"
	"fmt"
	"log"

	pb "github.com/trixky/hypertube/api/proto"
	grpcMetadata "google.golang.org/grpc/metadata"
)

func (s *AuthServer) InternalRegister(ctx context.Context, in *pb.InternalRegisterRequest) (*pb.InternalLoginResponse, error) {
	log.Printf("Greet function awas invoked with %v\n", in)

	// -----------------------
	md, ok := grpcMetadata.FromIncomingContext(ctx)

	if !ok {
		log.Println("arg 111")
		return nil, nil
	}

	auth := md.Get("auth")
	chat := md.Get("chat")

	fmt.Println("auth:", auth)
	fmt.Println("chat:", chat)

	// -----------------------
	// grpcMetadata.Pairs()
	// -----------------------
	// toto, ok := grpcMetadata.FromIncomingContext(ctx)

	// if !ok {
	// 	fmt.Println("pas content")
	// }

	// strs := toto.Get("asdf")

	// fmt.Println("strs:", strs)
	// -----------------------

	return &pb.InternalLoginResponse{
		Jwt: "ouiiiiiii",
	}, nil
}
