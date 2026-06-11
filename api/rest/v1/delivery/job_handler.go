package delivery

import (
    "construction_transport_server/internal/job/usecase"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

type JobHandler struct {
    uc *usecase.JobUsecase
}

func NewJobHandler(uc *usecase.JobUsecase) *JobHandler {
    return &JobHandler{uc: uc}
}

func (h *JobHandler) ListMyJobs(c *gin.Context) {
    userID := c.GetInt64("auth_id")
    jobs, err := h.uc.ListAssignedJobs(c.Request.Context(), userID)
    if err != nil {
        SendError(c, http.StatusInternalServerError, err.Error())
        return
    }
    SendResponse(c, http.StatusOK, "jobs fetched", jobs)
}

func (h *JobHandler) UpdateStatus(c *gin.Context) {
    userID := c.GetInt64("auth_id")
    bookingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        SendError(c, http.StatusBadRequest, "invalid job id")
        return
    }
    var req struct {
        Status string `json:"status" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        SendError(c, http.StatusBadRequest, "status required")
        return
    }
    if err := h.uc.UpdateJobStatus(c.Request.Context(), bookingID, userID, req.Status); err != nil {
        SendError(c, http.StatusBadRequest, err.Error())
        return
    }
    SendResponse(c, http.StatusOK, "job status updated", nil)
}