package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func createActivity(c *gin.Context) {
	userID := c.GetUint("userID")

	var req Activity
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	activity := Activity{
		ProfessorID: userID,
		Title:       req.Title,
		Description: req.Description,
		Language:    req.Language,
		TimeLimit:   req.TimeLimit,
		InviteToken: generateInviteToken(),
	}

	if err := db.Create(&activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity"})
		return
	}

	c.JSON(http.StatusCreated, activity)
}

func getActivities(c *gin.Context) {
	userID := c.GetUint("userID")

	var activities []Activity
	db.Where("professor_id = ?", userID).Find(&activities)

	// Add submission count
	type ActivityWithCount struct {
		Activity
		SubmissionCount int64 `json:"submissionCount"`
	}

	var result []ActivityWithCount
	for _, activity := range activities {
		var count int64
		db.Model(&Submission{}).Where("activity_id = ?", activity.ID).Count(&count)
		result = append(result, ActivityWithCount{
			Activity:        activity,
			SubmissionCount: count,
		})
	}

	c.JSON(http.StatusOK, result)
}

func getActivity(c *gin.Context) {
	activityID := c.Param("id")

	var activity Activity
	if err := db.First(&activity, activityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	c.JSON(http.StatusOK, activity)
}

func getSubmissions(c *gin.Context) {
	activityID := c.Param("id")

	var submissions []Submission
	db.Where("activity_id = ?", activityID).Order("created_at desc").Find(&submissions)

	c.JSON(http.StatusOK, submissions)
}

func joinActivity(c *gin.Context) {
	inviteToken := c.Param("inviteToken")

	var activity Activity
	if err := db.Where("invite_token = ?", inviteToken).First(&activity).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired invite link"})
		return
	}

	// Create anonymous student if not logged in
	student := User{
		Email: generateAnonymousEmail(),
		Name:  "Anonymous Student",
		Role:  "student",
	}
	db.Create(&student)

	c.JSON(http.StatusOK, gin.H{
		"activity": activity,
		"student":  student,
	})
}
