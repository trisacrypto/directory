package models

import (
	"time"

	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/types/known/anypb"
)

// HealthCheckExtra represents "extra" health check data.
type HealthCheckExtra struct {
	CheckAfter  string
	CheckBefore string
	Attempts    int32
	LastChecked string
}

// DelayCheck returns true if the health check
// should be skipped.
func (h *HealthCheckExtra) DelayCheck() bool {
	ca, err := time.Parse(
		time.RFC3339,
		h.CheckAfter)
	if err != nil {
		return false
	}
	return ca.After(time.Now())
}

// GetHealthCheckInfo from the extra data on the VASP record.
func GetHealthCheckInfo(vasp *pb.VASP) (*HealthCheckExtra, error) {
	// If the extra data is nil, return empty string with no error
	if vasp.Extra == nil {
		return nil, nil
	}

	// Unmarshal the extra data field on the VASP
	extra := &GDSExtraData{}
	if err := vasp.Extra.UnmarshalTo(extra); err != nil {
		return nil, nil
	}

	return &HealthCheckExtra{
		CheckAfter:  extra.GetHealthCheckAfter(),
		CheckBefore: extra.GetHealthCheckBefore(),
		Attempts:    extra.GetHealthCheckAttempts(),
		LastChecked: extra.GetHealthCheckLastChecked(),
	}, nil
}

// SetHealthCheckInfo on the extra data on the VASP record.
func SetHealthCheckInfo(vasp *pb.VASP, healthCheck HealthCheckExtra) (err error) {
	// maintains any other fields already in extra
	extra := &GDSExtraData{}
	if err = vasp.Extra.UnmarshalTo(extra); err != nil {
		return err
	}

	extra.HealthCheckAfter = healthCheck.CheckAfter
	extra.HealthCheckBefore = healthCheck.CheckBefore
	extra.HealthCheckAttempts = healthCheck.Attempts
	extra.HealthCheckLastChecked = healthCheck.LastChecked
	if vasp.Extra, err = anypb.New(extra); err != nil {
		return err
	}
	return nil
}
