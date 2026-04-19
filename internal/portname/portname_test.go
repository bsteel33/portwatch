package portname

import "testing"

func TestResolve_KnownPort(t *testing.T) {
	r := New()
	e, ok := r.Resolve(80, "tcp")
	if !ok {
		t.Fatal("expected port 80/tcp to be known")
	}
	if e.Name != "http" {
		t.Errorf("expected http, got %s", e.Name)
	}
	if e.Description == "" {
		t.Error("expected non-empty description")
	}
}

func TestResolve_UnknownPort(t *testing.T) {
	r := New()
	_, ok := r.Resolve(9999, "tcp")
	if ok {
		t.Error("expected port 9999/tcp to be unknown")
	}
}

func TestResolve_ProtoDistinct(t *testing.T) {
	r := New()
	// 53 exists for both tcp and udp
	for _, proto := range []string{"tcp", "udp"} {
		e, ok := r.Resolve(53, proto)
		if !ok {
			t.Errorf("expected 53/%s to be known", proto)
		}
		if e.Name != "dns" {
			t.Errorf("expected dns for 53/%s, got %s", proto, e.Name)
		}
	}
}

func TestName_Known(t *testing.T) {
	r := New()
	if got := r.Name(443, "tcp"); got != "https" {
		t.Errorf("expected https, got %s", got)
	}
}

func TestName_Unknown(t *testing.T) {
	r := New()
	if got := r.Name(12345, "tcp"); got != "12345/tcp" {
		t.Errorf("unexpected fallback: %s", got)
	}
}

func TestName_CaseInsensitiveProto(t *testing.T) {
	r := New()
	if got := r.Name(22, "TCP"); got != "ssh" {
		t.Errorf("expected ssh, got %s", got)
	}
}
