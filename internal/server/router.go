package server

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"interview-backend/cmd/flags"
	"interview-backend/internal/server/handles"
	"interview-backend/internal/server/handles/middlewares"
	"interview-backend/public"
	"interview-backend/util"
	"io/fs"
	"net/http"
	"path/filepath"
)

func RouterRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				util.ErrorResp(c, "Internal Server Error", http.StatusInternalServerError, false)
				c.Abort()
			}
		}()
		c.Next()
	}
}

func Init(e *gin.Engine) {
	rootFS, _ := fs.Sub(public.StaticFS, "dist")
	assetsFS, _ := fs.Sub(public.StaticFS, "dist/assets")
	tpl := template.Must(template.ParseFS(rootFS, "*.html"))
	e.StaticFS("/assets", http.FS(assetsFS))
	e.SetHTMLTemplate(tpl)

	api := e.Group("/api")
	auth := api.Group("/auth")
	upload := api.Group("/upload")
	user := api.Group("/user")
	conversation := api.Group("/conversation")
	message := api.Group("/message")
	download := api.Group("/download")

	auth.POST("/login/pwd", handles.LoginPwdHandle)
	auth.POST("/login/token", handles.LoginTokenHandle)
	auth.POST("/register", handles.RegisterHandle)
	auth.POST("/logout", handles.LogoutHandle)
	upload.POST("/image", middlewares.Auth, middlewares.FileChecker("image"), handles.FileSaver("image"))
	upload.POST("/audio", middlewares.Auth, middlewares.FileChecker("audio"), handles.FileSaver("audio"))
	user.GET("/get", middlewares.Auth, handles.GetUser)
	user.POST("/update/info", middlewares.Auth, handles.UpdateUserInfo)
	user.POST("/update/pwd", middlewares.Auth, handles.UpdateUserPwd)
	conversation.GET("/all", middlewares.Auth, handles.GetAllConversations)
	conversation.GET("/get", middlewares.Auth, handles.GetConversation)
	conversation.POST("/new", middlewares.Auth, handles.CreateConversation)
	conversation.DELETE("/del", middlewares.Auth, handles.DeleteConversation)
	message.POST("/new", middlewares.Auth, handles.CreateMessage)
	message.DELETE("/del", middlewares.Auth, handles.DeleteMessage)
	download.Static("/image", filepath.Join(flags.DataDir, "image"))
	download.Static("/audio", filepath.Join(flags.DataDir, "audio"))

	e.NoRoute(middlewares.APINoRoute, StaticHandler)
}

func StaticHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
