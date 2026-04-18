package snapshot

// Diff holds the changes between two snapshots.
type Diff struct {
	Opened []PortInfo
	Closed []PortInfo
}

// HasChanges returns true if there are any opened or closed ports.
func (d *Diff) HasChanges() bool {
	return len(d.Opened) > 0 || len(d.Closed) > 0
}

// Compare computes the difference between an old and a new snapshot.
func Compare(old, new *Snapshot) *Diff {
	oldMap := indexPorts(old.Ports)
	newMap := indexPorts(new.Ports)

	diff := &Diff{}

	for key, info := range newMap {
		if _, exists := oldMap[key]; !exists {
			diff.Opened = append(diff.Opened, info)
		}
	}

	for key, info := range oldMap {
		if _, exists := newMap[key]; !exists {
			diff.Closed = append(diff.Closed, info)
		}
	}

	return diff
}

func indexPorts(ports []PortInfo) map[string]PortInfo {
	m := make(map[string]PortInfo, len(ports))
	for _, p := range ports {
		key := p.Proto + ":" + itoa(p.Port)
		m[key] = p
	}
	return m
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}
