package server

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"interview-backend/cmd/flags"
	"interview-backend/public"
	"interview-backend/server/handles"
	"io/fs"
	"net/http"
)

func Init(e *gin.Engine) {
	rootFS, _ := fs.Sub(public.StaticFS, "dist")
	assetsFS, _ := fs.Sub(public.StaticFS, "dist/assets")
	e.StaticFS("/assets", http.FS(assetsFS))
	tpl := template.Must(template.ParseFS(rootFS, "*.html"))
	e.SetHTMLTemplate(tpl)

	api := e.Group("/api")
	auth := api.Group("/auth")

	auth.POST("/login/pwd", handles.LoginPwdHandle)
	auth.POST("/login/token", handles.LoginTokenHandle)
	auth.POST("/register", handles.RegisterHandle)
	api.POST("/upload/image", handles.PicUploaderHandle)
	api.POST("/upload/audio", handles.AudioUploaderHandle)
	api.Static("/image", flags.DataDir+"image")
	api.Static("/audio", flags.DataDir+"audio")

	e.NoRoute(func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })
}
