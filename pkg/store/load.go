package store

import (
	"context"
	"encoding/csv"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/trisacrypto/directory/pkg/utils"
	"github.com/trisacrypto/trisa/pkg/iso3166"
	"github.com/trisacrypto/trisa/pkg/ivms101"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// Load the database from a CSV fixture specified on disk by the path.
// TODO: this method needs to be updated for the new data structure.
func Load(db Store, path string) (err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return err
	}
	defer f.Close()

	rows := 0
	reader := csv.NewReader(f)

	for {
		var record []string
		if record, err = reader.Read(); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		rows++
		if rows == 1 {
			// Skip the expected header: entity name,country,category,url,address
			// TODO: validate the header
			continue
		}

		record[3] = strings.TrimSpace(record[3])
		if record[3] == "" {
			// Skip entries without websites
			continue
		}

		if !strings.HasPrefix(record[3], "http") {
			record[3] = "http://" + record[3]
		}

		// TODO: ensure that ID generation is correct
		vasp := &pb.VASP{
			Id:                  uuid.New().String(),
			RegisteredDirectory: "trisa.directory",
			Entity: &ivms101.LegalPerson{
				Name: &ivms101.LegalPersonName{
					NameIdentifiers: []*ivms101.LegalPersonNameId{
						{
							LegalPersonName:               record[0],
							LegalPersonNameIdentifierType: ivms101.LegalPersonLegal,
						},
					},
				},
				GeographicAddresses: []*ivms101.Address{
					{
						AddressLine: []string{
							record[4],
						},
					},
				},
			},
			Website:            record[3],
			VerificationStatus: pb.VerificationState_NO_VERIFICATION,
			ServiceStatus:      pb.ServiceState_UNKNOWN,
		}

		if record[1] != "#N/A" {

			var alphaCode iso3166.AlphaCode
			if alphaCode, err = iso3166.Find(record[1]); err != nil {
				return err
			}
			vasp.Entity.CountryOfRegistration = alphaCode.Alpha2
			vasp.Entity.GeographicAddresses[0].Country = alphaCode.Alpha2
		} else {
			vasp.Entity.CountryOfRegistration = "XX"
			vasp.Entity.GeographicAddresses[0].Country = "XX"
		}

		// TODO: better handling of VASP category in record[2]
		vasp.VaspCategories = []string{record[2]}

		var website *url.URL
		if website, err = url.Parse(record[3]); err != nil {
			return err
		}
		vasp.CommonName = website.Hostname()

		var id string
		ctx, cancel := utils.WithDeadline(context.Background())
		defer cancel()
		if id, err = db.CreateVASP(ctx, vasp); err != nil {
			return err
		}

		ctx, cancel = utils.WithDeadline(ctx)
		defer cancel()
		if _, err = db.RetrieveVASP(ctx, id); err != nil {
			return err
		}
	}
}
