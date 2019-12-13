package hikari

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/types"
)

// ServiceName is the name of the current microservice
const ServiceName = types.Hikari

// storage stores the DNS A records in the form of Key : Value pairs
// with Domain Name as the key and the IPv4 Address as the value
var storage = types.NewRecordStorage()

type handler struct{}

func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	// Only serve A records
	if r.Question[0].Qtype == dns.TypeA {
		msg.Authoritative = true
		domain := msg.Question[0].Name
		if address, ok := storage.Get(domain); ok {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(address),
			})
			w.WriteMsg(&msg)
		}
	}
}

// NewService returns a new instance of the current microservice
func NewService() *dns.Server {
	server := &dns.Server{
		Addr: fmt.Sprintf(":%d", configs.ServiceConfig.Hikari.Port),
		Net:  "udp",
	}
	server.Handler = &handler{}
	return server
}
