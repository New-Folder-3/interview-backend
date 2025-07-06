package handles

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"interview/internal/conf"
	"interview/op/realtime"
	"interview/util"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 64,
	WriteBufferSize: 1024 * 64,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域请求
	},
}

func Realtime(c *gin.Context) {
	user := c.GetHeader("User")
	userShould, ok := c.Get("user")
	if !ok || userShould.(string) != user {
		util.ErrorResp(c, "unmatched user", http.StatusUnauthorized)
		return
	}
	model := c.Query("model")
	frontend, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		util.ErrorRespWS(frontend, "websocket upgrade error", http.StatusInternalServerError)
	}
	defer frontend.Close()

	if conf.Conf.API.AliyunAPIKey == "" {
		util.ErrorRespWS(frontend, "aliyun api key not set", http.StatusInternalServerError)
		return
	}
	header := http.Header{}
	header.Set("Authorization", "Bearer "+conf.Conf.API.AliyunAPIKey)

	remote, resp, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s?model=%s", conf.AliWSUrl, model), header)
	if err != nil || resp.StatusCode != http.StatusSwitchingProtocols {
		var status string
		if resp != nil {
			status = resp.Status
		} else {
			status = "none"
		}
		util.ErrorPrinter(errors.WithMessage(err, status))
		util.ErrorRespWS(frontend, "cannot connect to aliyun", http.StatusInternalServerError)
		return
	}

	aliDone := make(chan error)
	frontendDone := make(chan error)
	conversationIDChan := make(chan string)

	go realtime.RemoteToFrontend(aliDone, remote, frontend, conversationIDChan, user)
	conversationID := <-conversationIDChan

	go realtime.FrontendToRemote(frontendDone, remote, frontend, conversationID)

	select {
	case err = <-aliDone:
		logrus.Println("disconnected from remote")
	case err = <-frontendDone:
		logrus.Println("disconnected from frontend")
	}

	if err != nil {
		util.ErrorPrinter(errors.WithStack(err))
		util.ErrorRespWS(frontend, err.Error(), http.StatusInternalServerError)
	} else {
		util.SuccessRespWS(frontend, conversationID, "connection closed")
	}
}
