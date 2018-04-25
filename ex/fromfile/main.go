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
	PictureSmall string `xml:"pictureSmall"`
	Price        Price  `xml:"price"`
}

type Price struct {
	WholePrice     float32 `xml:"WholePrice,attr"`
	BaseWholePrice float32 `xml:"BaseWholePrice,attr"`
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
