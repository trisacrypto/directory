package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/trisacrypto/directory/pkg/gds/secrets"
	"github.com/trisacrypto/directory/pkg/models/v1"
	storerr "github.com/trisacrypto/directory/pkg/store/errors"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/urfave/cli/v2"
)

//===========================================================================
// Contact Functions
//===========================================================================

func migrateContacts(c *cli.Context) (err error) {
	dryrun := c.Bool("dryrun")
	compare := c.Bool("compare")

	// Iterate through all vasps in the database
	vasps := db.ListVASPs(context.Background())
	defer vasps.Release()
	for vasps.Next() {
		var vasp *pb.VASP
		if vasp, err = vasps.VASP(); err != nil {
			return cli.Exit(err, 1)
		}

		// Iterate through all contacts on the vasp
		contacts := models.NewContactIterator(vasp.Contacts, models.SkipNoEmail())
		for contacts.Next() {
			vaspContact, _ := contacts.Value()

			var contact *models.Contact
			if contact, err = db.RetrieveContact(context.Background(), vaspContact.Email); err != nil {
				if errors.Is(err, storerr.ErrEntityNotFound) {
					if dryrun {
						fmt.Printf("contact %s missing and needs to be migrated\n\n", vaspContact.Email)
					} else {
						// Create the contact if it doesn't exist
						extra := &models.GDSContactExtraData{}
						if err := vaspContact.Extra.UnmarshalTo(extra); err != nil {
							return cli.Exit(fmt.Errorf("could not unmarshal extra for %s: %s", vaspContact.Email, err), 1)
						}

						contact = &models.Contact{
							Email:    vaspContact.Email,
							Name:     vaspContact.Name,
							Vasps:    []string{vasp.CommonName},
							Verified: extra.Verified,
							Token:    extra.Token,
							EmailLog: extra.EmailLog,
							Created:  time.Now().Format(time.RFC3339),
							Modified: time.Now().Format(time.RFC3339),
						}

						if _, err := db.CreateContact(context.Background(), contact); err != nil {
							return cli.Exit(err, 1)
						}
						fmt.Printf("contact %s created!\n", vaspContact.Email)
					}
					continue
				}
				return cli.Exit(err, 1)
			}

			if !compare {
				continue
			}

			// Check contact equality
			updates := make(map[string]map[string]interface{})
			if vaspContact.Email != contact.Email {
				updates["email"] = map[string]interface{}{
					"contact": contact.Email,
					"vasp":    vaspContact.Email,
				}
			}

			if vaspContact.Name != contact.Name {
				updates["name"] = map[string]interface{}{
					"contact": contact.Name,
					"vasp":    vaspContact.Name,
				}
			}

			found := false
			for _, included := range contact.Vasps {
				if vasp.CommonName == included {
					found = true
					break
				}
			}

			if !found {
				updates["vasps"] = map[string]interface{}{
					"contact": contact.Vasps,
					"vasp":    append(contact.Vasps, vasp.CommonName),
				}
			}

			extra := &models.GDSContactExtraData{}
			if err := vaspContact.Extra.UnmarshalTo(extra); err != nil {
				return cli.Exit(fmt.Errorf("could not unmarshal extra for %s: %s", vaspContact.Email, err), 1)
			}

			if extra.Verified != contact.Verified {
				updates["verified"] = map[string]interface{}{
					"contact": contact.Verified,
					"vasp":    extra.Verified,
				}
			}

			if extra.Token != contact.Token {
				updates["token"] = map[string]interface{}{
					"contact": contact.Token,
					"vasp":    extra.Token,
				}
			}

			if !slices.Equal(extra.EmailLog, contact.EmailLog) {
				updates["email_log"] = map[string]interface{}{
					"contact": len(contact.EmailLog),
					"vasp":    len(extra.EmailLog),
				}
			}

			if len(updates) > 0 {
				if dryrun {
					fmt.Printf("contact %s does not match vasp contact and needs to be updated\n", contact.Email)
					w := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', tabwriter.AlignRight)
					fmt.Fprintln(w, "field\tcontact\tvasp")
					for field, vals := range updates {
						fmt.Fprintf(w, "%s\t%v\t%v\n", field, vals["contact"], vals["vasp"])
					}
					w.Flush()
					fmt.Println()
				} else {
					fmt.Printf("contact %s updated!\n", contact.Email)
				}
			}
		}
	}

	if err = vasps.Error(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func fixVerifyToken(c *cli.Context) (err error) {
	dryrun := c.Bool("dryrun")

	vasps := db.ListVASPs(context.Background())
	defer vasps.Release()
	for vasps.Next() {
		var vasp *pb.VASP
		if vasp, err = vasps.VASP(); err != nil {
			return cli.Exit(err, 1)
		}

		contacts := models.NewContactIterator(vasp.Contacts, models.SkipNoEmail(), models.SkipDuplicates())
		for contacts.Next() {
			vaspContact, kind := contacts.Value()

			var contact *models.Contact
			if contact, err = db.RetrieveContact(context.Background(), vaspContact.Email); err != nil {
				return cli.Exit(err, 1)
			}

			if !contact.Verified && contact.Token == "" {
				fmt.Printf("contact %s is not verified and has no verification token\n", contact.Email)

				if !dryrun {
					contact.Token = secrets.CreateToken(models.VerificationTokenLength)
					if err = db.UpdateContact(context.Background(), contact); err != nil {
						return cli.Exit(err, 1)
					}
				}
			}

			// Ensure the vaspContact matches the contact
			token, verified, err := models.GetContactVerification(vaspContact)
			if err != nil {
				return cli.Exit(err, 1)
			}

			if contact.Verified != verified || contact.Token != token {
				vaspName, _ := vasp.Name()
				fmt.Printf("vasp %s contact %s (%s) does not match contact record\n", vaspName, kind, vaspContact.Email)

				if !dryrun {
					if err = models.SetContactVerification(vaspContact, contact.Token, contact.Verified); err != nil {
						return cli.Exit(err, 1)
					}

					if err = db.UpdateVASP(context.Background(), vasp); err != nil {
						return cli.Exit(err, 1)
					}
				}
			}

		}
	}

	if err = vasps.Error(); err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}

func contactList(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	contacts := db.ListContacts(ctx)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Name\tEmail\tVerified\tVASP(s)")

	for _, contact := range contacts {
		row := []string{
			contact.Name,
			contact.Email,
			fmt.Sprintf("%t", contact.Verified),
			strings.Join(contact.Vasps, ", "),
		}
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	w.Flush()
	return nil
}

func contactExport(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	vasps := db.ListVASPs(ctx)

	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	w.Write([]string{"ID", "VASP", "Administrative", "Technical", "Legal", "Billing"})

	for vasps.Next() {
		vasp, err := vasps.VASP()
		if err != nil {
			continue
		}

		row := make([]string, 6)
		row[0] = vasp.Id
		row[1], _ = vasp.Name()

		contacts := vasp.Contacts

		if contacts.Administrative != nil {
			row[2] = fmt.Sprintf("%q <%s>", contacts.Administrative.Name, contacts.Administrative.Email)
		}

		if contacts.Technical != nil {
			row[3] = fmt.Sprintf("%q <%s>", contacts.Technical.Name, contacts.Technical.Email)
		}

		if contacts.Legal != nil {
			row[4] = fmt.Sprintf("%q <%s>", contacts.Legal.Name, contacts.Legal.Email)
		}

		if contacts.Billing != nil {
			row[5] = fmt.Sprintf("%q <%s>", contacts.Billing.Name, contacts.Billing.Email)
		}

		w.Write(row)
	}

	return nil
}

func contactDetail(c *cli.Context) (err error) {
	if c.NArg() == 0 {
		return cli.Exit("specify at least one email address", 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if c.NArg() == 1 {
		var contact *models.Contact
		if contact, err = db.RetrieveContact(ctx, c.Args().First()); err != nil {
			return cli.Exit(err, 1)
		}
		return printJSON(contact)
	}

	contacts := make([]*models.Contact, 0, c.NArg())
	for i := 0; i < c.NArg(); i++ {
		var contact *models.Contact
		if contact, err = db.RetrieveContact(ctx, c.Args().Get(i)); err != nil {
			return cli.Exit(err, 1)
		}

		contacts = append(contacts, contact)
	}
	return printJSON(contacts)
}
