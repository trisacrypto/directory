package server

import (
	"math/big"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/trisacrypto/directory/pkg/sectigo"
)

func (s *Server) CreateSingleCertBatch(c *gin.Context) {
	info := &sectigo.CreateSingleCertBatchRequest{}
	if err := c.ShouldBindJSON(info); err != nil {
		c.JSON(http.StatusBadRequest, Err(err))
		return
	}

	// TODO: handle profile validation and authorities/profiles

	rep := &sectigo.BatchResponse{}
	c.JSON(http.StatusCreated, rep)
}

func (s *Server) UploadCSRBatch(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, Err("this endpoint has not been implemented yet"))
}

func (s *Server) BatchDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, Err("bach id not found"))
		return
	}

	rep := &sectigo.BatchResponse{
		BatchID: int(id),
	}
	c.JSON(http.StatusOK, rep)
}

func (s *Server) BatchStatus(c *gin.Context) {
	_, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, Err("bach id not found"))
		return
	}

	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte("READY_FOR_DOWNLOAD"))
}

func (s *Server) ProcessingInfo(c *gin.Context) {
	_, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, Err("bach id not found"))
		return
	}

	rep := &sectigo.ProcessingInfoResponse{}
	c.JSON(http.StatusOK, rep)
}

func (s *Server) Download(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, Err("not implemented yet"))
}

func (s *Server) LicensesUsed(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, Err("not implemented yet"))
}

func (s *Server) UserAuthorities(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, Err("not implemented yet"))
}

func (s *Server) AuthorityAvailableBalance(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, Err("not implemented yet"))
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
	c.JSON(http.StatusNotImplemented, Err("not implemented yet"))
}

func (s *Server) FindCertificate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, Err("not implemented yet"))
}

func (s *Server) RevokeCertificate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, Err("not implemented yet"))
}

func SerialNumber() *big.Int {
	sn := make([]byte, 16)
	rand.Read(sn)

	i := &big.Int{}
	return i.SetBytes(sn)
}
