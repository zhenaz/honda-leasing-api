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

/*
====================================
CREATE CONTRACT
====================================
*/

func (h *LeasingHandler) CreateContract(c *gin.Context) {

	var req service.CreateContractRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	contract, err := h.service.CreateContract(
		c.Request.Context(),
		&req,
	)
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

/*
====================================
GET CONTRACT DETAIL
====================================
*/

func (h *LeasingHandler) GetContractDetail(c *gin.Context) {

	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid contract id",
		})
		return
	}

	result, err := h.service.GetContractDetail(
		c.Request.Context(),
		id,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

/*
====================================
APPROVE CONTRACT
====================================
*/

func (h *LeasingHandler) ApproveContract(c *gin.Context) {

	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid contract id",
		})
		return
	}

	err = h.service.ApproveContract(
		c.Request.Context(),
		id,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "contract approved",
	})
}

func (h *LeasingHandler) ListContracts(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	contracts, total, err := h.service.ListContracts(
		c.Request.Context(),
		page,
		limit,
		status,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": contracts,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func (h *LeasingHandler) CompleteTask(c *gin.Context) {

	idStr := c.Param("id")

	taskID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task id"})
		return
	}

	err = h.service.CompleteTask(
		c.Request.Context(),
		taskID,
	)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "task completed",
	})
}
