package memberlist

import (
	"time"

	"github.com/miekg/dns"
)

//go:generate mockgen -destination ./mock/mock_dns_client_interface.go -package mock -source=./dns_client_interface.go

type IDNSClient interface {
	Exchange(m *dns.Msg, address string) (r *dns.Msg, rtt time.Duration, err error)
}
