package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	XML_FILE = "../../tmp/p5s-one.xml"
)

type XMLka struct {
	//XMLName  xml.Name  `xml:"p5s"`
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
	xmlFile, err := os.Open(XML_FILE)
	if err != nil {
		fmt.Println("Error opening file '%s': %s", XML_FILE, err)
		return
	}
	defer xmlFile.Close()
	buf, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Println("Error reading xml data: %s", err)
		return
	}
	//fmt.Printf("File content: %s\n", buf.String())

	var Struct XMLka
	err = xml.Unmarshal(buf, &Struct)
	if err != nil {
		fmt.Println("Error parsing xml: %s", err)
		return
	}
	fmt.Printf("list: %#v\n", Struct)

}
