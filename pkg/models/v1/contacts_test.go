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
)

func TestContactKind(t *testing.T) {
	testCases := []struct {
		kind   string
		assert require.BoolAssertionFunc
	}{
		{models.TechnicalContact, require.True},
		{models.AdministrativeContact, require.True},
		{models.LegalContact, require.True},
		{models.BillingContact, require.True},
		{"foo", require.False},
		{"ADMINISTRATIVE", require.False},
		{"Technical", require.False},
		{"  legal\t", require.False},
	}

	for i, tc := range testCases {
		tc.assert(t, models.ContactKindIsValid(tc.kind), "test case %d failed", i)
	}
}

func TestContacts(t *testing.T) {
	kinds := []string{
		models.TechnicalContact,
		models.AdministrativeContact,
		models.LegalContact,
		models.BillingContact,
	}

	// Create some simple fixtures for tests
	allNil := &models.Contacts{Contacts: &pb.Contacts{}, Emails: nil}
	allZero := &models.Contacts{Contacts: &pb.Contacts{Technical: &pb.Contact{}, Administrative: &pb.Contact{}, Legal: &pb.Contact{}, Billing: &pb.Contact{}}, Emails: make([]*models.Email, 0)}
	allPop := &models.Contacts{Contacts: &pb.Contacts{Technical: &pb.Contact{Email: "technical@example.com"}, Administrative: &pb.Contact{Email: "administrative@example.com"}, Legal: &pb.Contact{Email: "legal@example.com"}, Billing: &pb.Contact{Email: "billing@example.com"}}, Emails: []*models.Email{{Email: "technical@example.com", Token: "1"}, {Email: "administrative@example.com", Token: "2"}, {Email: "legal@example.com", Token: "3"}, {Email: "billing@example.com", Token: "4"}}}

	t.Run("Has", func(t *testing.T) {
		for _, kind := range kinds {
			require.False(t, allNil.Has(kind))
			require.False(t, allZero.Has(kind))
			require.True(t, allPop.Has(kind))
		}

		// Expect a panic when an unknown kind is passed in
		require.Panics(t, func() { allNil.Has("foo") }, "expected panic with unknown kind")
	})

	t.Run("Get", func(t *testing.T) {
		for _, kind := range kinds {
			require.Nil(t, allNil.Get(kind))

			contact := allZero.Get(kind)
			require.True(t, contact.IsZero())

			contact = allPop.Get(kind)
			require.False(t, contact.IsZero())
			require.NotEmpty(t, contact.Email)
			require.True(t, strings.HasPrefix(contact.Contact.Email, kind))
			require.True(t, strings.HasPrefix(contact.Email.Email, kind))
		}

		// Expect a panic when an unknown kind is passed in
		require.Panics(t, func() { allNil.Get("foo") }, "expected panic with unknown kind")
	})
}

func TestContactsIter(t *testing.T) {
	testCases, err := loadContactFixtures()
	require.NoError(t, err, "could not load contact fixtures")

	allNil := &models.Contacts{Contacts: &pb.Contacts{}, Emails: nil}
	allZero := &models.Contacts{Contacts: &pb.Contacts{Technical: &pb.Contact{}, Administrative: &pb.Contact{}, Legal: &pb.Contact{}, Billing: &pb.Contact{}}, Emails: make([]*models.Email, 0)}
	allNames := &models.Contacts{Contacts: &pb.Contacts{Technical: &pb.Contact{Name: "Tech Person"}, Administrative: &pb.Contact{Name: "Admin Person"}, Legal: &pb.Contact{Name: "Legal Person"}, Billing: &pb.Contact{Name: "Billing Person"}}, Emails: make([]*models.Email, 0)}
	allPop := &models.Contacts{Contacts: &pb.Contacts{Technical: &pb.Contact{Email: "technical@example.com"}, Administrative: &pb.Contact{Email: "administrative@example.com"}, Legal: &pb.Contact{Email: "legal@example.com"}, Billing: &pb.Contact{Email: "billing@example.com"}}, Emails: []*models.Email{{Email: "technical@example.com", Token: "1"}, {Email: "administrative@example.com", Token: "2"}, {Email: "legal@example.com", Token: "3"}, {Email: "billing@example.com", Token: "4"}}}

	t.Run("Empty", func(t *testing.T) {
		for _, contacts := range []*models.Contacts{allNil, allZero} {
			iter := contacts.NewIterator()
			require.False(t, iter.Next())
		}

		require.Equal(t, 4, allNames.Length())
		require.Equal(t, 4, allPop.Length())
	})

	t.Run("All", func(t *testing.T) {
		for i, tc := range testCases {
			iter := tc.Contacts.NewIterator()
			count := 0

			for iter.Next() {
				contact := iter.Contact()
				require.NotNil(t, contact, "iter returned nil contact in test case %d", i)
				count++
			}

			require.Equal(t, tc.Length, count)
			require.Equal(t, count, tc.Contacts.Length())
		}
	})

	t.Run("SkipNoEmail", func(t *testing.T) {
		// All test cases have emails so none will be skipped.
		for i, tc := range testCases {
			iter := tc.Contacts.NewIterator(models.SkipNoEmail())
			count := 0

			for iter.Next() {
				contact := iter.Contact()
				require.NotNil(t, contact, "iter returned nil contact in test case %d", i)
				require.NotEmpty(t, contact.Contact.Email, "expected email to not be empty in test case %d", i)
				count++
			}

			require.Equal(t, count, tc.Contacts.Length(models.SkipNoEmail()), "incorrect length on test case %d", i)
			require.Equal(t, tc.Length, count, "test case %d failed", i)
		}

		// Ensure that contacts without emails are skipped
		require.Equal(t, 0, allNames.Length(models.SkipNoEmail()))
	})

	t.Run("SkipUnverified", func(t *testing.T) {
		for i, tc := range testCases {
			iter := tc.Contacts.NewIterator(models.SkipUnverified())
			count := 0

			for iter.Next() {
				contact := iter.Contact()
				require.NotNil(t, contact, "iter returned nil contact in test case %d", i)
				require.True(t, contact.IsVerified())
				count++
			}

			if tc.Duplicates == 0 {
				require.Equal(t, tc.Verified, count)
			} else {
				require.GreaterOrEqual(t, count, tc.Verified)
			}

			require.Equal(t, count, tc.Contacts.Length(models.SkipUnverified()))
		}
	})

	t.Run("SkipVerified", func(t *testing.T) {
		for i, tc := range testCases {
			iter := tc.Contacts.NewIterator(models.SkipVerified())
			count := 0

			for iter.Next() {
				contact := iter.Contact()
				require.NotNil(t, contact, "iter returned nil contact in test case %d", i)
				require.False(t, contact.IsVerified())
				count++
			}

			if tc.Duplicates == 0 {
				require.Equal(t, tc.Unverified, count)
			} else {
				require.GreaterOrEqual(t, count, tc.Unverified)
			}

			require.Equal(t, count, tc.Contacts.Length(models.SkipVerified()))
		}
	})

	t.Run("SkipDuplicates", func(t *testing.T) {
		for i, tc := range testCases {
			iter := tc.Contacts.NewIterator(models.SkipDuplicates())
			count := 0

			for iter.Next() {
				contact := iter.Contact()
				require.NotNil(t, contact, "iter returned nil contact in test case %d", i)
				count++
			}

			require.Equal(t, tc.Length-tc.Duplicates, count)
			require.Equal(t, count, tc.Contacts.Length(models.SkipDuplicates()))
		}
	})

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

func TestContactVerifications(t *testing.T) {
	testCases, err := loadContactFixtures()
	require.NoError(t, err, "could not load contact fixtures")

	for i, tc := range testCases {
		contacts := tc.Contacts.ContactVerifications()
		require.Len(t, contacts, tc.Length, "test case %d failed", i)

		for kind, verified := range contacts {
			contact := tc.Contacts.Get(kind)
			require.NotNil(t, contact, "expected %s contact to not be nil in test case %d", kind, i)
			require.Equal(t, contact.IsVerified(), verified, "expected %s contact verified to be %t in test case %d", kind, verified, i)
		}
	}
}

type ContactsTestCase struct {
	Contacts   *models.Contacts `json:"contacts"`
	Length     int              `json:"length"`
	Verified   int              `json:"verified"`
	Unverified int              `json:"unverified"`
	Duplicates int              `json:"duplicates"`
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
