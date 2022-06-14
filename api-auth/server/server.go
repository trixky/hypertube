package server

import (
	pb "github.com/trixky/hypertube/api-auth/proto"
)

type AuthServer struct {
	pb.AuthServiceServer
}
