package api

import (
	"net/http"

	"crud/util"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func AuthRoutes(r *gin.Engine, db *sqlx.DB) {
	r.POST("/login", func(c *gin.Context) {
		var account struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		var user struct {
			ID        int    `db:"id"`
			FirstName string `db:"firstName"`
			LastName  string `db:"lastName"`
			Email     string `db:"email"`
			Role      string `db:"role"`
			Password  string `db:"password"`
		}

		err := db.Get(&user, "SELECT id, firstName, lastName, email, password, role FROM accounts WHERE email = ?", account.Email)
		if err != nil || account.Password != user.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Admins only."})
			return
		}

		token, err := util.GenerateJWT(user.ID, user.Email, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   "Login successful",
			"token":     token,
			"id":        user.ID,
			"role":      user.Role,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
		})
	})
}
