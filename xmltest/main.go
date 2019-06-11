package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

type XMLFucks struct {
	XMLName     xml.Name `xml:"http://faas.unnecessary.tech/schema ListOfFucks"`
	Status      string   `xml:"status"`
	Fucks       []string `xml:"fucks>item"`
	Observation string   `xml:"observation"`
}

func main() {
	w := &bytes.Buffer{}
	enc := xml.NewEncoder(w)
	procInst := xml.ProcInst{
		Target: "xml",
		Inst:   []byte("version=\"1.0\" encoding=\"UTF-8\""),
	}
	if err := enc.EncodeToken(procInst); err != nil {
		panic(err)
	}

	ssInst := xml.ProcInst{
		Target: "xml-stylesheet",
		Inst:   []byte("type=\"text/xsl\" href=\"https://faas.unnecessary.tech/fuckformat.xslt\""),
	}
	if err := enc.EncodeToken(ssInst); err != nil {
		panic(err)
	}

	v := XMLFucks{
		Status:      "ok",
		Fucks:       []string{"fuck", "fuck", "fuck"},
		Observation: "Why the fuck are you still using XML?",
	}

	if err := enc.Encode(v); err != nil {
		panic(err)
	}
	fmt.Println(w.String())
}
