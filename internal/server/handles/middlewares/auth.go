package middlewares

import (
	"github.com/gin-gonic/gin"
	"interview-backend/internal/db"
	"interview-backend/util"
	"net/http"
	"strings"
)

func Auth(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	token, ok := strings.CutPrefix(token, "Bearer ")
	if !ok || len(token) == 0 {
		util.ErrorResp(c, "Empty Token", http.StatusUnauthorized, false)
		c.Abort()
		return
	}
	tokenDB, err := db.GetToken(token)
	if err != nil {
		util.ErrorResp(c, "Unauthorized", http.StatusUnauthorized, false)
		c.Abort()
		return
	}
	c.Set("user", tokenDB.UserID)
	c.Next()
}
