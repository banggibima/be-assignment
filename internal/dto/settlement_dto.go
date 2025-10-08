package dto

type CreateSettlementJobRequest struct {
	From string `json:"from" binding:"required"`
	To   string `json:"to" binding:"required"`
}

type CreateSettlementJobResponse struct {
	JobID   string `json:"job_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type SettlementProgressResponse struct {
	JobID       string  `json:"job_id"`
	Status      string  `json:"status"`
	Progress    int     `json:"progress"`
	Processed   int     `json:"processed"`
	Total       int     `json:"total"`
	DownloadURL *string `json:"download_url,omitempty"`
}
