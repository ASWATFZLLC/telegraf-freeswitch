package utils

import (
	"bytes"
	"encoding/xml"
	"strings"

	"golang.org/x/net/html/charset"
)

type Gateway struct {
	Name     string `xml:"name"`
	Profile  string `xml:"profile"`
	Realm    string `xml:"realm"`
	PingTime string `xml:"pingtime"`
	State    string `xml:"state"`
	Status   string `xml:"status"`
}
type Gateways struct {
	Gateways []*Gateway `xml:"gateway"`
}

type SofiaGateway struct {
	Name   string
	Ping   string
	Status string
}

func (sp *SofiaGateway) loadXMLGateway(p *Gateway) error {
	sp.Name = p.Name
	sp.Status = "0"
	if p.Status == "UP" {
		sp.Status = "1"
	}
	sp.Ping = p.PingTime
	return nil
}

func ParseSofiaStatusGateways(data string) ([]*SofiaGateway, error) {
	data = strings.TrimSpace(data)
	dec := xml.NewDecoder(bytes.NewBufferString(data))
	dec.CharsetReader = charset.NewReaderLabel

	gateways := &Gateways{}
	err := dec.Decode(gateways)
	if err != nil {
		return nil, err
	}
	sofiaGateways := make([]*SofiaGateway, len(gateways.Gateways))
	for i, p := range gateways.Gateways {
		sp := &SofiaGateway{}
		if err := sp.loadXMLGateway(p); err != nil {
			return nil, err
		}
		sofiaGateways[i] = sp
	}
	return sofiaGateways, nil
}
