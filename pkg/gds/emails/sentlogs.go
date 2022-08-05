// This contains helper functions to help confirm if emails were sent.
package emails

import (
	"time"

	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

// VerifyEmailSentVASPContacts checks if an email has been sent to the VASP contacts in some time window (e.g. the past 30 days).
// Returns false for either email has not been sent or an error occurred.
func VerifyEmailSentVASPContact(vasp *pb.VASP, kind string, timeWindowDays int) (status bool, err error) {

	contact := models.ContactFromType(vasp.Contacts, kind)
	emailLog, err := models.GetEmailLog(contact)
	if err != nil {
		status = false
		return status, err
	} else {
		// get latest timestamp
		strTimestamp := emailLog[len(emailLog)-1].Timestamp

		//format emailTimestamp
		timestamp, _ := time.Parse(time.RFC3339, strTimestamp)

		//check if an email has been send within time window.
		WithinTimeWindow := timestamp.After(time.Now().AddDate(0, 0, -timeWindowDays))
		if WithinTimeWindow {
			status = true
			return status, nil
		} else {
			status = false
			return status, nil
		}
	}
}
