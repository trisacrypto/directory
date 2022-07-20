package auth

import "encoding/json"

// AppMetadata makes it easier to serialize and deserialize JSON from the auth0
// app_metadata assigned to the user by the BFF (and ensures the data is structured).
type AppMetadata struct {
	OrgID string `json:"orgid"`
	VASPs VASPs  `json:"vasps"`
}

type VASPs struct {
	MainNet string `json:"mainnet"`
	TestNet string `json:"testnet"`
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
