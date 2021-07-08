package models

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestHealthCheckExtraDelayCheck(t *testing.T) {
	now := time.Now()

	// check after in the past
	hc := HealthCheckExtra{
		CheckAfter: now.Add(-5 * time.Minute).Format(time.RFC3339),
	}
	if hc.DelayCheck() {
		t.Fatal("expected check after time in the past to return false")
	}

	// check after in the future
	hc.CheckAfter = now.Add(1 * time.Hour).Format(time.RFC3339)
	if !hc.DelayCheck() {
		t.Fatal("expected check after time in the future to return true")
	}
}

func TestGetHealthCheckInfo(t *testing.T) {
	now := time.Now()
	checkAfter := now.Add(-5 * time.Minute).Format(time.RFC3339)
	checkBefore := now.Add(2 * time.Hour).Format(time.RFC3339)
	lastChecked := now.Add(-5 * time.Hour).Format(time.RFC3339)

	extra := &GDSExtraData{
		AdminVerificationToken: "token", // ensure that this is not overwritten
		HealthCheckAfter:       checkAfter,
		HealthCheckBefore:      checkBefore,
		HealthCheckAttempts:    2,
		HealthCheckLastChecked: lastChecked,
	}
	vasp := &pb.VASP{}
	var err error
	if vasp.Extra, err = anypb.New(extra); err != nil {
		t.Fatal(err)
	}

	if hce, err := GetHealthCheckInfo(vasp); err != nil {
		t.Fatal(err)
	} else if diff := cmp.Diff(hce, &HealthCheckExtra{
		CheckAfter:  checkAfter,
		CheckBefore: checkBefore,
		Attempts:    2,
		LastChecked: lastChecked,
	}); diff != "" {
		t.Fatal(diff)
	}

	updated := HealthCheckExtra{
		CheckAfter:  now.Add(5 * time.Hour).Format(time.RFC3339),
		CheckBefore: now.Add(7 * time.Hour).Format(time.RFC3339),
		Attempts:    3,
		LastChecked: now.Format(time.RFC3339),
	}
	if err := SetHealthCheckInfo(vasp, updated); err != nil {
		t.Fatal(err)
	}

	actualExtra := &GDSExtraData{}
	if err = vasp.Extra.UnmarshalTo(actualExtra); err != nil {
		t.Fatal(err)
	} else if actualExtra.AdminVerificationToken != "token" {
		t.Fatalf("unexpected token: %s", actualExtra.AdminVerificationToken)
	} else if actualExtra.HealthCheckAfter != updated.CheckAfter {
		t.Fatalf("unexpected check after: %s", actualExtra.HealthCheckAfter)
	} else if actualExtra.HealthCheckBefore != updated.CheckBefore {
		t.Fatalf("unexpected check before: %s", actualExtra.HealthCheckBefore)
	} else if actualExtra.HealthCheckAttempts != 3 {
		t.Fatalf("unexpected attempts: %d", actualExtra.HealthCheckAttempts)
	} else if actualExtra.HealthCheckLastChecked != updated.LastChecked {
		t.Fatalf("unexpected last checked: %s", actualExtra.HealthCheckLastChecked)
	}
}
