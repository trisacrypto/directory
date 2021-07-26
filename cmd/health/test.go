package main

import (
	"testing"
	"time"

	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/store/mockdb"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

func TestRun(t *testing.T) {
	db := mockdb.MockDB{}

	db.OnRetrieveAll = func(opts *models.RetrieveAllOpts, c chan *pb.VASP) error {
		if opts.VerificationStatus != nil && *opts.VerificationStatus != pb.VerificationState_VERIFIED {
			t.Fatalf("unexpected verification status option: %s", opts.VerificationStatus)
		} else if !opts.TrisaEndpointExists {
			t.Fatal("expected option to be true")
		}
		c <- &pb.VASP{Id: "one", TrisaEndpoint: "/trisa.gds.api.v1beta1.TRISADirectory"}
		return nil
	}

	db.OnUpdateStatus = func(id string, status int32) error {
		if id != "one" {
			t.Fatalf("unexpected id: %s", id)
		} else if status != int32(pb.ServiceState_HEALTHY) {
			t.Fatalf("unexpected status: %d", status)
		}
		return nil
	}

	run(db, 1*time.Minute)
}
