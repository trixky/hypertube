package server

import (
	pb "github.com/trixky/hypertube/api/proto"
)

type AuthServer struct {
	pb.AuthServiceServer
}
