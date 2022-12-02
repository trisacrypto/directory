package server

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/trisacrypto/directory/pkg/sectigo"
	"github.com/trisacrypto/trisa/pkg/trust"
)

// In-memory store to hold information about batches and certificates.
type Store struct {
	sync.RWMutex
	issued       int
	batchseq     int
	batches      map[int]*Batch
	certificates map[string]*Certificate
	commonName   map[string]StringSet
}

type Batch struct {
	BatchID       int
	OrderNumber   int
	CreationDate  string
	Profile       string
	Size          int
	Status        string
	BatchName     string
	RejectReason  string
	UserID        string
	Active        int
	Failed        int
	SerialNumbers []string
	Params        Params
	Expires       time.Time
	Cleaned       bool
}

type Certificate struct {
	DeviceID     string
	CommonName   string
	SerialNumber string
	CreationDate string
	Status       string
	Data         []byte
}

type StringSet map[string]struct{}

var exists = struct{}{}

func NewStore() (*Store, error) {
	store := &Store{
		batchseq:     9999,
		batches:      make(map[int]*Batch),
		certificates: make(map[string]*Certificate),
		commonName:   make(map[string]StringSet),
	}

	go store.cleaner()
	return store, nil
}

func (s *Store) Issued() int {
	s.RLock()
	defer s.RUnlock()
	return s.issued
}

func (s *Store) AddBatch(profile string, info *sectigo.CreateSingleCertBatchRequest) (Batch, error) {
	s.Lock()
	defer s.Unlock()

	// Increment batch sequence and create batch with new ID
	// NOTE: batch will expire in 2 days and will not be able to be downloaded after
	s.batchseq++
	batch := &Batch{
		BatchID:       s.batchseq,
		BatchName:     info.BatchName,
		OrderNumber:   0,
		CreationDate:  time.Now().Format(time.RFC3339Nano),
		Profile:       profile,
		Size:          1,
		Status:        sectigo.BatchStatusProcessing,
		Expires:       time.Now().AddDate(0, 0, 2),
		Active:        1,
		Failed:        0,
		SerialNumbers: nil,
		Params:        info.ProfileParams,
		Cleaned:       false,
	}

	s.batches[batch.BatchID] = batch
	return *batch, nil
}

func (s *Store) GetBatch(id int) (Batch, error) {
	s.RLock()
	defer s.RUnlock()
	if batch, ok := s.batches[id]; ok {
		return *batch, nil
	}
	return Batch{}, errors.New("batch not found")
}

func (s *Store) RejectBatch(batchID int, rejectReason string) error {
	s.Lock()
	defer s.Unlock()

	batch, ok := s.batches[batchID]
	if !ok {
		return errors.New("no batch found with that id")
	}

	if batch.Status != sectigo.BatchStatusProcessing {
		return errors.New("batch is not in processing mode")
	}

	batch.Status = sectigo.BatchStatusRejected
	batch.RejectReason = rejectReason
	batch.Active = 0
	batch.Failed = 1
	batch.SerialNumbers = nil
	return nil
}

func (s *Store) AddCert(batchID int, data []byte) error {
	s.Lock()
	defer s.Unlock()

	batch, ok := s.batches[batchID]
	if !ok {
		return errors.New("no batch found with that id")
	}

	cert := &Certificate{
		DeviceID:     uuid.NewString(),
		CreationDate: time.Now().Format(time.RFC3339Nano),
		Data:         make([]byte, len(data)),
	}

	// Copy the data into the cert
	copy(cert.Data, data)

	// Parse the certificate to get the subject info and serial number
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		return errors.New("could not decode first pem encoded block of data")
	}

	crt, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	// Get the common name and serial number from the certificate
	cert.CommonName = crt.Subject.String()
	cert.SerialNumber = fmt.Sprintf("%X", crt.SerialNumber)
	cert.Status = sectigo.BatchStatusReadyForDownload

	// Get the common name from the certificate and add to index
	if _, ok := s.commonName[crt.Subject.CommonName]; ok {
		s.commonName[crt.Subject.CommonName][cert.SerialNumber] = exists
	} else {
		s.commonName[crt.Subject.CommonName] = StringSet{cert.SerialNumber: {}}
	}

	s.issued++
	s.certificates[cert.SerialNumber] = cert

	// Update the batch with the status
	batch.Status = sectigo.BatchStatusReadyForDownload
	batch.SerialNumbers = []string{cert.SerialNumber}
	batch.Active = 0
	return nil
}

func (s *Store) GetCertData(serialNumber string) (cert *trust.Provider, err error) {
	data, ok := s.certificates[serialNumber]
	if !ok {
		return nil, errors.New("could not find certificate with serial number")
	}

	if data.Status != sectigo.BatchStatusReadyForDownload {
		return nil, fmt.Errorf("cannot get certs in status %s", data.Status)
	}

	if len(data.Data) == 0 {
		return nil, errors.New("no certificate data available")
	}

	cert = &trust.Provider{}
	if err = cert.Decode(data.Data); err != nil {
		return nil, err
	}
	return cert, nil
}

func (s *Store) Find(commonName, serialNumber string) []Certificate {
	s.RLock()
	defer s.RUnlock()

	sns := make(StringSet)
	if serialNumber != "" {
		sns[serialNumber] = exists
	}

	if commonName != "" {
		for sn := range s.commonName[commonName] {
			sns[sn] = exists
		}
	}

	certs := make([]Certificate, 0, len(sns))
	for sn := range sns {
		if cert, ok := s.certificates[sn]; ok {
			certs = append(certs, *cert)
		}
	}

	return certs
}

func (s *Store) Revoke(serialNumber string) error {
	s.Lock()
	defer s.Unlock()

	cert, ok := s.certificates[serialNumber]
	if !ok {
		return errors.New("certificate not found")
	}

	cert.Status = sectigo.BatchStatusRevoked
	cert.Data = nil
	return nil
}

// Go routine that cleans up the in-memory store every 8 hours.
func (s *Store) cleaner() {
	for {
		// Use time after to wait 8 hours (not using sleep so we can add select in the future)
		<-time.After(8 * time.Hour)

		now := time.Now()
		batches := make([]int, 0)
		s.RLock()
		for id, batch := range s.batches {
			if now.After(batch.Expires) && !batch.Cleaned {
				batches = append(batches, id)
			}
		}
		s.RUnlock()

		if len(batches) > 0 {
			for _, id := range batches {
				s.Lock()
				batch := s.batches[id]
				for _, sn := range batch.SerialNumbers {
					if cert, ok := s.certificates[sn]; ok {
						cert.Data = nil
					}
				}

				batch.Status = sectigo.BatchStatusExpired
				batch.Cleaned = true
				s.Unlock()
			}
		}
	}
}
