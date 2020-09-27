package topology

import (
	"encoding/xml"
)

type loc struct {
	dcName   string
	rackName string
}
type rack struct {
	Name string   `xml:"name,attr"`
	Ips  []string `xml:"Ip"`
}
type dataCenter struct {
	Name  string `xml:"name,attr"`
	Racks []rack `xml:"Rack"`
}
type topology struct {
	DataCenters []dataCenter `xml:"DataCenter"`
}
type Configuration struct {
	XMLName     xml.Name       `xml:"Configuration"`
	Topo        topology       `xml:"Topology"`
	ip2location map[string]loc // this is not used any more. leave it here for later.
}

func (c *Configuration) String() string {
	if b, e := xml.MarshalIndent(c, "  ", "  "); e == nil {
		return string(b)
	}
	return ""
}

func (c *Configuration) Locate(ip string, dcName string, rackName string) (dc string, rack string) {
	if c != nil && c.ip2location != nil {
		if loc, ok := c.ip2location[ip]; ok {
			return loc.dcName, loc.rackName
		}
	}

	if dcName == "" {
		dcName = "DefaultDataCenter"
	}

	if rackName == "" {
		rackName = "DefaultRack"
	}

	return dcName, rackName
}
