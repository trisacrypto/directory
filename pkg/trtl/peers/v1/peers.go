package peers

import "fmt"

func (p *Peer) Key() string {
	if p.Id == 0 {
		return ""
	}
	return fmt.Sprintf("%04x", p.Id)
}
