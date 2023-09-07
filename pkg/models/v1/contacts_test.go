package models_test

import (
	"compress/gzip"
	"encoding/json"
	"math/bits"
	"net/mail"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/directory/pkg/gds/secrets"
	"github.com/trisacrypto/directory/pkg/models/v1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestIterContacts(t *testing.T) {
	contacts := &pb.Contacts{
		Technical: &pb.Contact{
			Name: "technical",
		},
		Administrative: &pb.Contact{
			Email: "administrative@example.com",
		},
		Billing: &pb.Contact{
			Name: "billing",
		},
		Legal: &pb.Contact{
			Email: "legal@example.com",
		},
	}
	expectedContacts := []*pb.Contact{
		contacts.Technical,
		contacts.Administrative,
		contacts.Legal,
		contacts.Billing,
	}
	expectedKinds := []string{
		models.TechnicalContact,
		models.AdministrativeContact,
		models.LegalContact,
		models.BillingContact,
	}

	actualContacts := []*pb.Contact{}
	actualKinds := []string{}

	// Should iterate over all contacts.
	iter := models.NewContactIterator(contacts)
	for iter.Next() {
		contact, kind := iter.Value()
		actualContacts = append(actualContacts, contact)
		actualKinds = append(actualKinds, kind)
	}
	require.NoError(t, iter.Error())
	require.Equal(t, expectedContacts, actualContacts)
	require.Equal(t, expectedKinds, actualKinds)

	actualContacts = []*pb.Contact{}
	actualKinds = []string{}

	// Should skip contacts without an email address.
	expectedContacts = []*pb.Contact{
		contacts.Administrative,
		contacts.Legal,
	}
	expectedKinds = []string{
		models.AdministrativeContact,
		models.LegalContact,
	}
	iter = models.NewContactIterator(contacts, models.SkipNoEmail())
	for iter.Next() {
		contact, kind := iter.Value()
		actualContacts = append(actualContacts, contact)
		actualKinds = append(actualKinds, kind)
	}
	require.NoError(t, iter.Error())
	require.Equal(t, expectedContacts, actualContacts)
	require.Equal(t, expectedKinds, actualKinds)

	actualContacts = []*pb.Contact{}
	actualKinds = []string{}

	// Should skip nil contacts.
	contacts.Technical = nil
	contacts.Billing = nil
	iter = models.NewContactIterator(contacts)
	for iter.Next() {
		contact, kind := iter.Value()
		actualContacts = append(actualContacts, contact)
		actualKinds = append(actualKinds, kind)
	}
	require.NoError(t, iter.Error())
	require.Equal(t, expectedContacts, actualContacts)
	require.Equal(t, expectedKinds, actualKinds)
}

func TestIterVerifiedContacts(t *testing.T) {
	contacts := &pb.Contacts{
		Technical: &pb.Contact{
			Email: "technical@example.com",
		},
		Administrative: &pb.Contact{
			Email: "administrative@example.com",
		},
		Billing: &pb.Contact{
			Email: "billing@example.com",
		},
		Legal: &pb.Contact{
			Email: "legal@example.com",
		},
	}

	actualContacts := []*pb.Contact{}
	actualKinds := []string{}

	// No contacts are verified.
	iter := models.NewContactIterator(contacts, models.SkipUnverified())
	for iter.Next() {
		contact, kind := iter.Value()
		actualContacts = append(actualContacts, contact)
		actualKinds = append(actualKinds, kind)
	}
	require.NoError(t, iter.Error())
	require.Equal(t, []*pb.Contact{}, actualContacts)
	require.Equal(t, []string{}, actualKinds)

	actualContacts = []*pb.Contact{}
	actualKinds = []string{}

	// Should only iterate through the verified contacts.
	require.NoError(t, models.SetContactVerification(contacts.Technical, "", true))
	require.NoError(t, models.SetContactVerification(contacts.Legal, "", true))
	expectedContacts := []*pb.Contact{
		contacts.Technical,
		contacts.Legal,
	}
	expectedKinds := []string{
		models.TechnicalContact,
		models.LegalContact,
	}
	iter = models.NewContactIterator(contacts, models.SkipUnverified())
	for iter.Next() {
		contact, kind := iter.Value()
		actualContacts = append(actualContacts, contact)
		actualKinds = append(actualKinds, kind)
	}
	require.NoError(t, iter.Error())
	require.Equal(t, expectedContacts, actualContacts)
	require.Equal(t, expectedKinds, actualKinds)
}

func TestIterDuplicates(t *testing.T) {
	contacts := &pb.Contacts{
		Technical: &pb.Contact{
			Email: "data@enterprised.com",
		},
		Administrative: &pb.Contact{
			Email: "picard@enterpised.com",
		},
		Legal: &pb.Contact{
			Email: "troi@enterprised.com",
		},
		Billing: &pb.Contact{
			Email: "riker@enterprised.com",
		},
	}

	actualContacts := []*pb.Contact{}
	actualKinds := []string{}

	// Should iterate over all contacts if there are no duplicates.
	expectedContacts := []*pb.Contact{
		contacts.Technical,
		contacts.Administrative,
		contacts.Legal,
		contacts.Billing,
	}
	expectedKinds := []string{
		models.TechnicalContact,
		models.AdministrativeContact,
		models.LegalContact,
		models.BillingContact,
	}

	iter := models.NewContactIterator(contacts, models.SkipDuplicates())
	for iter.Next() {
		contact, kind := iter.Value()
		actualContacts = append(actualContacts, contact)
		actualKinds = append(actualKinds, kind)
	}

	require.NoError(t, iter.Error())
	require.Equal(t, expectedContacts, actualContacts)
	require.Equal(t, expectedKinds, actualKinds)

	// Should skip duplicate contacts.
	actualContacts = []*pb.Contact{}
	actualKinds = []string{}
	contacts.Legal.Email = "riker@enterprised.com"
	expectedContacts = []*pb.Contact{
		contacts.Technical,
		contacts.Administrative,
		contacts.Legal,
	}
	expectedKinds = []string{
		models.TechnicalContact,
		models.AdministrativeContact,
		models.LegalContact,
	}

	iter = models.NewContactIterator(contacts, models.SkipDuplicates())
	for iter.Next() {
		contact, kind := iter.Value()
		actualContacts = append(actualContacts, contact)
		actualKinds = append(actualKinds, kind)
	}

	require.NoError(t, iter.Error())
	require.Equal(t, expectedContacts, actualContacts)
	require.Equal(t, expectedKinds, actualKinds)
}

func TestVerifiedContacts(t *testing.T) {
	testCases, err := loadContactFixtures()
	require.NoError(t, err, "could not load contact fixtures")

	for i, tc := range testCases {
		verified := tc.Contacts.VerifiedContacts()

		if tc.Duplicates == 0 {
			// If there are no duplicates the verified contacts should exactly match.
			require.Len(t, verified, tc.Verified, "test case %d failed", i)
		} else {
			// Otherwise there should be at least as many verified contacts, more indicates duplicates.
			require.GreaterOrEqual(t, len(verified), tc.Verified, "test case %d failed", i)
		}

		for kind, email := range verified {
			contact := tc.Contacts.Get(kind)
			require.NotNil(t, contact, "expected %s contact (%s) to not be nil in test case %d", kind, email, i)
			require.True(t, contact.Email.Verified, "expected %s contact (%s) to be verified in test case %d", kind, email, i)
		}
	}
}

func TestGetSentEmailCount(t *testing.T) {
	contacts := &pb.Contacts{
		Technical: &pb.Contact{
			Email: "technical@example.com",
		},
		Administrative: &pb.Contact{
			Email: "administrative@example.com",
		},
		Billing: &pb.Contact{
			Email: "billing@example.com",
		},
		Legal: &pb.Contact{
			Email: "legal@example.com",
		},
	}

	// Log should initially be empty
	emailLog, err := models.GetEmailLog(contacts.Administrative)
	require.NoError(t, err)
	require.Len(t, emailLog, 0)

	// Error should be returned if the reason is empty
	_, err = models.CountSentEmails(emailLog, "", 30)
	require.EqualError(t, err, "cannot match on empty reason string")

	// Error should be returned if the time window is invalid
	_, err = models.CountSentEmails(emailLog, "test", -1)
	require.EqualError(t, err, "time window must be a positive number of days")

	// Append an entry to an empty log
	err = models.AppendEmailLog(contacts.Administrative, "verify_contact", "verification")
	require.NoError(t, err)

	// Append an entry to an empty log
	err = models.AppendEmailLog(contacts.Administrative, "verify_contact", "verification")
	require.NoError(t, err)

	// Get email log for contact
	emailLog, err = models.GetEmailLog(contacts.Administrative)
	require.NoError(t, err)
	require.Len(t, emailLog, 2)
	require.Equal(t, "verify_contact", emailLog[0].Reason)
	require.Equal(t, "verification", emailLog[0].Subject)

	// Should return 2 emails sent for contact
	sent, err := models.CountSentEmails(emailLog, "verify_contact", 30)
	require.NoError(t, err)
	require.Equal(t, 2, sent)

	// Get the technical contact's email log
	emailLog, err = models.GetEmailLog(contacts.Technical)
	require.NoError(t, err)

	// Should return 0 emails when the log is empty
	sent, err = models.CountSentEmails(emailLog, "verify_contact", 30)
	require.NoError(t, err)
	require.Equal(t, 0, sent)

	// Construct an email log with entries at different times
	log := []*models.EmailLogEntry{
		{
			Timestamp: time.Now().Add(-time.Hour * 24 * 31).Format(time.RFC3339),
			Reason:    "verify_contact",
			Subject:   "verification",
		},
		{
			Timestamp: time.Now().Add(-time.Hour * 24 * 29).Format(time.RFC3339),
			Reason:    "verify_contact",
			Subject:   "verification",
		},
		{
			Timestamp: time.Now().Add(-time.Hour * 24 * 28).Format(time.RFC3339),
			Reason:    "verify_contact",
			Subject:   "verification",
		},
	}
	require.NoError(t, SetEmailLog(contacts.Billing, log))

	// Get the billing contact's email log
	emailLog, err = models.GetEmailLog(contacts.Billing)
	require.NoError(t, err)

	// Should only return a count of emails within the time window
	sent, err = models.CountSentEmails(emailLog, "verify_contact", 32)
	require.NoError(t, err)
	require.Equal(t, 3, sent, "expected 3 emails sent within the last 32 days")

	sent, err = models.CountSentEmails(emailLog, "verify_contact", 30)
	require.NoError(t, err)
	require.Equal(t, 2, sent, "expected 2 emails sent within the last 30 days")

	sent, err = models.CountSentEmails(emailLog, "verify_contact", 27)
	require.NoError(t, err)
	require.Equal(t, 0, sent, "expected 0 emails sent within the last 27 days")

	// Construct an email log with entries of different reasons
	log = []*models.EmailLogEntry{
		{
			Timestamp: time.Now().Add(-time.Hour * 24 * 31).Format(time.RFC3339),
			Reason:    "verify_contact",
			Subject:   "verification",
		},
		{
			Timestamp: time.Now().Add(-time.Hour * 24 * 29).Format(time.RFC3339),
			Reason:    "rejection",
			Subject:   "rejected registration",
		},
		{
			Timestamp: time.Now().Add(-time.Hour * 24 * 28).Format(time.RFC3339),
			Reason:    "verify_contact",
			Subject:   "verification",
		},
	}
	require.NoError(t, SetEmailLog(contacts.Legal, log))

	// Get the legal contact's email log
	emailLog, err = models.GetEmailLog(contacts.Legal)
	require.NoError(t, err)

	// Should only return a count of emails that match the reason and are within the time window
	sent, err = models.CountSentEmails(emailLog, "verify_contact", 32)
	require.NoError(t, err)
	require.Equal(t, 2, sent, "expected 2 emails sent within the last 32 days")

	sent, err = models.CountSentEmails(emailLog, "rejection", 30)
	require.NoError(t, err)
	require.Equal(t, 1, sent, "expected 1 emails sent within the last 30 days")
}

func TestEmailValidation(t *testing.T) {
	testCases := []struct {
		email         *models.Email
		err           error
		expectedName  string
		expectedEmail string
	}{
		{&models.Email{}, models.ErrNoEmailAddress, "", ""},
		{&models.Email{Email: "\t\n\t\n\t\t\t"}, models.ErrNoEmailAddress, "", ""},
		{&models.Email{Email: "ted@example.com", Verified: true, VerifiedOn: "", Token: ""}, models.ErrVerifiedInvalid, "", ""},
		{&models.Email{Email: "ted@example.com", Verified: true, VerifiedOn: "", Token: "foo"}, models.ErrVerifiedInvalid, "", ""},
		{&models.Email{Email: "ted@example.com", Verified: true, VerifiedOn: "2023-09-06T16:05:45-05:00", Token: "foo"}, models.ErrVerifiedInvalid, "", ""},
		{&models.Email{Email: "ted@example.com", Verified: true, VerifiedOn: "2023-09-06T16:05:45-05:00", Token: ""}, nil, "", "ted@example.com"},
		{&models.Email{Email: "ted@example.com", Verified: false, VerifiedOn: "", Token: ""}, models.ErrUnverifiedInvalid, "", ""},
		{&models.Email{Email: "ted@example.com", Verified: false, VerifiedOn: "2023-09-06T16:05:45-05:00", Token: ""}, models.ErrUnverifiedInvalid, "", ""},
		{&models.Email{Email: "ted@example.com", Verified: false, VerifiedOn: "2023-09-06T16:05:45-05:00", Token: "foo"}, models.ErrUnverifiedInvalid, "", ""},
		{&models.Email{Email: "ted@example.com", Verified: false, VerifiedOn: "", Token: "foo"}, nil, "", "ted@example.com"},
		{&models.Email{Email: "TED@example.com", Verified: false, VerifiedOn: "", Token: "foo"}, nil, "", "ted@example.com"},
		{&models.Email{Email: "Ted Tonks <TED@example.com>", Verified: false, VerifiedOn: "", Token: "foo"}, nil, "Ted Tonks", "ted@example.com"},
		{&models.Email{Name: "James Surry", Email: "Ted Tonks <TED@example.com>", Verified: false, VerifiedOn: "", Token: "foo"}, nil, "James Surry", "ted@example.com"},
		{&models.Email{Email: "TED@example.com", Verified: true, VerifiedOn: "2023-09-06T16:05:45-05:00", Token: ""}, nil, "", "ted@example.com"},
		{&models.Email{Email: "Ted Tonks <TED@example.com>", Verified: true, VerifiedOn: "2023-09-06T16:05:45-05:00", Token: ""}, nil, "Ted Tonks", "ted@example.com"},
		{&models.Email{Name: "James Surry", Email: "Ted Tonks <TED@example.com>", Verified: true, VerifiedOn: "2023-09-06T16:05:45-05:00", Token: ""}, nil, "James Surry", "ted@example.com"},
	}

	for i, tc := range testCases {
		err := tc.email.Validate()
		if tc.err == nil {
			require.NoError(t, err, "test case %d failed with error", i)
			require.Equal(t, tc.expectedName, tc.email.Name, "test case %d failed with name mismatch", i)
			require.Equal(t, tc.expectedEmail, tc.email.Email, "test case %d failed with email mismatch", i)
		} else {
			require.ErrorIs(t, err, tc.err, "test case %d failed with incorrect error", i)
		}
	}

}

func TestNormalizeEmail(t *testing.T) {
	testCases := []struct {
		email    string
		expected string
	}{
		{"support@trisa.io", "support@trisa.io"},
		{"Gary.Verdun@example.com", "gary.verdun@example.com"},
		{"   jessica@blankspace.net       ", "jessica@blankspace.net"},
		{"\t\t\nweird@foo.co.uk\t\n", "weird@foo.co.uk"},
		{"ALLCAPSCREAM@WILD.FR", "allcapscream@wild.fr"},
		{"Gary Verdun <gary@example.com>", "gary@example.com"},
	}

	for i, tc := range testCases {
		require.Equal(t, tc.expected, models.NormalizeEmail(tc.email), "test case %d failed", i)
	}
}

// Helper function to serialize an email log onto a contact's extra data.
func SetEmailLog(contact *pb.Contact, log []*models.EmailLogEntry) (err error) {
	extra := &models.GDSContactExtraData{}
	extra.EmailLog = log
	if contact.Extra, err = anypb.New(extra); err != nil {
		return err
	}
	return nil
}

const contactFixturesPath = "testdata/contacts.json.gz"

// Load Contact Fixtures from testdata, generating them if they do not exist.
// The contact fixtures are a series of test cases with different permutations of
// duplicate, verified, and unverified contacts and emails. Basically it exposes every
// possible combination of techincal, administrative, legal, and billing contact that
// is possible in the application: around 232 test cases in all.
func loadContactFixtures() (_ []*ContactsTestCase, err error) {
	// If the fixture does not exist, create it.
	if _, err = os.Stat(contactFixturesPath); os.IsNotExist(err) {
		if err = createContactFixtures(); err != nil {
			return nil, err
		}
	}

	var f *os.File
	if f, err = os.Open(contactFixturesPath); err != nil {
		return nil, err
	}
	defer f.Close()

	var r *gzip.Reader
	if r, err = gzip.NewReader(f); err != nil {
		return nil, err
	}
	defer r.Close()

	contacts := make([]*ContactsTestCase, 0)
	if err = json.NewDecoder(r).Decode(&contacts); err != nil {
		return nil, err
	}
	return contacts, nil
}

type ContactsTestCase struct {
	Contacts   *models.Contacts `json:"contacts"`
	Length     int              `json:"length"`
	Verified   int              `json:"verified"`
	Unverified int              `json:"unverified"`
	Duplicates int              `json:"duplicates"`
}

func createContactFixtures() (err error) {
	fixtures := make([]*ContactsTestCase, 0)
	people := []string{"James Franklin <jfranklin@example.com>", "Judy Bloomfeld <jbloomfeld@example.com>", "Christine Cradock <cradock@example.com>", "Richard Simlake <rsimlake@example.com>"}

	// Create all combinations of unique contacts, e.g.
	// T, A, L, B, TA, TL, AL, TB, AB, LB, TAL, TAB, TLB, ALB, TALB
	for i := 1; i < len(people)+1; i++ {
		for _, persons := range combinations(people, i) {
			// Create contacts to determine the number of emails
			initial, _, _ := makeContacts(persons...)
			nEmails := len(initial.Emails)

			// Create verified/unverified masks
			for _, mask := range verifiedMasks(nEmails) {
				contacts, nContacts, nDuplicates := makeContacts(persons...)
				nVerified, nUnverified := 0, 0

				for j, isVerified := range mask {
					if isVerified {
						contacts.Emails[j].Verified = true
						contacts.Emails[j].VerifiedOn = time.Now().Format(time.RFC3339Nano)
						nVerified++
					} else {
						contacts.Emails[j].Verified = false
						contacts.Emails[j].Token = secrets.CreateToken(models.VerificationTokenLength)
						nUnverified++
					}
				}

				fixtures = append(fixtures, &ContactsTestCase{
					Contacts:   contacts,
					Length:     nContacts,
					Unverified: nUnverified,
					Verified:   nVerified,
					Duplicates: nDuplicates,
				})
			}
		}
	}

	// Create combinations of contacts that are duplicated e.g.
	// TTTT, TATT, TTLT, TTTB, TALT, TATB, TTLB,
	// TAAA, AAAA, AALA, AAAB, TALA, TAAB, AALB,
	// TLLL, LALL, LLLL, LLLB, TALL, TLLB, LALB,
	// TBBB, BABB, BBLB, BBBB, TABB, TBLB, BALB,
	for _, persons := range duplicateCombinations(people) {
		// Create contacts to determine the number of emails
		initial, _, _ := makeContacts(persons...)
		nEmails := len(initial.Emails)

		// Create verified/unverified masks
		for _, mask := range verifiedMasks(nEmails) {
			contacts, nContacts, nDuplicates := makeContacts(persons...)
			nVerified, nUnverified := 0, 0

			for j, isVerified := range mask {
				if isVerified {
					contacts.Emails[j].Verified = true
					contacts.Emails[j].VerifiedOn = time.Now().Format(time.RFC3339Nano)
					nVerified++
				} else {
					contacts.Emails[j].Verified = false
					contacts.Emails[j].Token = secrets.CreateToken(models.VerificationTokenLength)
					nUnverified++
				}
			}

			fixtures = append(fixtures, &ContactsTestCase{
				Contacts:   contacts,
				Length:     nContacts,
				Unverified: nUnverified,
				Verified:   nVerified,
				Duplicates: nDuplicates,
			})
		}
	}

	// Write the contacts fixture to disk.
	var f *os.File
	if f, err = os.Create(contactFixturesPath); err != nil {
		return err
	}
	defer f.Close()

	w := gzip.NewWriter(f)
	defer w.Close()

	return json.NewEncoder(w).Encode(fixtures)
}

// Simplified way to quickly make contacts records from strings for the technical,
// administrative, legal, and billing contacts. Specify names and emails in the form:
// "Full Name <email@addr.com>"; if you want to skip a contact type, use an empty string.
func makeContacts(addresses ...string) (*models.Contacts, int, int) {
	contacts := &models.Contacts{
		Contacts: &pb.Contacts{},
		Emails:   make([]*models.Email, 0, len(addresses)),
	}

	nContacts := 0
	nDuplicates := 0
	for idx, address := range addresses {
		if address == "" {
			continue
		}

		email, err := mail.ParseAddress(address)
		if err != nil {
			panic(err)
		}

		switch idx {
		case 0:
			contacts.Contacts.Technical = &pb.Contact{Name: email.Name, Email: email.Address}
		case 1:
			contacts.Contacts.Administrative = &pb.Contact{Name: email.Name, Email: email.Address}
		case 2:
			contacts.Contacts.Legal = &pb.Contact{Name: email.Name, Email: email.Address}
		case 3:
			contacts.Contacts.Billing = &pb.Contact{Name: email.Name, Email: email.Address}
		default:
			panic("too many email addresses")
		}

		// Count the number of contacts added
		nContacts++

		// Add email to list of emails if it hasn't already been added
		found := false
		for _, addr := range contacts.Emails {
			if addr.Email == email.Address {
				nDuplicates++
				found = true
				break
			}
		}

		if !found {
			ts := time.Now().Format(time.RFC3339Nano)
			contacts.Emails = append(contacts.Emails, &models.Email{Name: email.Name, Email: email.Address, Created: ts, Modified: ts})
		}
	}

	return contacts, nContacts, nDuplicates
}

func combinations(set []string, n int) (subsets [][]string) {
	length := uint(len(set))

	if n > len(set) {
		n = len(set)
	}

	// Go through all possible combinations of objects
	// from 1 (only first object in subset) to 2^length (all objects in subset)
	for subsetBits := 1; subsetBits < (1 << length); subsetBits++ {
		if n > 0 && bits.OnesCount(uint(subsetBits)) != n {
			continue
		}

		var subset []string

		for object := uint(0); object < length; object++ {
			// checks if object is contained in subset
			// by checking if bit 'object' is set in subsetBits
			if (subsetBits>>object)&1 == 1 {
				// add object to subset
				subset = append(subset, set[object])
			} else {
				subset = append(subset, "")
			}
		}
		// add subset to subsets
		subsets = append(subsets, subset)
	}
	return subsets
}

func duplicateCombinations(set []string) (subsets [][]string) {
	seen := make(map[string]struct{})
	for _, dup := range set {
		for i := 1; i < len(set)-1; i++ {
			for _, combo := range combinations(set, i) {
				// Create the duplicate
				for j, e := range combo {
					if e == "" {
						combo[j] = dup
					}
				}

				s := strings.Join(combo, "")
				if _, ok := seen[s]; !ok {
					subsets = append(subsets, combo)
					seen[s] = struct{}{}
				}
			}
		}
	}
	return subsets
}

var smask = []string{"T", "A", "L", "B"}

func verifiedMasks(n int) (subsets [][]bool) {
	elems := smask[0:n]
	subsets = append(subsets, make([]bool, n))
	for i := 1; i < n+1; i++ {
		for _, combo := range combinations(elems, i) {
			subset := make([]bool, n)
			for j, e := range combo {
				subset[j] = e != ""
			}
			subsets = append(subsets, subset)
		}
	}
	return subsets
}
