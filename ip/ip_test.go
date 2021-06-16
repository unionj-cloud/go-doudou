package ip

import (
	"net"
	"testing"
)

func TestGetOutboundIP(t *testing.T) {
	tests := []struct {
		name string
		want net.IP
	}{
		{
			name: "",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetOutboundIP(); len(got) == 0 {
				t.Errorf("GetOutboundIP() got nothing")
			}
		})
	}
}
