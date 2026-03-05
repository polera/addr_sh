package addr

import (
	"net/http"
	"testing"
)

func TestCalculateV4CIDR(t *testing.T) {
	tests := []struct {
		name           string
		cidr           string
		wantAddresses  float64
		wantUsable     float64
		wantErr        bool
	}{
		{"slash 0", "0", 4294967296, 4294967294, false},
		{"slash 16", "16", 65536, 65534, false},
		{"slash 24", "24", 256, 254, false},
		{"slash 32", "32", 1, 1, false},
		{"negative", "-1", 0, 0, true},
		{"too large", "33", 0, 0, true},
		{"non-numeric", "abc", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateV4CIDR(tt.cidr)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CalculateV4CIDR(%q) error = %v, wantErr %v", tt.cidr, err, tt.wantErr)
			}
			if !tt.wantErr {
				if got.NumAddresses != tt.wantAddresses {
					t.Errorf("NumAddresses = %v, want %v", got.NumAddresses, tt.wantAddresses)
				}
				if got.NumUsableAddresses != tt.wantUsable {
					t.Errorf("NumUsableAddresses = %v, want %v", got.NumUsableAddresses, tt.wantUsable)
				}
			}
		})
	}
}

func TestGetIpv4CIDR(t *testing.T) {
	tests := []struct {
		name          string
		addr          string
		cidr          string
		wantNetwork   string
		wantBroadcast string
		wantFirst     string
		wantLast      string
		wantHostOnly  bool
		wantErr       bool
	}{
		{
			name:          "slash 24",
			addr:          "192.168.1.0",
			cidr:          "24",
			wantNetwork:   "192.168.1.0/24",
			wantBroadcast: "192.168.1.255",
			wantFirst:     "192.168.1.1",
			wantLast:      "192.168.1.254",
		},
		{
			name:          "slash 16",
			addr:          "10.0.0.0",
			cidr:          "16",
			wantNetwork:   "10.0.0.0/16",
			wantBroadcast: "10.0.255.255",
			wantFirst:     "10.0.0.1",
			wantLast:      "10.0.255.254",
		},
		{
			name:         "slash 32 host only",
			addr:         "192.168.1.1",
			cidr:         "32",
			wantNetwork:  "192.168.1.1/32",
			wantFirst:    "192.168.1.1",
			wantLast:     "192.168.1.1",
			wantHostOnly: true,
		},
		{
			name:    "invalid address",
			addr:    "not-an-ip",
			cidr:    "24",
			wantErr: true,
		},
		{
			name:    "invalid cidr",
			addr:    "192.168.1.0",
			cidr:    "abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIpv4CIDR(tt.addr, tt.cidr)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetIpv4CIDR(%q, %q) error = %v, wantErr %v", tt.addr, tt.cidr, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got.Network != tt.wantNetwork {
				t.Errorf("Network = %q, want %q", got.Network, tt.wantNetwork)
			}
			if got.FirstUsableAddress != tt.wantFirst {
				t.Errorf("FirstUsableAddress = %q, want %q", got.FirstUsableAddress, tt.wantFirst)
			}
			if got.LastUsableAddress != tt.wantLast {
				t.Errorf("LastUsableAddress = %q, want %q", got.LastUsableAddress, tt.wantLast)
			}
			if tt.wantBroadcast != "" && got.BroadcastAddress != tt.wantBroadcast {
				t.Errorf("BroadcastAddress = %q, want %q", got.BroadcastAddress, tt.wantBroadcast)
			}
			if got.HostOnly != tt.wantHostOnly {
				t.Errorf("HostOnly = %v, want %v", got.HostOnly, tt.wantHostOnly)
			}
		})
	}
}

func TestGetRemoteHost(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		xff        string
		wantIP     string
	}{
		{
			name:       "direct client",
			remoteAddr: "203.0.113.1:12345",
			wantIP:     "203.0.113.1",
		},
		{
			name:       "loopback with single XFF",
			remoteAddr: "127.0.0.1:12345",
			xff:        "198.51.100.1",
			wantIP:     "198.51.100.1",
		},
		{
			name:       "loopback with XFF chain",
			remoteAddr: "127.0.0.1:12345",
			xff:        "198.51.100.1, 10.0.0.1, 172.16.0.1",
			wantIP:     "198.51.100.1",
		},
		{
			name:       "loopback without XFF",
			remoteAddr: "127.0.0.1:12345",
			wantIP:     "127.0.0.1",
		},
		{
			name:       "IPv6 loopback with XFF",
			remoteAddr: "[::1]:12345",
			xff:        "203.0.113.50",
			wantIP:     "203.0.113.50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest("GET", "/", nil)
			r.RemoteAddr = tt.remoteAddr
			if tt.xff != "" {
				r.Header.Set("X-Forwarded-For", tt.xff)
			}

			got := GetRemoteHost(r)
			if got.String() != tt.wantIP {
				t.Errorf("GetRemoteHost() = %q, want %q", got.String(), tt.wantIP)
			}
		})
	}
}
