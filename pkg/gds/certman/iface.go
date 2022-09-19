package certman

import "sync"

// Service defines the CertMan go routine interface for outside users to interact with
// the certificate manager directly.
type Service interface {
	Run(*sync.WaitGroup) error
	Stop()
	CertManager()
	HandleCertificateRequests()
	HandleCertificateReissuance()
}
