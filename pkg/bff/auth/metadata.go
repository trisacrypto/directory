package auth

import (
	"encoding/json"

	"github.com/trisacrypto/directory/pkg/bff/models/v1"
)

// AppMetadata makes it easier to serialize and deserialize JSON from the auth0
// app_metadata assigned to the user by the BFF (and ensures the data is structured).
type AppMetadata struct {
	OrgID         string              `json:"orgid"`
	VASPs         VASPs               `json:"vasps"`
	Organizations map[string]struct{} `json:"organizations"`
}

type VASPs struct {
	MainNet string `json:"mainnet"`
	TestNet string `json:"testnet"`
}

// TODO: Hash-based method might be more maintainable, but this avoids error handling for now
func (meta *AppMetadata) Equals(other *AppMetadata) bool {
	if meta.OrgID != other.OrgID {
		return false
	}

	if meta.VASPs != other.VASPs {
		return false
	}

	if len(meta.Organizations) != len(other.Organizations) {
		return false
	}

	for i, org := range meta.Organizations {
		if value, ok := other.Organizations[i]; !ok || value != org {
			return false
		}
	}

	return true
}

func (meta *AppMetadata) Load(appdata map[string]interface{}) (err error) {
	// Serialize appdata back to JSON
	var data []byte
	if data, err = json.Marshal(appdata); err != nil {
		return err
	}

	// Deserialize app metadata from struct
	if err = json.Unmarshal(data, meta); err != nil {
		return err
	}

	return nil
}

func (meta *AppMetadata) Dump() (appdata map[string]interface{}, err error) {
	// Serialize meta back to JSON
	var data []byte
	if data, err = json.Marshal(meta); err != nil {
		return nil, err
	}

	appdata = make(map[string]interface{})
	if err = json.Unmarshal(data, &appdata); err != nil {
		return nil, err
	}

	return appdata, nil
}

// ClearOrganization removes all organization-related data from the app metadata.
func (meta *AppMetadata) ClearOrganization() {
	meta.OrgID = ""
	meta.VASPs.TestNet = ""
	meta.VASPs.MainNet = ""
}

// UpdateOrganization completely replaces the organization data in the app metadata
// with data from the organization record.
func (meta *AppMetadata) UpdateOrganization(org *models.Organization) {
	meta.OrgID = org.Id

	if org.Testnet != nil && org.Testnet.Id != "" {
		meta.VASPs.TestNet = org.Testnet.Id
	} else {
		meta.VASPs.TestNet = ""
	}

	if org.Mainnet != nil && org.Mainnet.Id != "" {
		meta.VASPs.MainNet = org.Mainnet.Id
	} else {
		meta.VASPs.MainNet = ""
	}
}

// AddOrganization adds an organization ID to the set of organizations the user is a
// part of. This method is idempotent and will not add the organization ID if it
// already exists.
func (meta *AppMetadata) AddOrganization(orgID string) {
	if meta.Organizations == nil {
		meta.Organizations = make(map[string]struct{})
	}
	meta.Organizations[orgID] = struct{}{}
}

// RemoveOrganization removes an organization ID from the set of organizations the user
// is a part of. This method is idempotent and will not error if the organization ID
// does not exist in the metadata.
func (meta *AppMetadata) RemoveOrganization(orgID string) {
	delete(meta.Organizations, orgID)
}
