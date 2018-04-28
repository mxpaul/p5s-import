package main

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Product struct {
	Name         string `xml:"name,attr"`
	Description  string `xml:"description"`
	PictureSmall string `xml:"pictureSmall"`
}

func main() {
	prod := Product{
		Name:         "Good long product",
		Description:  "Very long, very good, check it up",
		PictureSmall: "http://pic.io/long.jpg",
	}

	conn, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		log.Fatalf("Mongo connect error: %s", err)
	}
	defer conn.Close()

	ccnn := conn.DB("p5stest").C("stock")

	info, err := ccnn.Upsert(bson.M{"name": "Good long product"}, &prod)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Upsert ChangeInfo: %#v", info)

	result := Product{}
	err = ccnn.Find(bson.M{"name": "Good long product"}).One(&result)
	if err != nil {
		log.Fatalf("Select error: %s", err)
	}

	log.Printf("Result: %#v", result)

}
