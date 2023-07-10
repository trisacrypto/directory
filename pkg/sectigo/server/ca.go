package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/trisacrypto/directory/pkg/sectigo"
	"github.com/trisacrypto/trisa/pkg/trust"
)

const (
	EcosystemID = 21
	UserID      = 295
)

func (s *Server) CreateSingleCertBatch(c *gin.Context) {
	info := &sectigo.CreateSingleCertBatchRequest{}
	if err := c.ShouldBindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, Err(err))
		return
	}

	// Determine the profile from the authority
	var profile string
	switch info.AuthorityID {
	case 423:
		profile = sectigo.ProfileCipherTraceEE
	case 489:
		profile = sectigo.ProfileCipherTraceEndEntityCertificate
	default:
		c.JSON(http.StatusBadRequest, Err("unknown profile"))
		return
	}

	// TODO: handle profile validation for params
	if _, ok := info.ProfileParams[sectigo.ParamCommonName]; !ok {
		c.JSON(http.StatusBadRequest, Err("common name required"))
		return
	}

	if _, ok := info.ProfileParams[sectigo.ParamPassword]; !ok {
		c.JSON(http.StatusBadRequest, Err("pkcs12 password required"))
		return
	}

	// Create a batch from the certs
	batch, err := s.store.AddBatch(profile, info)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Err(err))
		return
	}

	rep := &sectigo.BatchResponse{
		BatchID:         batch.BatchID,
		OrderNumber:     batch.OrderNumber,
		CreationDate:    batch.CreationDate,
		Profile:         batch.Profile,
		Size:            batch.Size,
		Status:          batch.Status,
		Active:          batch.Active > 0,
		BatchName:       batch.BatchName,
		RejectReason:    batch.RejectReason,
		GeneratorValues: nil,
		UserID:          UserID,
		Downloadable:    batch.Status == sectigo.BatchStatusReadyForDownload,
		Rejectable:      batch.Status == sectigo.BatchStatusProcessing,
	}

	// After a suitable delay to mimic the real service, issue certs
	go s.CreateCerts(info.ProfileParams, profile, batch.BatchID)

	// Return the batch response
	c.JSON(http.StatusCreated, rep)
}

func (s *Server) UploadCSRBatch(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, Err("this endpoint has not been implemented yet"))
}

func (s *Server) BatchDetail(c *gin.Context) {
	id, err := ParseID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, Err("batch id not found"))
		return
	}

	batch, err := s.store.GetBatch(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Err(err))
		return
	}

	rep := &sectigo.BatchResponse{
		BatchID:         batch.BatchID,
		OrderNumber:     batch.OrderNumber,
		CreationDate:    batch.CreationDate,
		Profile:         batch.Profile,
		Size:            batch.Size,
		Status:          batch.Status,
		Active:          batch.Active > 0,
		BatchName:       batch.BatchName,
		RejectReason:    batch.RejectReason,
		GeneratorValues: nil,
		UserID:          UserID,
		Downloadable:    batch.Status == sectigo.BatchStatusReadyForDownload,
		Rejectable:      batch.Status == sectigo.BatchStatusProcessing,
	}
	c.JSON(http.StatusOK, rep)
}

func (s *Server) BatchStatus(c *gin.Context) {
	id, err := ParseID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, Err("batch id not found"))
		return
	}

	batch, err := s.store.GetBatch(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Err(err))
		return
	}

	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(batch.Status))
}

func (s *Server) ProcessingInfo(c *gin.Context) {
	id, err := ParseID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, Err("batch id not found"))
		return
	}

	batch, err := s.store.GetBatch(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Err(err))
		return
	}

	rep := &sectigo.ProcessingInfoResponse{
		Active:  batch.Active,
		Success: len(batch.SerialNumbers),
		Failed:  batch.Failed,
	}
	c.JSON(http.StatusOK, rep)
}

func (s *Server) Download(c *gin.Context) {
	id, err := ParseID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, Err("batch id not found"))
		return
	}

	batch, err := s.store.GetBatch(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Err(err))
		return
	}

	if batch.Status != sectigo.BatchStatusReadyForDownload {
		c.JSON(http.StatusBadRequest, Err("batch is not in downloadable state"))
		return
	}

	// TODO: handle multiple certificates to download
	if len(batch.SerialNumbers) != 1 {
		c.JSON(http.StatusInternalServerError, Err("batch does not have associated certs"))
		return
	}

	// Get the pkcs12 password
	pkcs12password := batch.Params.Get(sectigo.ParamPassword, "")
	if pkcs12password == "" {
		c.JSON(http.StatusInternalServerError, Err("batch does not have pkcs12 password"))
		return
	}

	certData, err := s.store.GetCertData(batch.SerialNumbers[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, Err(err))
		return
	}

	sz, err := trust.NewSerializer(true, pkcs12password, trust.CompressionZIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Err(err))
		return
	}

	data, err := sz.Compress(certData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Err(err))
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%d.zip\"", batch.BatchID))
	c.Data(http.StatusOK, "application/octet-stream", data)
}

func (s *Server) LicensesUsed(c *gin.Context) {
	out := &sectigo.LicensesUsedResponse{
		Ordered: 1000,
		Issued:  s.store.Issued(),
	}
	c.JSON(http.StatusOK, out)
}

func (s *Server) UserAuthorities(c *gin.Context) {
	out := []*sectigo.AuthorityResponse{
		{
			ID:                  489,
			EcosystemID:         EcosystemID,
			SignerCertificateID: 0,
			EcosystemName:       "Staging",
			Balance:             0,
			Enabled:             true,
			ProfileID:           489,
			ProfileName:         "CipherTrace End Entity Certificate_#85",
		},
		{
			ID:                  423,
			EcosystemID:         EcosystemID,
			SignerCertificateID: 0,
			EcosystemName:       "Staging",
			Balance:             100,
			Enabled:             true,
			ProfileID:           423,
			ProfileName:         "CipherTrace EE_#17",
		},
	}
	c.JSON(http.StatusOK, out)
}

func (s *Server) AuthorityAvailableBalance(c *gin.Context) {
	c.JSON(http.StatusOK, 100)
}

func (s *Server) Profiles(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, Err("not implemented yet"))
}

func (s *Server) ProfileParams(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, Err("not implemented yet"))
}

func (s *Server) ProfileDetail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, Err("not implemented yet"))
}

func (s *Server) Organization(c *gin.Context) {
	out := &sectigo.OrganizationResponse{
		OrganizationID:      1042,
		OrganizationName:    "Staging",
		Address:             "123 Main Street",
		PrimaryContactName:  "Jane Wilkerson",
		PrimaryContactEmail: "jw@example.com",
		PrimaryContactPhone: "555-555-5555",
		ManufactureID:       "",
		Logo:                "",
		Authorities: []*sectigo.AuthorityResponse{
			{
				ID:                  489,
				EcosystemID:         EcosystemID,
				SignerCertificateID: 0,
				EcosystemName:       "Staging",
				Balance:             0,
				Enabled:             true,
				ProfileID:           489,
				ProfileName:         "CipherTrace End Entity Certificate_#85",
			},
			{
				ID:                  423,
				EcosystemID:         EcosystemID,
				SignerCertificateID: 0,
				EcosystemName:       "Staging",
				Balance:             100,
				Enabled:             true,
				ProfileID:           423,
				ProfileName:         "CipherTrace EE_#17",
			},
		},
		EcosystemID: EcosystemID,
		Parameters: map[string]string{
			"Organization": "Staging",
		},
		Status: "ACTIVE",
	}

	c.JSON(http.StatusOK, out)
}

func (s *Server) FindCertificate(c *gin.Context) {
	in := &sectigo.FindCertificateRequest{}
	if err := c.ShouldBindJSON(in); err != nil {
		c.JSON(http.StatusBadRequest, Err(err))
		return
	}

	if in.CommonName == "" && in.SerialNumber == "" {
		c.JSON(http.StatusBadRequest, Err("specify common name or serial number"))
		return
	}

	certs := s.store.Find(in.CommonName, in.SerialNumber)
	out := &sectigo.FindCertificateResponse{
		TotalCount: len(certs),
		Items:      make([]*sectigo.FindCertificateItem, 0, len(certs)),
	}

	for _, cert := range certs {
		out.Items = append(out.Items, &sectigo.FindCertificateItem{
			DeviceID:     cert.DeviceID,
			CommonName:   cert.CommonName,
			SerialNumber: cert.SerialNumber,
			CreationDate: cert.CreationDate,
			Status:       cert.Status,
		})
	}

	c.JSON(http.StatusOK, out)
}

func (s *Server) RevokeCertificate(c *gin.Context) {
	in := &sectigo.RevokeCertificateRequest{}
	if err := c.ShouldBindJSON(in); err != nil {
		c.JSON(http.StatusBadRequest, Err(err))
		return
	}

	if in.SerialNumber == "" {
		c.JSON(http.StatusBadRequest, Err("no serial number specified"))
		return
	}

	if err := s.store.Revoke(in.SerialNumber); err != nil {
		c.JSON(http.StatusNotFound, Err(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func ParseID(c *gin.Context) (id int, err error) {
	var id64 int64
	if id64, err = strconv.ParseInt(c.Param("id"), 10, 64); err != nil {
		return id, err
	}
	return int(id64), nil
}
