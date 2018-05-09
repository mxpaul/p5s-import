package main

import (
	"context"
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	stock "github.com/mxpaul/p5stock/ex/grpc/stock"
)

var (
	listenAddr = flag.String("listen", "127.0.0.1:10000", "Listen address:port")
	mongoAddr  = flag.String("mongo", "127.0.0.1:27017", "Mongo server address:port")
)

type Stock struct {
	// Mongo Connect
	Db *mgo.Collection
}

func (self *Stock) ImportProduct(ctx context.Context, prod *stock.FullProduct) (*stock.ImportReply, error) {
	info, err := self.Db.Upsert(bson.M{"prodid": prod.ProdId}, &prod)
	if err != nil {
		log.Printf("mongo error: %s", err)
		return &stock.ImportReply{OK: false, Description: err.Error()}, err
	}
	log.Printf("Upsert [ok] %s %s M:%d U:%d", prod.ProdId, prod.Name, info.Matched, info.Updated)
	return &stock.ImportReply{OK: true}, nil
}

func main() {
	//go:generate protoc stock.proto --go_out=plugins=grpc:.
	log.Printf("parse command line")
	flag.Parse()

	log.Printf("Listen %s for network connections", *listenAddr)
	lis, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Create mongo connector for %s", *mongoAddr)
	conn, err := mgo.Dial(*mongoAddr)
	if err != nil {
		log.Fatalf("Mongo connect error: %s", err)
	}
	defer conn.Close()
	mongoCollection := conn.DB("p5stest").C("stock")

	grpcServer := grpc.NewServer()
	stock.RegisterStockServer(grpcServer, &Stock{Db: mongoCollection})
	//... // determine whether to use TLS
	log.Printf("Ready to process grpc requests")
	grpcServer.Serve(lis)
	log.Printf("Exiting...")
}
