package handler

import (
	"net/http"
	"strconv"

	service "honda-leasing-api/internal/services"

	"github.com/gin-gonic/gin"
)

type LeasingHandler struct {
	service service.LeasingService
}

func NewLeasingHandler(service service.LeasingService) *LeasingHandler {
	return &LeasingHandler{service: service}
}

func (h *LeasingHandler) CreateContract(c *gin.Context) {

	var req service.CreateContractRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	contract, err := h.service.CreateContract(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": contract,
	})
}

func (h *LeasingHandler) GetContractDetail(c *gin.Context) {

	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid contract id"})
		return
	}

	result, err := h.service.GetContractDetail(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": result})
}
