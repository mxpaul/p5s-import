package p5stock

//go:generate sh -c "protoc *.proto --go_out=plugins=grpc:."
//go:generate mockgen -source import.go -package mocks -destination mocks/storer.mock.go Storer

import (
	"context"
	//"google.golang.org/grpc"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Storer interface {
	Upsert(interface{}, interface{}) (mgo.ChangeInfo, error)
}

type Stock struct {
	// Mongo Connect
	Db Storer
}

func (self *Stock) ImportProduct(ctx context.Context, prod *FullProduct) (*ImportReply, error) {
	_, err := self.Db.Upsert(bson.M{"prodid": prod.ProdId}, prod)
	if err != nil {
		return &ImportReply{OK: false, Description: "database error"}, nil
	}
	return &ImportReply{OK: true}, nil
}
