package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/banggibima/be-assignment/internal/dto"
	"github.com/banggibima/be-assignment/internal/services"
	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	JobService *services.JobService
}

func NewJobHandler(jobService *services.JobService) *JobHandler {
	return &JobHandler{
		JobService: jobService,
	}
}

func (h *JobHandler) Register(r *gin.Engine) {
	r.GET("/jobs/:id", h.GetJob)
	r.POST("/jobs/:id/cancel", h.CancelJob)
	r.POST("/jobs/settlement", h.StartJob)
	r.GET("/downloads/:job_id", h.Download)
}

// StartJob godoc
// @Summary Create Settlement Job
// @Description Create a new settlement job
// @Tags Job
// @Accept json
// @Produce json
// @Param request body dto.CreateSettlementJobRequest true "Job request"
// @Success 202 {object} dto.CreateSettlementJobResponse
// @Failure 400 {object} dto.ErrorResponse "Bad Request"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /jobs/settlement [post]
func (h *JobHandler) StartJob(c *gin.Context) {
	var req dto.CreateSettlementJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.JobService.CreateJob(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, res)
}

// GetJob godoc
// @Summary Get Job Status
// @Description Get the status of a specific job
// @Tags Job
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} dto.SettlementProgressResponse
// @Failure 404 {object} dto.ErrorResponse "JOB_NOT_FOUND"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /jobs/{id} [get]
func (h *JobHandler) GetJob(c *gin.Context) {
	jobID := c.Param("id")

	job, err := h.JobService.GetJobStatus(c.Request.Context(), jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "JOB_NOT_FOUND"})
		return
	}

	resp := map[string]interface{}{
		"job_id":    job.JobID,
		"status":    job.Status,
		"progress":  job.Progress,
		"processed": job.Processed,
		"total":     job.Total,
	}

	if job.DownloadURL != nil {
		resp["download_url"] = *job.DownloadURL
	}

	c.JSON(http.StatusOK, resp)
}

// CancelJob godoc
// @Summary Cancel Job
// @Description Cancel a running job
// @Tags Job
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} dto.CancelJobResponse
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /jobs/{id}/cancel [post]
func (h *JobHandler) CancelJob(c *gin.Context) {
	jobID := c.Param("id")

	res, err := h.JobService.CancelJob(jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// Download godoc
// @Summary Download Job Result
// @Description Download CSV file of completed job
// @Tags Job
// @Produce octet-stream
// @Param job_id path string true "Job ID"
// @Success 200 {file} string
// @Router /downloads/{job_id} [get]
func (h *JobHandler) Download(c *gin.Context) {
	jobID := c.Param("job_id")
	filename := jobID + ".csv"
	path := filepath.Join("/tmp/settlements", filename)

	c.File(path)
}
