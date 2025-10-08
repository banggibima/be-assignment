package dto

type CreateOrderRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	BuyerID   string `json:"buyer_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type CreateOrderResponse struct {
	ID         string `json:"id"`
	ProductID  string `json:"product_id"`
	BuyerID    string `json:"buyer_id"`
	Quantity   int    `json:"quantity"`
	TotalPrice int    `json:"total_price"`
	Status     string `json:"status"`
}

type GetOrderResponse struct {
	ID         string `json:"id"`
	ProductID  string `json:"product_id"`
	BuyerID    string `json:"buyer_id"`
	Quantity   int    `json:"quantity"`
	TotalPrice int    `json:"total_price"`
}
