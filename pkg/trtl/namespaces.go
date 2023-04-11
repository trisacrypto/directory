package trtl

import "github.com/trisacrypto/directory/pkg/utils/wire"

// TODO: need functionality to actually extract namespaces from db
// ???
const (
	NamespacePeers         = wire.NamespaceReplicas
	NamespaceIndex         = wire.NamespaceIndices
	NamespaceDefault       = "default"
	NamespaceUnknown       = "unknown"
	NamespaceSequence      = wire.NamespaceSequence
	NamespaceVASPs         = wire.NamespaceVASPs
	NamespaceCertReqs      = wire.NamespaceCertReqs
	NamespaceCerts         = wire.NamespaceCerts
	NamespaceAnnouncements = wire.NamespaceAnnouncements
	NamespaceOrganizations = wire.NamespaceOrganizations
)

// Reserved namespaces that cannot be used by the caller since they are in use by trtl.
// If necessary this can be moved to configuration in the future.
var reservedNamespaces = map[string]struct{}{
	NamespacePeers:    {},
	NamespaceSequence: {},
	NamespaceDefault:  {}, // if the user does not specify a namespace
	NamespaceUnknown:  {}, // the "unknown" namespace for debugging

	// TODO: add index namespace back to reserved namespaces when trtl does indexing.
	// NamespaceIndex:    {},
}

// Replicated namespaces are the namespaces that are used in anti-entropy by default.
var replicatedNamespaces = []string{
	NamespaceVASPs,
	NamespaceCertReqs,
	NamespaceCerts,
	NamespaceAnnouncements,
	NamespaceOrganizations,
	NamespacePeers,
	NamespaceDefault,
}

// Measured namespaces are the namespaces that are measured by the monitor.
var measuredNamespaces = []string{
	NamespacePeers,
	NamespaceIndex,
	NamespaceDefault,
	NamespaceSequence,
	NamespaceVASPs,
	NamespaceCertReqs,
	NamespaceCerts,
	NamespaceAnnouncements,
	NamespaceOrganizations,
}
