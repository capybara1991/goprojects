package httptransport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"auth-service/internal/service"
)

func IssueHandler(s *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Query("user_id")
		userID, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}
		access, refresh, err := s.Issue(userID, c.GetHeader("User-Agent"), c.ClientIP())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"access": access, "refresh": refresh})
	}
}

func RefreshHandler(s *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct{ Access, Refresh string }
		if c.BindJSON(&req) != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}
		newAccess, newRefresh, err := s.Refresh(req.Access, req.Refresh, c.GetHeader("User-Agent"), c.ClientIP())
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"access": newAccess, "refresh": newRefresh})
	}
}

func WhoamiHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	}
}

func LogoutHandler(s *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetString("userID")
		userID, _ := uuid.Parse(uid)
		_ = s.Logout(userID)
		c.JSON(http.StatusOK, gin.H{"message": "logged out"})
	}
}
