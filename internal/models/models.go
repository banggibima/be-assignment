package models

import "time"

type Product struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Order struct {
	ID         string    `json:"id"`
	ProductID  string    `json:"product_id"`
	BuyerID    string    `json:"buyer_id"`
	Quantity   int       `json:"quantity"`
	TotalPrice int       `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Transaction struct {
	ID         string    `json:"id"`
	OrderID    string    `json:"order_id"`
	MerchantID string    `json:"merchant_id"`
	Amount     int       `json:"amount"`
	Fee        int       `json:"fee"`
	Status     string    `json:"status"`
	PaidAt     time.Time `json:"paid_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Settlement struct {
	ID          string    `json:"id"`
	MerchantID  string    `json:"merchant_id"`
	Date        time.Time `json:"date"`
	GrossAmount int       `json:"gross_amount"`
	FeeAmount   int       `json:"fee_amount"`
	NetAmount   int       `json:"net_amount"`
	TxnCount    int       `json:"txn_count"`
	UniqueRunID string    `json:"unique_run_id"`
	GeneratedAt time.Time `json:"generated_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Job struct {
	ID         string    `json:"id"`
	JobID      string    `json:"job_id"`
	Status     string    `json:"status"`
	Processed  int       `json:"processed"`
	Total      int       `json:"total"`
	Progress   int       `json:"progress"`
	From       string    `json:"from"`
	To         string    `json:"to"`
	ResultPath *string   `json:"result_path"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
