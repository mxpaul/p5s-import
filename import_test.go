package p5stock

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kylelemons/godebug/pretty"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/mxpaul/p5stock/mocks"
)

func serverAddr() string {
	return "127.0.0.1:10203"
}

func NewMockDb(t *testing.T, testCase *TestCase) (Done func(), mockDb *mocks.MockStorer, ctx context.Context) {
	mockCtrl, ctx := gomock.WithContext(context.Background(), t)

	mockDb = mocks.NewMockStorer(mockCtrl)
	testCase.MockUpsertSetup(mockDb)

	Done = func() {
		mockCtrl.Finish()
	}
	return
}

func GoFuckers() {
	_ = bson.M{}
	_ = log.Printf
}

type TestCase struct {
	Desc            string
	Prod            *FullProduct
	MockUpsertSetup func(mockDb *mocks.MockStorer)
	Want            ImportReply
}

func RunTestCase(testCase *TestCase, t *testing.T) {
	lis, err := net.Listen("tcp", serverAddr())
	if err != nil {
		t.Errorf("failed to listen: %v", err)
	}

	mockDone, mockStorer, ctx := NewMockDb(t, testCase)
	defer mockDone()

	grpcServer := grpc.NewServer()
	RegisterStockServer(grpcServer, &Stock{Db: mockStorer})
	go func(server *grpc.Server, lis net.Listener) {
		grpcServer.Serve(lis)
	}(grpcServer, lis)

	go func(ctx context.Context, t *testing.T, server *grpc.Server) {
		select {
		case <-ctx.Done():
			server.Stop()
		}
	}(ctx, t, grpcServer)

	conn, err := grpc.Dial(serverAddr(), grpc.WithInsecure())
	if err != nil {
		t.Errorf("Dial failed: %s", err)
	}
	client := NewStockClient(conn)

	got, err := client.ImportProduct(context.Background(), testCase.Prod)
	if err != nil {
		t.Errorf("ImportProduct error: %s", err)
	}
	log.Printf("Reply: %#v", got)
	if diff := pretty.Compare(got, testCase.Want); diff != "" {
		t.Fatalf("grpc reply unexpected: %s", diff)
	}

	grpcServer.Stop()
}

func TestDummy(t *testing.T) {
	prodId := "123"
	cases := []TestCase{
		TestCase{
			Desc: "first test case",
			Prod: &FullProduct{ProdId: prodId},
			MockUpsertSetup: func(mockDb *mocks.MockStorer) {
				mockDb.EXPECT().
					Upsert(bson.M{"prodid": prodId}, &FullProduct{ProdId: prodId}).
					Return(mgo.ChangeInfo{}, nil).
					Times(1)
			},
			Want: ImportReply{OK: true},
		},
		TestCase{
			Desc: "database upsert error",
			Prod: &FullProduct{ProdId: prodId},
			MockUpsertSetup: func(mockDb *mocks.MockStorer) {
				mockDb.EXPECT().
					Upsert(bson.M{"prodid": prodId}, &FullProduct{ProdId: prodId}).
					Return(mgo.ChangeInfo{}, fmt.Errorf("fuckup")).
					Times(1)
			},
			Want: ImportReply{OK: false, Description: "database error"},
		},
	}

	for _, testCase := range cases {
		RunTestCase(&testCase, t)
	}
}
