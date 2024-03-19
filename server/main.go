package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sample-manager/constants"
	"sample-manager/models"
	pb "sample-manager/proto"

	"sample-manager/db"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)
type Server struct {
	pb.SampleManagerServer
	DB *gorm.DB
}

func main() {
	db := database.Connection()
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSampleManagerServer(grpcServer, &Server{DB: db})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
	}
	log.Println("Server started running on port 9000")
}

func (s *Server) GetSampleId(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	sample_item_id, err := database.GetSampleID(s.DB, req.Clm, req.ItemId)

	if err != nil {
		errorString := fmt.Sprintf("No mapping found: %v", err)
		return nil, status.Error(codes.Unavailable, errorString)
	}

	resp := pb.GetResponse {
		SampleItemId: sample_item_id,
	}

	return &resp, nil
}

func (s *Server) CreateMapping(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	var segments []models.Segment
	for _, s := range(req.Segments) {
		segment := models.Segment {
			Segment: s,
		}
		segments = append(segments, segment)
	}
	mapping := &models.Mapping {
		Segments: segments,
		SampleItemID: req.SampleItemId,
		ItemID: req.ItemId,
	}

	err := s.DB.Create(&mapping).Error

	if err != nil {
		errorString := fmt.Sprintf("Error storing the mapping: %v", err)
		return nil, status.Errorf(codes.Unknown, errorString)
	}

	response := &pb.CreateResponse{
		Message: constants.CREATE_MAPPING_SUCCESS,
	}

	return response, nil
}