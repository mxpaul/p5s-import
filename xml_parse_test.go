package p5stock

import (
	"log"
	//"reflect"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestHelloWorld(t *testing.T) {
	// t.Fatal("not implemented")
}

var singleProductXML = `
<?xml version="1.0" encoding="utf-8" ?>
<p5s FileDate="22:31:47 24.04.2018">
<product prodID="11"  vendorCode="VC" name="Продукт 1" vendor="vendor name"
brutto="" batteries="" pack="box" material="материал 2" lenght="14.00" diameter="1.50" CollectionName="коллекция название">
	<description>Description text. [It goes here]!</description>
	<categories>
		<category Name="категория" subName="подкатегория"/>
	</categories>
	<price RetailPrice="160.00" BaseRetailPrice="310.00" WholePrice="75.00" BaseWholePrice="150.00" Discount="50" StopPromo="0"/>
	<pictureSmall>http://im.io/small.jpg</pictureSmall>
	<pictures>
		<picture>http://im.io/big1.jpg</picture>
		<picture>http://im.io/big2.jpg</picture>
		<picture>http://im.io/big3.jpg</picture>
		<picture>http://im.io/big4.jpg</picture>
		<picture>http://im.io/big5.jpg</picture>
		<picture>http://im.io/big6.jpg</picture>
	</pictures>
	<assortiment>
		<assort aID="25" sklad="0" color="red" size="" barcode="123456780123" ShippingDate=""/>
	</assortiment>
</product>
</p5s>
`
var parsedProduct = XMLProduct{
	ProdId:               "11",
	Name:                 "Продукт 1",
	Description:          "Description text. [It goes here]!",
	PictureSmall:         "http://im.io/small.jpg",
	VendorProdId:         "VC",
	VendorName:           "vendor name",
	VendorCollectionName: "коллекция название",
	Material:             "материал 2",
	Diameter:             1.50,
	Length:               14,
	Price: XMLPrice{
		BaseWholePrice:  150.0,
		BaseRetailPrice: 310.0,
		WholePrice:      75.0,
		RetailPrice:     160.0,
		StopPromo:       0,
		Discount:        50,
	},
	Pictures: []XMLPicture{
		XMLPicture{URL: "http://im.io/big1.jpg"},
		XMLPicture{URL: "http://im.io/big2.jpg"},
		XMLPicture{URL: "http://im.io/big3.jpg"},
		XMLPicture{URL: "http://im.io/big4.jpg"},
		XMLPicture{URL: "http://im.io/big5.jpg"},
		XMLPicture{URL: "http://im.io/big6.jpg"},
	},
	Variants: []XMLVariant{
		XMLVariant{AID: "25", Count: 0, ShippingDate: ""},
	},
	Categories: []XMLCategory{
		XMLCategory{Name: "категория", SubName: "подкатегория"},
	},
}

func TestReadOneProduct(t *testing.T) {
	xmlreader := strings.NewReader(singleProductXML)
	importer := NewStreamedImporter(xmlreader)

	gotList := make([]XMLProduct, 0, 0)
	for prod := range importer.C {
		gotList = append(gotList, prod)
	}
	//log.Printf("product list: %v", list)
	_ = log.Printf
	wantList := []XMLProduct{parsedProduct}
	//if !reflect.DeepEqual(gotList, wantList) {
	if diff := pretty.Compare(gotList, wantList); diff != "" {
		t.Fatalf("Got not what we want: %s", diff)
	}
}
