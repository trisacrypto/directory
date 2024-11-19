package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/trisacrypto/directory/pkg/models/v1"
	"github.com/trisacrypto/directory/pkg/utils"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/urfave/cli/v2"
)

//===========================================================================
// VASP Functions
//===========================================================================

func vaspList(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	iter := db.ListVASPs(ctx)
	defer iter.Release()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Name\tID\tCommon Name")

	for iter.Next() {
		vasp, _ := iter.VASP()
		name, _ := vasp.Name()
		fmt.Fprintln(w, strings.Join([]string{name, vasp.Id, vasp.CommonName}, "\t"))
	}

	if err = iter.Error(); err != nil {
		return cli.Exit(err, 1)
	}

	w.Flush()
	return nil
}

func vaspDetail(c *cli.Context) (err error) {
	if c.NArg() == 0 {
		return cli.Exit("specify at least one vasp uuid to retrieve", 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if c.NArg() == 1 {
		var vasp *pb.VASP
		if vasp, err = db.RetrieveVASP(ctx, c.Args().First()); err != nil {
			return cli.Exit(err, 1)
		}

		return printJSON(vasp)
	}

	vasps := make([]*pb.VASP, 0, c.NArg())
	for i := 0; i < c.NArg(); i++ {
		var vasp *pb.VASP
		if vasp, err = db.RetrieveVASP(ctx, c.Args().Get(i)); err != nil {
			return cli.Exit(err, 1)
		}
		vasps = append(vasps, vasp)
	}
	return printJSON(vasps)
}

func rereview(c *cli.Context) (err error) {
	var newStatus pb.VerificationState
	if state := c.String("verification-state"); state != "" {
		state = strings.Replace(strings.ToUpper(state), " ", "_", -1)
		newStatus = pb.VerificationState(pb.VerificationState_value[state])
	}

	if newStatus == pb.VerificationState_NO_VERIFICATION {
		return cli.Exit("verification-state needs to be specified", 1)
	}

	vaspID := c.String("vasp")
	fmt.Printf("lookup vasp with id %s\n", vaspID)

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	var vasp *pb.VASP
	if vasp, err = db.RetrieveVASP(ctx, vaspID); err != nil {
		return cli.Exit(fmt.Errorf("could not find VASP record: %s", err), 1)
	}

	if vasp.VerificationStatus == pb.VerificationState_VERIFIED {
		return cli.Exit(fmt.Errorf("VASP is %q -- use revoke instead", vasp.VerificationStatus), 1)
	}

	// Check with the user if we should continue with the certificate revocation
	fmt.Printf("updating verification state for %s\n", vasp.CommonName)
	if !c.Bool("yes") {
		if !askForConfirmation(fmt.Sprintf("set VASP status to %q?", newStatus)) {
			return cli.Exit(fmt.Errorf("canceled by user"), 1)
		}
	}

	if err = models.UpdateVerificationStatus(
		vasp,
		newStatus,
		"verification state updated by admins",
		"support@rotational.io",
	); err != nil {
		return cli.Exit(fmt.Errorf("could not update VASP status: %s", err), 1)
	}

	if newStatus < pb.VerificationState_VERIFIED {
		vasp.VerifiedOn = ""
	}

	if newStatus == pb.VerificationState_VERIFIED {
		vasp.VerifiedOn = time.Now().Format(time.RFC3339Nano)
	}

	if err = db.UpdateVASP(ctx, vasp); err != nil {
		return cli.Exit(fmt.Errorf("could not save VASP: %s", err), 1)
	}

	fmt.Println("VASP verification state updated")
	return nil
}

func destroy(c *cli.Context) (err error) {
	vaspID := c.String("vasp")
	fmt.Printf("lookup vasp with id %s\n", vaspID)

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	var vasp *pb.VASP
	if vasp, err = db.RetrieveVASP(ctx, vaspID); err != nil {
		return cli.Exit(fmt.Errorf("could not find VASP record: %s", err), 1)
	}

	if vasp.VerificationStatus == pb.VerificationState_VERIFIED {
		return cli.Exit(fmt.Errorf("VASP is %q -- use revoke before destroying VASP record", vasp.VerificationStatus), 1)
	}

	fmt.Printf("destroying VASP record for %s\n", vasp.CommonName)
	if !c.Bool("yes") {
		if !askForConfirmation("continue with operation?") {
			return cli.Exit(fmt.Errorf("canceled by user"), 1)
		}
	}

	if err = db.DeleteVASP(ctx, vaspID); err != nil {
		return cli.Exit(fmt.Errorf("could not delete record: %s", err), 1)
	}
	return nil
}

func vaspStatus(c *cli.Context) (err error) {
	vaspID := c.String("vasp")
	fmt.Printf("lookup vasp with id %s\n", vaspID)

	ctx, cancel := utils.WithDeadline(context.Background())
	defer cancel()

	var vasp *pb.VASP
	if vasp, err = db.RetrieveVASP(ctx, vaspID); err != nil {
		return cli.Exit(fmt.Errorf("could not find VASP record: %s", err), 1)
	}

	name, _ := vasp.Name()
	fmt.Printf("Name: %s\nCommon Name: %s\nStatus: %s\n\n", name, vasp.CommonName, vasp.VerificationStatus)

	certreqs, err := models.GetCertReqIDs(vasp)
	if err != nil {
		return cli.Exit(err, 1)
	}

	for i, certreq := range certreqs {
		ca, err := db.RetrieveCertReq(ctx, certreq)
		if err != nil {
			return cli.Exit(err, 1)
		}

		fmt.Printf("Certificate Request %d:\n  ID: %s\n  Common Name: %s\n  Status: %s\n  SANs: %s\n\n", i+1, ca.Id, ca.CommonName, ca.Status, strings.Join(ca.DnsNames, ", "))

	}
	return nil
}
