package handler

import (
	"context"
	log "gin-template/pkg/logger"
	"io"
	"time"

	pb "gin-template/service/proto"
)

type Service struct{}

func (e *Service) Call(ctx context.Context, req *pb.CallRequest, rsp *pb.CallResponse) error {
	log.Info("Received Service.Call request: %v", req)
	rsp.Msg = "Hello " + req.Name
	return nil
}

func (e *Service) ClientStream(ctx context.Context, stream pb.Service_ClientStreamStream) error {
	var count int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Info("Got %v pings total", count)
			return stream.SendMsg(&pb.ClientStreamResponse{Count: count})
		}
		if err != nil {
			return err
		}
		log.Info("Got ping %v", req.Stroke)
		count++
	}
}

func (e *Service) ServerStream(ctx context.Context, req *pb.ServerStreamRequest, stream pb.Service_ServerStreamStream) error {
	log.Info("Received Service.ServerStream request: %v", req)
	for i := 0; i < int(req.Count); i++ {
		log.Info("Sending %d", i)
		if err := stream.Send(&pb.ServerStreamResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 250)
	}
	return nil
}

func (e *Service) BidiStream(ctx context.Context, stream pb.Service_BidiStreamStream) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Info("Got ping %v", req.Stroke)
		if err := stream.Send(&pb.BidiStreamResponse{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
