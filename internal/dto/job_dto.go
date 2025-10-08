package dto

type JobStatusResponse struct {
	JobID       string  `json:"job_id"`
	Status      string  `json:"status"`
	Progress    int     `json:"progress"`
	Processed   int     `json:"processed"`
	Total       int     `json:"total"`
	ResultPath  *string `json:"result_path,omitempty"`
	DownloadURL *string `json:"download_url,omitempty"`
}

type CancelJobResponse struct {
	JobID   string `json:"job_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
