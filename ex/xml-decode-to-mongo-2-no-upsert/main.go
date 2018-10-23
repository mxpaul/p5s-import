package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	XML_FILE = "../../tmp/p5s_full.xml"
)

type XMLka struct {
	FileDate string    `xml:"FileDate,attr"`
	Products []Product `xml:"product"`
}

type Product struct {
	Name                 string     `xml:"name,attr"`
	Description          string     `xml:"description"`
	PictureSmall         string     `xml:"pictureSmall"`
	Price                Price      `xml:"price"`
	ProdId               string     `xml:"prodID,attr"`
	Variants             []Variant  `xml:"assortiment>assort"`
	Pictures             []Picture  `xml:"pictures>picture"`
	VendorProdId         string     `xml:"vendorCode,attr"`
	VendorName           string     `xml:"vendor,attr"`
	VendorCollectionName string     `xml:"CollectionName,attr"`
	Material             string     `xml:"material,attr"`
	Diameter             float32    `xml:"diameter,attr"`
	Length               float32    `xml:"lenght,attr"`
	Categories           []Category `xml:"categories>category"`
}

type Price struct {
	WholePrice      float32 `xml:"WholePrice,attr"`
	BaseWholePrice  float32 `xml:"BaseWholePrice,attr"`
	BaseRetailPrice float32 `xml:"BaseRetailPrice,attr"`
	RetailPrice     float32 `xml:"RetailPrice,attr"`
	Discount        float32 `xml:"Discount,attr"`
	StopPromo       float32 `xml:"StopPromo,attr"`
}

type Variant struct {
	AID          string `xml:"aID,attr"`
	Count        int    `xml:"sklad,attr"`
	ShippingDate string `xml:"ShippingDate,attr"`
}

type Picture struct {
	URL string `xml:",innerxml"`
}

type Category struct {
	Name    string `xml:"Name,attr"`
	SubName string `xml:"subName,attr"`
}

func main() {
	conn, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		log.Fatalf("Mongo connect error: %s", err)
	}
	defer conn.Close()
	ccnn := conn.DB("p5stest").C("stock")

	xmlFile, err := os.Open(XML_FILE)
	if err != nil {
		fmt.Println("Error opening file '%s': %s", XML_FILE, err)
		return
	}
	defer xmlFile.Close()
	decoder := xml.NewDecoder(xmlFile)
	var inElement string
	for {
		t, err := decoder.Token()
		if err != nil && err != io.EOF {
			log.Fatalf("xml parse error: %s", err)
		}
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			inElement = se.Name.Local
			var prod Product
			if inElement == "product" {
				decoder.DecodeElement(&prod, &se)
				//info, err := ccnn.Upsert(bson.M{"prodid": prod.ProdId}, &prod)
				err := importProductToMongo(ccnn, &prod)
				if err != nil {
					log.Fatal(err)
				}
				//log.Printf("Upsert [ok] %s %s M:%d U:%d", prod.ProdId, prod.Name, info.Matched, info.Updated)
			}
		default:
		}
	}

}

func importProductToMongo(ccnn *mgo.Collection, prod *Product) error {
	query := ccnn.Find(bson.M{"prodid": prod.ProdId})
	existingProd := Product{}
	err := query.One(&existingProd)
	switch err {
	case mgo.ErrNotFound: // Import new product
		log.Printf("not found, insert new product: [%s] %s", prod.ProdId, prod.Name)
		err = ccnn.Insert(prod)
		return err
		//	return err
	case nil: // prodId already exists
		log.Printf("prod exists, update product: [%s] %s", prod.ProdId, prod.Name)
		err = ccnn.Update(bson.M{"prodid": prod.ProdId}, prod)
		return err
	default: // Real error
		return err
	}
	//log.Printf("Query result: %v", query)
	////info, err := ccnn.Upsert(bson.M{"prodid": prod.ProdId}, &prod)
	////if err != nil {
	////	return err
	////}
	////log.Printf("Upsert [ok] %s %s M:%d U:%d", prod.ProdId, prod.Name, info.Matched, info.Updated)
	//return nil
}
