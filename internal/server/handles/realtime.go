package handles

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"interview/internal/conf"
	"interview/op/realtime"
	"interview/util"
	"net/http"
	"net/url"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Realtime(c *gin.Context) {
	user := c.GetHeader("User")
	userShould, ok := c.Get("user")
	if !ok || userShould.(string) != user {
		util.ErrorResp(c, "unmatched user", http.StatusUnauthorized)
		return
	}
	frontend, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		util.ErrorResp(c, "websocket upgrade error", http.StatusInternalServerError)
	}
	defer frontend.Close()

	if conf.Conf.API.AliyunAPIKey == "" {
		util.ErrorResp(c, "aliyun api key not set", http.StatusInternalServerError)
		return
	}
	u := url.URL{Scheme: "wss", Host: "dashscope.aliyuncs.com", Path: "/api-ws/v1/realtime"}
	header := http.Header{}
	header.Set("Authorization", "Bearer "+conf.Conf.API.AliyunAPIKey)

	remote, resp, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil || resp.StatusCode != http.StatusSwitchingProtocols {
		util.ErrorPrinter(errors.WithMessage(err, resp.Status))
		util.ErrorResp(c, "cannot connect to aliyun", http.StatusInternalServerError)
		return
	}

	aliDone := make(chan error)
	frontendDone := make(chan error)
	conversationIDChan := make(chan string)

	go realtime.RemoteToFrontend(aliDone, remote, frontend, conversationIDChan, user)
	go realtime.FrontendToRemote(frontendDone, remote, frontend)

	conversationID := <-conversationIDChan

	select {
	case err = <-aliDone:
	case err = <-frontendDone:
	}

	if err != nil {
		util.ErrorPrinter(errors.WithStack(err))
		util.ErrorResp(c, err.Error(), http.StatusInternalServerError)
	} else {
		util.SuccessResp(c, conversationID, "connection closed")
	}
}
