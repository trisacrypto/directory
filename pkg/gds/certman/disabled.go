package certman

import (
	"sync"

	"github.com/rs/zerolog/log"
)

// Disabled implements the certman.Service interface but is essentially a no-op that
// warns that the certificate manager is disabled. This allows outsider users to
// interact with certman without having to check if it's enabled.
type Disabled struct{}

// Compile time interface implementation check.
var _ Service = &Disabled{}

func (d *Disabled) Run(*sync.WaitGroup) error {
	log.Warn().Msg("certman is disabled")
	return nil
}

func (d *Disabled) Stop() {
	log.Debug().Msg("stopping disabled certman")
}

func (d *Disabled) CertManager() {
	log.Trace().Msg("certman is disabled: cannot start cert manager go routine")
}

func (d *Disabled) HandleCertificateRequests() {
	log.Trace().Msg("certman is disabled: cannot handle certificate requests")
}

func (d *Disabled) HandleCertificateReissuance() {
	log.Trace().Msg("certman is disabled: cannot handle certificate reissuance")
}
