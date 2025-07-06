package util

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func ErrorResp(c *gin.Context, err string, code int) {
	c.JSON(code, gin.H{
		"error": err,
	})
}

func SuccessResp(c *gin.Context, data interface{}, message string) {
	c.JSON(200, gin.H{
		"message": message,
		"data":    data,
	})
}

func ErrorRespWS(ws *websocket.Conn, message string, code int) {
	resp := map[string]interface{}{
		"error": message,
		"code":  code,
	}
	payload, _ := json.Marshal(resp)
	ws.WriteMessage(websocket.TextMessage, payload)
}

func SuccessRespWS(ws *websocket.Conn, data interface{}, message string) {
	resp := map[string]interface{}{
		"message": message,
		"data":    data,
	}
	payload, _ := json.Marshal(resp)
	ws.WriteMessage(websocket.TextMessage, payload)
}
