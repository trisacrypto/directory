package peers

import "fmt"

func (p *Peer) Key() string {
	return fmt.Sprintf("%04x", p.Id)
}
