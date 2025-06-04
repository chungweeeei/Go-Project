package controllers

import (
	"net/http"

	"example.com/restful-server/database"
	"example.com/restful-server/helper"
	"example.com/restful-server/models"
	"github.com/gin-gonic/gin"
)

type taskRequest struct {
	Title       string `json:"title";binding:"required"`
	Description string `json:"description";binding:"required"`
	Status      string `json:"status";binding:"required"`
	Priority    string `json:"priority";binding:"required"`
	Category    string `json:"category"`
	DueDate     string `json:"dueDate";binding:"required"` // Expecting date in YYYY-MM-DD format
}

func CreateTask(context *gin.Context) {
	var req taskRequest
	err := context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to parse request body"})
	}

	// get user email from context
	userEmail := context.GetString("userEmail")

	// Create a new task instance
	newTask := models.Task{
		UserEmail:   userEmail,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		Category:    req.Category,
		DueDate:     helper.FromISO(req.DueDate), // Convert string to time.Time
	}

	result := database.DB.Create(&newTask)

	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create a new task"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Task created successfully"})
}

func GetTasks(context *gin.Context) {

	var tasks []models.Task
	userEmail := context.GetString("userEmail")

	// Fetch tasks for the user
	result := database.DB.Where("user_email = ?", userEmail).Find(&tasks)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch tasks"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Fetched tasks successfully", "tasks": tasks})
}

func GetTaskByID(context *gin.Context) {
	taskID := context.Param("id")
	var task models.Task

	// Fetch the task by ID
	result := database.DB.Where("id = ?", taskID).First(&task)
	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Task not found"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Fetched task successfully", "task": task})
}

func UpdateTask(context *gin.Context) {
	taskID := context.Param("id")
	var req taskRequest
	err := context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to parse request body"})
		return
	}

	// get user email from context
	userEmail := context.GetString("userEmail")

	// Update the task
	result := database.DB.Model(&models.Task{}).Where("id = ? AND user_email = ?", taskID, userEmail).Updates(models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		Category:    req.Category,
		DueDate:     helper.FromISO(req.DueDate), // Convert string to time.Time
	})

	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update task"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func DeleteTask(context *gin.Context) {
	taskID := context.Param("id")
	userEmail := context.GetString("userEmail")

	// Delete the task
	result := database.DB.Where("id = ? AND user_email = ?", taskID, userEmail).Delete(&models.Task{})
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete task"})
		return
	}

	if result.RowsAffected == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Task not found"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
