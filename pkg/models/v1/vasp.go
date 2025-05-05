package models

import (
	"github.com/trisacrypto/trisa/pkg/ivms101"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// Temporary struct to specify the VASP fields that we allow updating using JSON on
// the command line. This struct is only used for the gdsutil update command and can
// be removed if that command is removed.
type VASPCLIUpdate struct {
	Entity             *ivms101.LegalPerson   `json:"entity,omitempty"`
	Website            string                 `json:"website,omitempty"`
	BusinessCategory   string                 `json:"business_category,omitempty"`
	VASPCategories     []string               `json:"vasp_categories,omitempty"`
	EstablishedOn      string                 `json:"established_on,omitempty"`
	TRIXO              *pb.TRIXOQuestionnaire `json:"trixo,omitempty"`
	CertificateWebhook string                 `json:"certificate_webhook,omitempty"`
	NoEmailDelivery    *bool                  `json:"no_email_delivery,omitempty"`
}

func (v *VASPCLIUpdate) Update(vasp *pb.VASP) (err error) {
	if v.Entity != nil {
		vasp.Entity = v.Entity
	}
	if v.Website != "" {
		vasp.Website = v.Website
	}
	if v.BusinessCategory != "" {
		if vasp.BusinessCategory, err = pb.ParseBusinessCategory(v.BusinessCategory); err != nil {
			return err
		}
	}
	if len(v.VASPCategories) > 0 {
		for i, cat := range v.VASPCategories {
			if v.VASPCategories[i], err = pb.ValidVASPCategory(cat); err != nil {
				return err
			}
		}

		vasp.VaspCategories = v.VASPCategories
	}
	if v.EstablishedOn != "" {
		vasp.EstablishedOn = v.EstablishedOn
	}
	if v.TRIXO != nil {
		vasp.Trixo = v.TRIXO
	}
	if v.CertificateWebhook != "" {
		vasp.CertificateWebhook = v.CertificateWebhook
	}
	if v.NoEmailDelivery != nil {
		vasp.NoEmailDelivery = *v.NoEmailDelivery
	}

	return vasp.Validate(false)

}
