package controllers

import (
	"net/http"
	"strconv"

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

type taskStatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
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

	// Create task via orm engine
	result := database.DB.Create(&newTask)

	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create a new task"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Task created successfully"})
}

func GetTasks(context *gin.Context) {

	// Check query string for filtering
	var pageStr string = context.Query("page")
	if pageStr == "" {
		pageStr = "1" // Default to page 1 if not provided
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid page parameter"})
		return
	}

	statusStr := context.Query("status")

	var tasks []models.Task
	userEmail := context.GetString("userEmail")

	// Build base query
	query := database.DB.Model(&models.Task{}).Where("user_email = ?", userEmail)

	// Apply status filter if provided
	if statusStr != "" {
		query = query.Where("status = ?", statusStr)
	}

	// Get total count of tasks for pagination
	var totalRecords int64
	query.Model(&models.Task{}).Count(&totalRecords)

	// Calculate total pages
	size := 5
	totalPages := int((totalRecords + int64(size) - 1) / int64(size)) // Ceiling division
	if totalPages == 0 {
		totalPages = 1 // At least 1 page even if no records
	}

	if page > totalPages {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Page number exceeds total pages"})
		return
	}

	// Fetch tasks with pagination
	order := context.Query("order")
	if order == "" {
		order = "due_date DESC" // Default order by due date descending
	} else {
		order = "due_date " + order // Use the provided order
	}

	result := query.Offset((page - 1) * size).Limit(size).Order(order).Find(&tasks)

	// Fetch tasks depending on the user email and pagination parameters
	// result := database.DB.Where("user_email = ?", userEmail).Find(&tasks)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch tasks"})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Fetched tasks successfully",
		"tasks":   tasks,
		"total":   totalRecords,
		"page":    page,
		"size":    size,
		"pages":   totalPages,
	})
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

func GetTaskStatusCounts(context *gin.Context) {
	userEmail := context.GetString("userEmail")

	var statusCounts []taskStatusCount
	// Query to get count of tasks grouped by status
	result := database.DB.Model(&models.Task{}).
		Select("status, COUNT(*) as count").
		Where("user_email = ?", userEmail).
		Group("status").
		Scan(&statusCounts)

	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch task status counts"})
		return
	}

	// Create a map for easier frontend consumption
	statusMap := make(map[string]int64)
	var totalTasks int64 = 0

	for _, sc := range statusCounts {
		statusMap[sc.Status] = sc.Count
		totalTasks += sc.Count
	}

	// Ensure all possible statuses are included (even if count is 0)
	possibleStatuses := []string{"Todo", "In Progress", "Done"}
	for _, status := range possibleStatuses {
		if _, exists := statusMap[status]; !exists {
			statusMap[status] = 0
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"message":       "Fetched task status counts successfully",
		"total_tasks":   totalTasks,
		"status_counts": statusMap,
	})
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
