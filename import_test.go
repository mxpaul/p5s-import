package p5stock

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
)

func serverAddr() string {
	return "127.0.0.1:10203"
}

func TestDummy(t *testing.T) {
	_ = log.Printf
	lis, err := net.Listen("tcp", serverAddr())
	if err != nil {
		t.Errorf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	RegisterStockServer(grpcServer, &Stock{})
	go grpcServer.Serve(lis)

	conn, err := grpc.Dial(serverAddr(), grpc.WithInsecure())
	if err != nil {
		t.Errorf("Dial failed: %s", err)
	}
	client := NewStockClient(conn)

	prod := FullProduct{ProdId: "123"}
	reply, err := client.ImportProduct(context.Background(), &prod)
	if err != nil {
		t.Errorf("ImportProduct error: %s", err)
	}
	log.Printf("Reply: %#v", reply)

	grpcServer.Stop()
}
