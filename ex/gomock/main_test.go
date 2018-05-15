package main

import (
	//"context"
	//"log"
	"testing"

	"github.com/golang/mock/gomock"
	//"github.com/kylelemons/godebug/pretty"
)

func TestExample(t *testing.T) {
	//mockCtrl, ctx := gomock.WithContext(context.Background(), t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	//go func(ctx context.Context, t *testing.T, mock *gomock.Controller) {
	//	select {
	//	case <-ctx.Done():
	//		//t.Errorf("context done: %v", ctx.Err())
	//		//mock.Finish()
	//		t.Errorf("Context cancel: %s", ctx.Err())
	//		panic("CONTEXT CANCEL")
	//	}
	//}(ctx, t, mockCtrl)

	mockStorer := NewMockStorer(mockCtrl)
	mockStorer.EXPECT().
		//Upsert(bson.M{"prodid": "123"}, FullProduct{ProdId: "123"}).
		Upsert(1234).
		Return(0, nil)
	//.
	//		Times(1)
	user := User{Db: mockStorer}
	user.TestCall()
}
