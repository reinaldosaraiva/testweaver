package server

import "testing"

func TestSetExecHostValidation(t *testing.T) {
	m, err := NewManager(t.TempDir(), "")
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	defer m.Close()

	cases := []struct {
		host string
		ok   bool
	}{
		{"macewindu@192.168.3.140", true},
		{"root@host.example.com", true},
		{"u@[2001:db8::1]", true},
		{"", true}, // empty = local execution, always legal
		{"192.168.3.140", false},
		{"u@h -oProxyCommand=evil", false},
		{"u@h;touch /tmp/pwn", false},
		{"u ser@h", false},
		{"u@h`id`", false},
		{"-oBatchMode=yes@h", false},
	}
	for _, c := range cases {
		err := m.SetExecHost(c.host, "/wd")
		if c.ok && err != nil {
			t.Errorf("SetExecHost(%q) = %v, want nil", c.host, err)
		}
		if !c.ok && err == nil {
			t.Errorf("SetExecHost(%q) = nil, want rejection", c.host)
		}
	}
}

func TestExecHostPersistsAcrossReopen(t *testing.T) {
	dir := t.TempDir()
	m1, err := NewManager(dir, "")
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	if err := m1.SetExecHost("macewindu@192.168.3.140", "/Users/macewindu/sw-exec"); err != nil {
		t.Fatalf("SetExecHost: %v", err)
	}
	if h, wd := m1.ExecHost(); h != "macewindu@192.168.3.140" || wd != "/Users/macewindu/sw-exec" {
		t.Fatalf("in-memory state = %q/%q", h, wd)
	}
	if err := m1.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	m2, err := NewManager(dir, "")
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer m2.Close()
	if h, wd := m2.ExecHost(); h != "macewindu@192.168.3.140" || wd != "/Users/macewindu/sw-exec" {
		t.Fatalf("persisted state = %q/%q, want the values set before close", h, wd)
	}

	// Clearing the host returns to local execution and persists as empty.
	if err := m2.SetExecHost("", ""); err != nil {
		t.Fatalf("clear: %v", err)
	}
	if h, _ := m2.ExecHost(); h != "" {
		t.Fatalf("after clear host = %q, want empty", h)
	}
}
