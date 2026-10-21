package traffic

import "testing"

func TestNormalizeGoScheme(t *testing.T) {
	cases := map[string]string{
		"socks5h": "socks5",
		"socks5":  "socks5",
		"http":    "http",
		"https":   "https",
		"":        "",
		"weird":   "weird",
	}
	for in, want := range cases {
		if got := NormalizeGoScheme(in); got != want {
			t.Errorf("NormalizeGoScheme(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestValidateProxyURLAcceptsSocks5h(t *testing.T) {
	u, err := ValidateProxyURL("socks5h://192.168.3.140:9050")
	if err != nil {
		t.Fatalf("socks5h proxy rejected: %v", err)
	}
	if u.Scheme != "socks5h" {
		t.Errorf("scheme = %q, want socks5h preserved for CLI consumers", u.Scheme)
	}
	if _, err := ValidateProxyURL("socks6h://h"); err == nil {
		t.Errorf("unknown scheme must still be rejected")
	}
}
