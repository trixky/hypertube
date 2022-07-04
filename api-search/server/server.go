package server

import (
	pb "github.com/trixky/hypertube/api-search/proto"
)

type SearchServer struct {
	pb.SearchServiceServer
}
