package auth

import (
	"encoding/json"
	"sort"

	"github.com/trisacrypto/directory/pkg/bff/models/v1"
)

// AppMetadata makes it easier to serialize and deserialize JSON from the auth0
// app_metadata assigned to the user by the BFF (and ensures the data is structured).
type AppMetadata struct {
	OrgID         string   `json:"orgid"`
	VASPs         VASPs    `json:"vasps"`
	Organizations []string `json:"organizations"`
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

	for i, value := range meta.Organizations {
		if value != other.Organizations[i] {
			return false
		}
	}

	return true
}

func (meta *AppMetadata) GetOrganizations() []string {
	return meta.Organizations
}

func (meta *AppMetadata) Load(appdata *map[string]interface{}) (err error) {
	if appdata == nil {
		return nil
	}

	// Serialize appdata back to JSON
	var data []byte
	if data, err = json.Marshal(*appdata); err != nil {
		return err
	}

	// Deserialize app metadata from struct
	if err = json.Unmarshal(data, meta); err != nil {
		return err
	}

	return nil
}

func (meta *AppMetadata) Dump() (_ *map[string]interface{}, err error) {
	// Serialize meta back to JSON
	var data []byte
	if data, err = json.Marshal(meta); err != nil {
		return nil, err
	}

	appdata := make(map[string]interface{})
	if err = json.Unmarshal(data, &appdata); err != nil {
		return nil, err
	}

	return &appdata, nil
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
	meta.ClearOrganization()
	meta.OrgID = org.Id

	if org.Testnet != nil {
		meta.VASPs.TestNet = org.Testnet.Id
	}

	if org.Mainnet != nil {
		meta.VASPs.MainNet = org.Mainnet.Id
	}
}

// AddOrganization adds an organization ID to the set of organizations the user is a
// part of. This method is idempotent and will not add the organization ID if it
// already exists.
func (meta *AppMetadata) AddOrganization(orgID string) {
	// Find the index of insertion using a binary search
	i := indexOf(meta.Organizations, orgID)

	// If the organization ID is already in the list, do nothing
	if i < len(meta.Organizations) && meta.Organizations[i] == orgID {
		return
	}

	// Otherwise, insert the organization ID into the list
	meta.Organizations = append(meta.Organizations, "")
	copy(meta.Organizations[i+1:], meta.Organizations[i:])
	meta.Organizations[i] = orgID
}

// RemoveOrganization removes an organization ID from the set of organizations the user
// is a part of. This method is idempotent and will not error if the organization ID
// does not exist in the metadata.
func (meta *AppMetadata) RemoveOrganization(orgID string) {
	// Find the index of removal using a binary search
	i := indexOf(meta.Organizations, orgID)

	// If the organization ID is not in the list, do nothing
	if i >= len(meta.Organizations) || meta.Organizations[i] != orgID {
		return
	}

	// Otherwise, remove the organization ID from the list
	copy(meta.Organizations[i:], meta.Organizations[i+1:])
	meta.Organizations[len(meta.Organizations)-1] = ""
	meta.Organizations = meta.Organizations[:len(meta.Organizations)-1]
}

// indexOf uses a binary search to return the index where the target string should be
// inserted or found in the list. The list must already be sorted in ascending order,
// otherwise this method will have undefined behavior.
func indexOf(list []string, target string) int {
	return sort.Search(len(list), func(i int) bool { return list[i] >= target })
}
