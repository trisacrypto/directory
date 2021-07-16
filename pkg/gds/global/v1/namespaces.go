package global

// Namespace constants for all managed objects in GDS
const (
	NamespaceVASPs    = "vasps"
	NamespaceCertReqs = "certreqs"
	NamespaceReplicas = "peers"
	NamespaceIndices  = "index"
	NamespaceSequence = "sequence"
)

// Namespaces defines all possible namespaces that GDS manages
var Namespaces = [3]string{NamespaceVASPs, NamespaceCertReqs, NamespaceReplicas}
