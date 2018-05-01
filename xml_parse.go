package p5s

import (
	"encoding/xml"
	"io"
	"log"
)

type XMLProduct struct {
	ProdId               string        `xml:"prodID,attr"`
	Name                 string        `xml:"name,attr"`
	Description          string        `xml:"description"`
	PictureSmall         string        `xml:"pictureSmall"`
	Price                XMLPrice      `xml:"price"`
	Variants             []XMLVariant  `xml:"assortiment>assort"`
	Pictures             []XMLPicture  `xml:"pictures>picture"`
	VendorProdId         string        `xml:"vendorCode,attr"`
	VendorName           string        `xml:"vendor,attr"`
	VendorCollectionName string        `xml:"CollectionName,attr"`
	Material             string        `xml:"material,attr"`
	Diameter             float32       `xml:"diameter,attr"`
	Length               float32       `xml:"lenght,attr"`
	Categories           []XMLCategory `xml:"categories>category"`
}

type XMLPrice struct {
	WholePrice      float32 `xml:"WholePrice,attr"`
	BaseWholePrice  float32 `xml:"BaseWholePrice,attr"`
	BaseRetailPrice float32 `xml:"BaseRetailPrice,attr"`
	RetailPrice     float32 `xml:"RetailPrice,attr"`
	Discount        float32 `xml:"Discount,attr"`
	StopPromo       float32 `xml:"StopPromo,attr"`
}

type XMLVariant struct {
	AID          string `xml:"aID,attr"`
	Count        int    `xml:"sklad,attr"`
	ShippingDate string `xml:"ShippingDate,attr"`
}

type XMLPicture struct {
	URL string `xml:",innerxml"`
}

type XMLCategory struct {
	Name    string `xml:"Name,attr"`
	SubName string `xml:"subName,attr"`
}

type ProductStreamChan chan XMLProduct
type XMLStreamImporter struct {
	r io.Reader
	C ProductStreamChan
}

func NewStreamedImporter(r io.Reader) *XMLStreamImporter {
	self := XMLStreamImporter{
		r: r,
		C: make(ProductStreamChan, 2),
	}
	go func(self *XMLStreamImporter) {
		//	_, err := ioutil.ReadAll(self.r)
		//	if err != nil {
		//		return
		//	}
		//	self.C <- XMLProduct{}

		decoder := xml.NewDecoder(self.r)
		var inElement string
		for {
			t, err := decoder.Token()
			if err != nil && err != io.EOF {
				log.Fatalf("xml parse error: %s", err)
				break
			}
			if t == nil {
				break
			}

			switch se := t.(type) {
			case xml.StartElement:
				inElement = se.Name.Local
				var prod XMLProduct
				if inElement == "product" {
					decoder.DecodeElement(&prod, &se)
					self.C <- prod
					//info, err := ccnn.Upsert(bson.M{"prodid": prod.ProdId}, &prod)
					//if err != nil {
					//	log.Fatal(err)
					//}
					//log.Printf("Upsert [ok] %s %s M:%d U:%d", prod.ProdId, prod.Name, info.Matched, info.Updated)
				}
			default:
			}
		}

		close(self.C)
	}(&self)
	return &self
}
