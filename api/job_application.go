package api

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// JobApplication model
type JobApplication struct {
	ID            int    `db:"id" json:"id"`
	FullName      string `db:"full_name" json:"full_name"`
	Age           int    `db:"age" json:"age"`
	Degree        string `db:"degree" json:"degree"`
	Experience    string `db:"experience" json:"experience"`
	Email         string `db:"email" json:"email"`
	Resume        string `db:"resume" json:"resume"`
	ApproveStatus string `db:"approve_status" json:"approve_status"`
	CreatedAt     string `db:"created_at" json:"created_at"`
}

// RegisterRoutes sets up job application endpoints
func RegisterRoutes(r *gin.Engine, db *sqlx.DB) {

	// Download a resume file
	r.GET("/job-applications/:id/download", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var application JobApplication

		err := db.Get(&application, "SELECT resume FROM job_applications WHERE id = ?", id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Job application not found"})
			return
		}

		resumeFilePath := application.Resume
		if resumeFilePath == "" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Resume file not found"})
			return
		}

		ctx.Header("Content-Disposition", "attachment; filename="+filepath.Base(resumeFilePath))
		ctx.File(resumeFilePath)
	})

	// Create a new job application
	r.POST("/job-applications", func(ctx *gin.Context) {
		var application JobApplication
		var resumeFilePath string

		file, header, err := ctx.Request.FormFile("resume")
		if err == nil {
			defer file.Close()
			resumeFilePath = filepath.Join("uploads", header.Filename)
			if err := ctx.SaveUploadedFile(header, resumeFilePath); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save resume"})
				return
			}
		} else {
			if err := ctx.ShouldBindJSON(&application); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON or missing resume file"})
				return
			}
			resumeFilePath = application.Resume
		}

		fullName := ctx.PostForm("full_name")
		age, _ := strconv.Atoi(ctx.PostForm("age"))
		degree := ctx.PostForm("degree")
		experience := ctx.PostForm("experience")
		email := ctx.PostForm("email")

		if fullName == "" {
			fullName = application.FullName
			age = application.Age
			degree = application.Degree
			experience = application.Experience
			email = application.Email
		}

		query := `INSERT INTO job_applications (full_name, age, degree, experience, email, resume, approve_status) 
		          VALUES (?, ?, ?, ?, ?, ?, 'Pending')`
		result, err := db.Exec(query, fullName, age, degree, experience, email, resumeFilePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		id, _ := result.LastInsertId()
		ctx.JSON(http.StatusCreated, gin.H{
			"id":             id,
			"full_name":      fullName,
			"age":            age,
			"degree":         degree,
			"experience":     experience,
			"email":          email,
			"resume":         resumeFilePath,
			"approve_status": "Pending",
		})
	})

	// Get all job applications
	r.GET("/job-applications", func(ctx *gin.Context) {
		var applications []JobApplication
		err := db.Select(&applications, "SELECT * FROM job_applications")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, applications)
	})

	// Get accepted job applications
	r.GET("/job-applications/accepted", func(ctx *gin.Context) {
		var acceptedApplications []JobApplication
		err := db.Select(&acceptedApplications, "SELECT * FROM job_applications WHERE approve_status = 'Accepted'")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, acceptedApplications)
	})

	// Accept a job application (Email sending removed)
	r.PUT("/job-applications/:id/accept", func(ctx *gin.Context) {
		id := ctx.Param("id")

		query := `UPDATE job_applications SET approve_status = 'Accepted' WHERE id = ?`
		_, err := db.Exec(query, id)
		if err != nil {
			log.Println("Error updating database for ID:", id, "Error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Application accepted",
		})
	})

	// Reject a job application (Email sending removed)
	r.PUT("/job-applications/:id/reject", func(ctx *gin.Context) {
		id := ctx.Param("id")

		query := `UPDATE job_applications SET approve_status = 'Rejected' WHERE id = ?`
		_, err := db.Exec(query, id)
		if err != nil {
			log.Println("Error updating database for ID:", id, "Error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Application rejected",
		})
	})

		// Get rejected job applications
		r.GET("/job-applications/rejected", func(ctx *gin.Context) {
			var rejectedApplications []JobApplication
			err := db.Select(&rejectedApplications, "SELECT * FROM job_applications WHERE approve_status = 'Rejected'")
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusOK, rejectedApplications)
		})

	// Delete an accepted job application
	r.DELETE("/job-applications/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		var count int
		err := db.Get(&count, "SELECT COUNT(*) FROM job_applications WHERE id = ? AND approve_status = 'Accepted'", id)
		if err != nil || count == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Accepted application not found"})
			return
		}

		_, err = db.Exec("DELETE FROM job_applications WHERE id = ?", id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Accepted application deleted successfully"})
	})

	// Delete a rejected job application
r.DELETE("/job-applications/:id/rejected", func(ctx *gin.Context) {
	id := ctx.Param("id")

	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM job_applications WHERE id = ? AND approve_status = 'Rejected'", id)
	if err != nil || count == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Rejected application not found"})
		return
	}

	_, err = db.Exec("DELETE FROM job_applications WHERE id = ?", id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Rejected application deleted successfully"})
})

// Restore a rejected job application
r.PUT("/job-applications/:id/restore", func(ctx *gin.Context) {
	id := ctx.Param("id")

	// Check if the application exists and is rejected
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM job_applications WHERE id = ? AND approve_status = 'Rejected'", id)
	if err != nil || count == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Rejected application not found"})
		return
	}

	// Update the status back to 'Pending'
	_, err = db.Exec("UPDATE job_applications SET approve_status = 'Pending' WHERE id = ?", id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Rejected application restored successfully"})
})


}
