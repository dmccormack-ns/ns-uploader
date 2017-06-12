package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// <soapenv:Envelope
//     xmlns:xsd='http://www.w3.org/2001/XMLSchema'
//     xmlns:xsi='http://www.w3.org/2001/XMLSchema-instance'
//     xmlns:soapenv='http://schemas.xmlsoap.org/soap/envelope/'
//     xmlns:platformCore='urn:core_2017_1.platform.webservices.netsuite.com'
//     xmlns:platformMsgs='urn:messages_2017_1.platform.webservices.netsuite.com'
//     xmlns:docFileCab='urn:filecabinet_2017_1.documents.webservices.netsuite.com'>
//     <soapenv:Header>
//         <passport xsi:type='platformCore:Passport'>
//             <email xsi:type='xsd:string'>dmccormack@netsuite.com</email>
//             <password xsi:type='xsd:string'>#*na9RdE3AqQuOG0</password>
//             <account xsi:type='xsd:string'>TSTDRV1620460</account>
//         </passport>
//         <applicationInfo xsi:type='platformMsgs:ApplicationInfo'>
//             <applicationId xsi:type='xsd:string'>6CEAF778-AA4C-4A21-9E59-510B940360FB</applicationId>
//         </applicationInfo>
//     </soapenv:Header>
//     <soapenv:Body>
//         <update xsi:type='platformMsgs:UpdateRequest'>
//             <record xsi:type='docFileCab:File' internalId='58725'>
//                 <altTagCaption xsi:type='xsd:string'>test_soap_update</altTagCaption>
//             </record>
//         </update>
//     </soapenv:Body>
// </soapenv:Envelope>

func main() {
	type BodyField struct {
		XMLName xml.Name
		Type    string `xml:"xsi:type,attr"`
		Value   string `xml:",innerxml"`
	}

	type Record struct {
		XMLName    xml.Name `xml:"record"`
		Type       string   `xml:"xsi:type,attr"`
		InternalID string   `xml:"internalId,attr"`
		BodyField  BodyField
	}

	type Update struct {
		XMLName xml.Name `xml:"update"`
		Type    string   `xml:"xsi:type,attr"`
		Record  Record
	}

	type Body struct {
		XMLName xml.Name `xml:"soapenv:Body"`
		Action  Update
	}
	type Element struct {
		XMLName xml.Name
		Type    string `xml:"xsi:type,attr"`
		Value   string `xml:",innerxml"`
	}
	type Passport struct {
		XMLName  xml.Name `xml:"passport"`
		Id       string   `xml:"xsi:type,attr"`
		Email    Element
		Password Element
		Account  Element
	}
	type ApplicationInfo struct {
		XMLName       xml.Name `xml:"applicationInfo"`
		Type          string   `xml:"xsi:type,attr"`
		ApplicationID Element
	}

	type Header struct {
		XMLName         xml.Name `xml:"soapenv:Header"`
		Passport        Passport
		ApplicationInfo ApplicationInfo
	}

	type Envelope struct {
		XMLName      xml.Name `xml:"soapenv:Envelope"`
		Xsd          string   `xml:"xmlns:xsd,attr"`
		Xsi          string   `xml:"xmlns:xsi,attr"`
		Soapenv      string   `xml:"xmlns:soapenv,attr"`
		PlatformCore string   `xml:"xmlns:platformCore,attr"`
		PlatformMsgs string   `xml:"xmlns:platformMsgs,attr"`
		DocFileCab   string   `xml:"xmlns:docFileCab,attr"`
		Header       Header
		Body         Body
	}
	//     <soapenv:Body>
	//         <update xsi:type='platformMsgs:UpdateRequest'>
	//             <record xsi:type='docFileCab:File' internalId='58725'>
	//                 <altTagCaption xsi:type='xsd:string'>test_soap_update</altTagCaption>
	//             </record>
	//         </update>
	//     </soapenv:Body>

	a := ApplicationInfo{
		XMLName: xml.Name{Local: "applicationInfo"},
		Type:    "platformMsgs:ApplicationInfo",
		ApplicationID: Element{
			XMLName: xml.Name{Local: "applicationId"},
			Type:    "xsd:string",
			Value:   "6CEAF778-AA4C-4A21-9E59-510B940360FB",
		},
	}

	b := Body{
		XMLName: xml.Name{Local: "soapenv:Body"},
		Action: Update{
			XMLName: xml.Name{Local: "update"},
			Type:    "platformMsgs:UpdateRequest",
			Record: Record{
				XMLName:    xml.Name{Local: "record"},
				Type:       "docFileCab:File",
				InternalID: "58725",
				BodyField: BodyField{
					XMLName: xml.Name{Local: "altTagCaption"},
					Type:    "xsd:string",
					Value:   "test_soap_update_from_go",
				},
			},
		},
	}

	p := Passport{
		XMLName: xml.Name{Local: "passport"},
		Id:      "platformCore:Passport",
		Email: Element{
			XMLName: xml.Name{Local: "email"},
			Type:    "xsd:string",
			Value:   "",
		},
		Password: Element{
			XMLName: xml.Name{Local: "password"},
			Type:    "xsd:string",
			Value:   "",
		},
		Account: Element{
			XMLName: xml.Name{Local: "account"},
			Type:    "xsd:string",
			Value:   "",
		},
	}
	h := Header{
		XMLName:         xml.Name{Local: "soapenv:Header"},
		Passport:        p,
		ApplicationInfo: a,
	}

	e := Envelope{
		XMLName:      xml.Name{Local: "soapenv:Envelope"},
		Xsd:          "http://www.w3.org/2001/XMLSchema",
		Xsi:          "http://www.w3.org/2001/XMLSchema-instance",
		Soapenv:      "http://schemas.xmlsoap.org/soap/envelope/",
		PlatformCore: "urn:core_2017_1.platform.webservices.netsuite.com",
		PlatformMsgs: "urn:messages_2017_1.platform.webservices.netsuite.com",
		DocFileCab:   "urn:filecabinet_2017_1.documents.webservices.netsuite.com",
		Header:       h,
		Body:         b,
	}
	output, err := xml.MarshalIndent(e, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	headerBytes := []byte(xml.Header)
	output = append(headerBytes, output...)
	//os.Stdout.Write([]byte(xml.Header))
	req, err := http.NewRequest("POST", "https://webservices.na1.netsuite.com/services/NetSuitePort_2017_1", bytes.NewBuffer(output))
	req.Header.Set("Content-Type", "text/xml")
	req.Header.Set("SoapAction", "update")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}
