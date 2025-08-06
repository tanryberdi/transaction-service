package http

import (
	"errors"
	"net/http"
	"strconv"

	"transaction-service/internal/application/services"
	"transaction-service/internal/domain/entities"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests
type Handler struct {
	transactionService *services.TransactionService
}

// NewHandler creates a new HTTP handler
func NewHandler(transactionService *services.TransactionService) *Handler {
	return &Handler{
		transactionService: transactionService,
	}
}

// SetupRoutes sets up the HTTP routes
func (h *Handler) SetupRoutes(router *gin.Engine) {

}

// ProcessTransaction handles POST /user/{userId}/transaction
func (h *Handler) ProcessTransaction(c *gin.Context) {
	// Extract user ID from the path
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Extract source type from header
	sourceTypeHeader := c.GetHeader("Source-Type")
	if sourceTypeHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Source-Type header is required",
		})
		return
	}

	sourceType := entities.SourceType(sourceTypeHeader)
	if !sourceType.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Source-Type header. Must be one of: game, server, payment",
		})
		return
	}

	// Parse the request body
	var req entities.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate required fields
	if req.State == "" || req.Amount == "" || req.TransactionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "All fields (state, amount, transactionID) are required",
		})
		return
	}

	// Process the transaction
	err = h.transactionService.ProcessTransaction(c.Request.Context(), userID, req, sourceType)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})

		case errors.Is(err, services.ErrInsufficientFunds):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Insufficient funds",
			})

		case errors.Is(err, services.ErrDuplicateTransaction):
			c.JSON(http.StatusConflict, gin.H{
				"error": "Transaction already processed",
			})

		case errors.Is(err, services.ErrInvalidAmount):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid amount format",
			})

		case errors.Is(err, services.ErrInvalidTransactionState):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid state. Must be 'win' or 'lose'",
			})

		case errors.Is(err, services.ErrInvalidSourceType):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Source-Type. Must be one of: game, server, payment",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error: " + err.Error(),
			})
		}

		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction processed successfully",
		"status":  "success",
	})
}

// GetUserBalance handles GET /user/{userId}/balance
func (h *Handler) GetUserBalance(c *gin.Context) {
	// Extract user ID from the path
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID. Must be a positive integer.",
		})
		return
	}

	// Get user balance
	balance, err := h.transactionService.GetUserBalance(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error: " + err.Error(),
		})
		return
	}

	// Return the balance
	c.JSON(http.StatusOK, gin.H{
		"userId":  userID,
		"balance": balance,
	})
}
