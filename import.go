package p5stock

//go:generate sh -c "protoc *.proto --go_out=plugins=grpc:."

import (
	"context"
	//"google.golang.org/grpc"
)

type Stock struct {
	// Mongo Connect
	//Db *mgo.Collection
}

func (self *Stock) ImportProduct(ctx context.Context, prod *FullProduct) (*ImportReply, error) {
	return &ImportReply{}, nil
	//info, err := self.Db.Upsert(bson.M{"prodid": prod.ProdId}, &prod)
	//if err != nil {
	//	log.Printf("mongo error: %s", err)
	//	return &stock.ImportReply{OK: false, Description: err.Error()}, err
	//}
	//log.Printf("Upsert [ok] %s %s M:%d U:%d", prod.ProdId, prod.Name, info.Matched, info.Updated)
	//return &stock.ImportReply{OK: true}, nil
}
