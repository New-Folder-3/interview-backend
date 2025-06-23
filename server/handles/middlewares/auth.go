package middlewares

import (
	"github.com/gin-gonic/gin"
	"interview-backend/internal/db"
	"interview-backend/util"
	"net/http"
)

func Auth(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	_, err := db.GetToken(token)
	if err != nil {
		util.ErrorResp(c, "Unauthorized", http.StatusUnauthorized, false)
		c.Abort()
	}
}
