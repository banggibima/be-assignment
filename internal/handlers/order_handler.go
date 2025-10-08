package handlers

import (
	"net/http"

	"github.com/banggibima/be-assignment/internal/dto"
	"github.com/banggibima/be-assignment/internal/services"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	OrderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		OrderService: orderService,
	}
}

func (h *OrderHandler) Register(r *gin.Engine) {
	r.GET("/orders/:id", h.GetByID)
	r.POST("/orders", h.Create)
}

// Create godoc
// @Summary Create Order
// @Description Create a new order
// @Tags Order
// @Accept json
// @Produce json
// @Param request body dto.CreateOrderRequest true "Order request"
// @Success 201 {object} dto.CreateOrderResponse
// @Failure 400 {object} dto.ErrorResponse "Bad Request"
// @Failure 404 {object} dto.ErrorResponse "PRODUCT_NOT_FOUND"
// @Failure 409 {object} dto.ErrorResponse "OUT_OF_STOCK"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /orders [post]
func (h *OrderHandler) Create(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.OrderService.CreateOrder(c.Request.Context(), req)
	if err != nil {
		switch err {
		case services.ErrOutOfStock:
			c.JSON(http.StatusConflict, gin.H{"error": "OUT_OF_STOCK"})
			return
		case services.ErrProductNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "PRODUCT_NOT_FOUND"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, resp)
}

// GetByID godoc
// @Summary Get Order By ID
// @Description Get order details by ID
// @Tags Order
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} dto.GetOrderResponse
// @Failure 404 {object} dto.ErrorResponse "ORDER_NOT_FOUND"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /orders/{id} [get]
func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.OrderService.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ORDER_NOT_FOUND"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
