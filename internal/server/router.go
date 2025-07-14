package server

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"interview/cmd/flags"
	"interview/internal/server/handles"
	"interview/internal/server/handles/middlewares"
	"interview/public"
	"interview/util"
	"io/fs"
	"net/http"
	"path/filepath"
)

func RouterRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				util.ErrorResp(c, "Internal Server Error", http.StatusInternalServerError)
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
	e.SetHTMLTemplate(tpl)

	api := e.Group("/api")
	auth := api.Group("/auth")
	upload := api.Group("/upload")
	user := api.Group("/user")
	conversation := api.Group("/conversation")
	message := api.Group("/message")
	download := api.Group("/download")
	interview := api.Group("/interview")

	// Use AuthRequest
	auth.POST("/login/pwd", handles.LoginPwdHandle)
	auth.POST("/login/token", handles.LoginTokenHandle)
	auth.POST("/register", handles.RegisterHandle)
	auth.POST("/logout", handles.LogoutHandle)

	// Use UserRequest
	user.GET("/get", middlewares.Auth, handles.GetUser)
	user.POST("/update/info", middlewares.Auth, handles.UpdateUserInfo)
	user.POST("/update/pwd", middlewares.Auth, handles.UpdateUserPwd)

	// Use ConversationRequest
	conversation.GET("/all", middlewares.Auth, handles.GetAllConversations)
	conversation.GET("/get", middlewares.Auth, handles.GetConversation)
	conversation.POST("/new", middlewares.Auth, handles.CreateConversation)
	conversation.GET("/realtime", middlewares.Auth, handles.Realtime)
	conversation.DELETE("/del", middlewares.Auth, handles.DeleteConversation)
	conversation.POST("/combine", middlewares.Auth, handles.CombineConversations)
	message.POST("/new", middlewares.Auth, handles.NewMessage)
	message.DELETE("/del", middlewares.Auth, handles.DeleteMessage)

	// Use InterviewRequest
	interview.POST("/new", middlewares.Auth, handles.NewInterview)
	interview.GET("/get", middlewares.Auth, handles.GetInterview)
	interview.GET("/all", middlewares.Auth, handles.GetAllInterviews)
	interview.DELETE("/del", middlewares.Auth, handles.DeleteInterview)
	interview.POST("/update", middlewares.Auth, handles.UpdateInterview)
	interview.GET("/keywords", middlewares.Auth, handles.GetKeywords)
	interview.GET("/dimension", middlewares.Auth, handles.GetDimension)
	interview.GET("/get-comment", middlewares.Auth, handles.GetComment)
	interview.POST("/add-comment", middlewares.Auth, handles.AddComment)

	// File upload and download
	upload.POST("/image", middlewares.Auth, middlewares.FileChecker("image"), handles.FileSaver("image"))
	upload.POST("/audio", middlewares.Auth, middlewares.FileChecker("audio"), handles.FileSaver("audio"))
	upload.POST("/audio64", middlewares.Auth, handles.Audio64Saver)
	upload.POST("/video", middlewares.Auth, middlewares.FileChecker("video"), handles.FileSaver("video"))
	download.Static("/image", filepath.Join(flags.DataDir, "image"))
	download.Static("/audio", filepath.Join(flags.DataDir, "audio"))
	download.Static("/video", filepath.Join(flags.DataDir, "audio"))

	// Static files
	e.StaticFS("/assets", http.FS(assetsFS))
	e.NoRoute(middlewares.APINoRoute, StaticHandler)
}

func StaticHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
