package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"

	p5s "github.com/mxpaul/p5s-import"
	stock "github.com/mxpaul/p5s-import/ex/grpc/stock"
)

var (
	serverAddr = flag.String("server", "127.0.0.1:10000", "Listen address:port")
	XML_FILE   = flag.String("file", "../../tmp/p5s-one.xml", "xml file to import")
)

func main() {
	//go:generate protoc stock.proto --go_out=plugins=grpc:.
	flag.Parse()
	// Connect to server
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial failed: %s", err)
	}
	defer conn.Close()
	client := stock.NewStockClient(conn)
	log.Printf("Ready to send requests")

	// Open XML file to import
	xmlFile, err := os.Open(*XML_FILE)
	if err != nil {
		fmt.Println("Error opening file '%s': %s", *XML_FILE, err)
		return
	}
	defer xmlFile.Close()
	log.Printf("XML file %s open success", *XML_FILE)
	importer := p5s.NewStreamedImporter(xmlFile)

	for p := range importer.C {
		prod := stock.FullProduct{Name: p.Name, ProdId: p.ProdId}
		reply, err := client.ImportProduct(context.Background(), &prod)
		if err != nil {
			log.Printf("ImportProduct error: %s", err)
			os.Exit(1)
		}
		if !reply.OK {
			log.Printf("Upsert [fail] %s %s: %#v", prod.ProdId, prod.Name, reply)
			os.Exit(1)
		}
		log.Printf("Upsert [OK] %s %s Stat:%#v", prod.ProdId, prod.Name, reply.Stat)
	}

}
