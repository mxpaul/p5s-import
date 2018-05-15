package main

//go:generate mockgen -source=main.go -destination=mock_storer.go -package=main Storer

import (
	"fmt"
)

type Storer interface {
	Upsert(interface{}) (int, error)
}

type User struct {
	Db Storer
}

func (self *User) TestCall() {
	self.Db.Upsert(123)
}

func main() {
	fmt.Println("")
}
