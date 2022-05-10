package main

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/secrets"
	"github.com/trisacrypto/directory/pkg/gds/store"
	"github.com/trisacrypto/directory/pkg/utils/logger"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/urfave/cli/v2"
)

func interactiveReissue(c *cli.Context) (err error) {
	var reissuer *InteractiveReissue
	if reissuer, err = NewInteractiveReissue(c); err != nil {
		return cli.Exit(fmt.Errorf("could not create reissuer: %s", err), 1)
	}

	if err = reissuer.Run(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func NewInteractiveReissue(c *cli.Context) (reissuer *InteractiveReissue, err error) {
	// Turn off logging to make it easier to read interactive ui
	logger.Discard()

	// Load the configuration from the environment and connect to the database store
	reissuer = &InteractiveReissue{c: c}
	if reissuer.conf, err = config.New(); err != nil {
		return nil, err
	}

	reissuer.conf.Database.ReindexOnBoot = false
	if reissuer.db, err = store.Open(reissuer.conf.Database); err != nil {
		return nil, err
	}
	return reissuer, nil
}

type InteractiveReissue struct {
	c    *cli.Context
	db   store.Store
	conf config.Config
}

func (r *InteractiveReissue) Run() (err error) {
	defer r.db.Close()

	// Step zero: get email address for audit log
	var email string
	if email, err = r.loadEmail(); err != nil {
		return err
	}

	// Step one: load the VASP
	var vasp *pb.VASP
	if vasp, err = r.loadVASP(); err != nil {
		return fmt.Errorf("could not find VASP: %s", err)
	}

	name, _ := vasp.Name()
	fmt.Printf("Editing VASP %q which is %s\n", name, vasp.VerificationStatus)

	// Step two: edit the VASP if necessary
	if err = r.editVASP(vasp, email); err != nil {
		return fmt.Errorf("could not update VASP: %s", err)
	}

	// Step three: get current certificate requests to cancel any that haven't been downloaded
	var crids []string
	if crids, err = r.getCertReqIDs(vasp); err != nil {
		return fmt.Errorf("could not get certificate request IDs: %s", err)
	}

	for _, crid := range crids {
		if err = r.handleOldCertReq(vasp, crid, email); err != nil {
			return fmt.Errorf("could not handle old certificate request: %s", err)
		}
	}

	// Step four: create a new certificate request to reissue the certs
	fmt.Printf("Creating new certificate ... ")
	var certreq *models.CertificateRequest
	if certreq, err = models.NewCertificateRequest(vasp); err != nil {
		return fmt.Errorf("could not create a new certificate request: %s", err)
	}

	if err = models.UpdateCertificateRequestStatus(certreq, models.CertificateRequestState_READY_TO_SUBMIT, "manually reissuing certificates", email); err != nil {
		return fmt.Errorf("could not update certificate request status: %s", err)
	}

	// TODO: add dNSName param here for extra SANs
	fmt.Printf("%s\n", certreq.Id)

	// Create a secret password for the new certs
	if err = r.createPKCS12Password(certreq.Id); err != nil {
		return fmt.Errorf("could not create pkcs12 password: %s", err)
	}

	if err = r.db.UpdateCertReq(certreq); err != nil {
		return fmt.Errorf("could not save new certificate request: %s", err)
	}

	// Step six: update audit log with new certificate reissuance and append cert request to VASP
	if err = models.UpdateVerificationStatus(vasp, vasp.VerificationStatus, "reissuing certificates", email); err != nil {
		return fmt.Errorf("could not update audit log for VASP: %s", err)
	}

	if err = models.AppendCertReqID(vasp, certreq.Id); err != nil {
		return fmt.Errorf("could not update certreqs for VASP: %s", err)
	}

	if err = r.db.UpdateVASP(vasp); err != nil {
		return fmt.Errorf("could not update VASP: %s", err)
	}

	// Job done!
	fmt.Println("Certificate request for reissued certificates completed")
	return nil
}

func (r *InteractiveReissue) loadEmail() (_ string, err error) {
	prompt := promptui.Prompt{
		Label: "Enter Your Email (for audit log)",
		Validate: func(input string) error {
			if input == "" {
				return errors.New("email address is required")
			}
			if _, err := mail.ParseAddress(input); err != nil {
				return err
			}
			return nil
		},
	}
	return prompt.Run()
}

func (r *InteractiveReissue) loadVASP() (vasp *pb.VASP, err error) {
	prompt := promptui.Prompt{
		Label: "Enter VASP ID",
		Validate: func(input string) error {
			if _, err := uuid.Parse(input); err != nil {
				return err
			}
			return nil
		},
	}

	var vaspID string
	if vaspID, err = prompt.Run(); err != nil {
		return nil, err
	}

	if vasp, err = r.db.RetrieveVASP(vaspID); err != nil {
		return nil, err
	}

	return vasp, nil
}

func (r *InteractiveReissue) editVASP(vasp *pb.VASP, email string) (err error) {
	if !confirm("Do you want to edit the common name?") {
		return nil
	}

	cnp := promptui.Prompt{
		Label:     "Common Name:",
		Default:   vasp.CommonName,
		AllowEdit: true,
		Validate: func(input string) error {
			re := regexp.MustCompile(`^[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,7}$`)
			if !re.MatchString(input) {
				return errors.New("enter a valid common name")
			}
			return nil
		},
	}

	if vasp.CommonName, err = cnp.Run(); err != nil {
		return err
	}

	epp := promptui.Prompt{
		Label:     "Endpoint:",
		Default:   vasp.TrisaEndpoint,
		AllowEdit: true,
		Validate: func(input string) error {
			re := regexp.MustCompile(`^[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,7}:\d+$`)
			if !re.MatchString(input) {
				return errors.New("enter a valid common name")
			}
			return nil
		},
	}

	if vasp.TrisaEndpoint, err = epp.Run(); err != nil {
		return err
	}

	if err = models.UpdateVerificationStatus(vasp, vasp.VerificationStatus, "changed common name and/or endpoint to reissue certificates", email); err != nil {
		return err
	}

	if err = r.db.UpdateVASP(vasp); err != nil {
		return err
	}
	return nil
}

func (r *InteractiveReissue) getCertReqIDs(vasp *pb.VASP) (ids []string, err error) {
	if ids, err = models.GetCertReqIDs(vasp); err != nil {
		return nil, err
	}

	if len(ids) > 0 {
		return ids, nil
	}

	// If the IDs are not stored on the VASP, do a manual search just to make sure and
	// update the VASP with any discovered cert req IDs. This should rarely happen and
	// only for old records.
	fmt.Println("VASP does not have any associated certificate requests")
	if !confirm("Search for unlinked certificate requests and update VASP?") {
		return nil, nil
	}

	ids = make([]string, 0)
	iter := r.db.ListCertReqs()
	defer iter.Release()
	for iter.Next() {
		var cr *models.CertificateRequest
		if cr, err = iter.CertReq(); err != nil {
			continue
		}

		if cr.Vasp == vasp.Id {
			ids = append(ids, cr.Id)
		}
	}

	if len(ids) > 0 {
		fmt.Printf("Attempting to append %d ids to the VASP ... ", len(ids))
		nappends := 0
		for _, id := range ids {
			if err = models.AppendCertReqID(vasp, id); err == nil {
				nappends++
			}
		}
		if err = r.db.UpdateVASP(vasp); err != nil {
			fmt.Printf("could not update VASP with certificate request IDs: %s\n", err)
		} else {
			fmt.Printf("appended %d certificate request IDs to VASP\n", nappends)
		}
	}
	return ids, nil
}

func (r *InteractiveReissue) handleOldCertReq(vasp *pb.VASP, crid, email string) (err error) {
	var certreq *models.CertificateRequest
	if certreq, err = r.db.RetrieveCertReq(crid); err != nil {
		return err
	}

	// Check the current certreq status; if it hasn't already been downloaded, then cancel it.
	if certreq.Status < models.CertificateRequestState_COMPLETED {
		reason := "superceded by new certificate request"
		fmt.Printf("Canceling certificate request %s and setting state %s from %s\n", certreq.Id, models.CertificateRequestState_CR_ERRORED, certreq.Status)
		if err = models.UpdateCertificateRequestStatus(certreq, models.CertificateRequestState_CR_ERRORED, reason, email); err != nil {
			return err
		}
		certreq.RejectReason = reason
		certreq.Modified = time.Now().Format(time.RFC3339)

		if err = r.db.UpdateCertReq(certreq); err != nil {
			return err
		}
	} else {
		fmt.Printf("Certificate request %s is in state %s - making no changes\n", certreq.Id, certreq.Status)
	}

	return nil
}

func (r *InteractiveReissue) createPKCS12Password(crid string) (err error) {
	var sm *secrets.SecretManager
	if sm, err = secrets.New(r.conf.Secrets); err != nil {
		return err
	}

	password := secrets.CreateToken(16)

	secretType := "password"
	if err = sm.With(crid).CreateSecret(context.Background(), secretType); err != nil {
		return err
	}
	if err = sm.With(crid).AddSecretVersion(context.TODO(), secretType, []byte(password)); err != nil {
		return err
	}

	fmt.Printf("The PKCS12 password is: %s\n", password)
	return nil
}

func confirm(label string) bool {
	check := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	if _, err := check.Run(); err != nil {
		return false
	}
	return true
}
