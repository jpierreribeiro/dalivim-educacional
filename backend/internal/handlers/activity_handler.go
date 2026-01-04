package handler

import (
	"net/http"
	"strconv"

	"dalivim/internal/service"

	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	activityService service.ActivityService
}

func NewActivityHandler(activityService service.ActivityService) *ActivityHandler {
	return &ActivityHandler{activityService: activityService}
}

type CreateActivityRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Language    string `json:"language" binding:"required"`
	TimeLimit   int    `json:"timeLimit" binding:"required,min=1"`
}

func (h *ActivityHandler) Create(c *gin.Context) {
	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")

	activity, err := h.activityService.Create(userID, req.Title, req.Description, req.Language, req.TimeLimit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, activity)
}

func (h *ActivityHandler) GetAll(c *gin.Context) {
	userID := c.GetUint("userID")

	activities, err := h.activityService.GetByProfessorID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, activities)
}

func (h *ActivityHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	activity, err := h.activityService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	c.JSON(http.StatusOK, activity)
}

func (h *ActivityHandler) Join(c *gin.Context) {
	inviteToken := c.Param("inviteToken")

	activity, student, err := h.activityService.JoinActivity(inviteToken)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired invite link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"activity": activity,
		"student":  student,
	})
}
