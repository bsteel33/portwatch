package portgroup_test

import (
	"testing"

	"github.com/user/portwatch/internal/portgroup"
	"github.com/user/portwatch/internal/scanner"
)

func samplePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 22, Proto: "tcp", Service: "ssh"},
		{Port: 8080, Proto: "tcp", Service: "http-alt"},
		{Port: 51000, Proto: "tcp", Service: ""},
	}
}

func TestApply_GroupsWellKnown(t *testing.T) {
	g := portgroup.New()
	groups := g.Apply(samplePorts())
	for _, grp := range groups {
		if grp.Name == "well-known" {
			if len(grp.Ports) != 1 || grp.Ports[0].Port != 22 {
				t.Errorf("expected port 22 in well-known, got %v", grp.Ports)
			}
			return
		}
	}
	t.Error("well-known group not found")
}

func TestApply_GroupsRegistered(t *testing.T) {
	g := portgroup.New()
	groups := g.Apply(samplePorts())
	for _, grp := range groups {
		if grp.Name == "registered" {
			if len(grp.Ports) != 1 || grp.Ports[0].Port != 8080 {
				t.Errorf("expected port 8080 in registered, got %v", grp.Ports)
			}
			return
		}
	}
	t.Error("registered group not found")
}

func TestApply_GroupsDynamic(t *testing.T) {
	g := portgroup.New()
	groups := g.Apply(samplePorts())
	for _, grp := range groups {
		if grp.Name == "dynamic" {
			if len(grp.Ports) != 1 || grp.Ports[0].Port != 51000 {
				t.Errorf("expected port 51000 in dynamic, got %v", grp.Ports)
			}
			return
		}
	}
	t.Error("dynamic group not found")
}

func TestApply_EmptyPorts(t *testing.T) {
	g := portgroup.New()
	groups := g.Apply([]scanner.Port{})
	if len(groups) != 0 {
		t.Errorf("expected no groups, got %d", len(groups))
	}
}
