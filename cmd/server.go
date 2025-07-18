package cmd

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"interview/cmd/flags"
	"interview/internal/client"
	"interview/internal/conf"
	"interview/internal/server"
	"interview/util"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the server",
	Long:  `Starts the server`,
	Run: func(cmd *cobra.Command, args []string) {
		serverStart()
	},
}

func serverStart() {
	Init()
	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	var r *gin.Engine
	if !flags.Dev {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	} else {
		gin.SetMode(gin.DebugMode)
		r = gin.New()
		config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
		config.AllowAllOrigins = true
		r.Use(cors.New(config))
	}
	r.Use(
		server.RouterRecovery(),
		gin.LoggerWithWriter(log.StandardLogger().Out),
		gin.RecoveryWithWriter(log.StandardLogger().Out))

	server.Init(r)
	client.Init()
	util.S3Init()

	var httpHandler http.Handler = r
	httpBase := fmt.Sprintf("%s:%d", conf.Conf.Schema.Listen, conf.Conf.Schema.Port)
	util.Log.Infof("Starting HTTP server on %s", httpBase)
	httpServer := &http.Server{
		Addr:    httpBase,
		Handler: httpHandler,
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			util.Log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 优雅地退出程序
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	util.Log.Infof("Shutting down server")
	Release()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		util.Log.Fatalf("Failed to shutdown server: %v", err)
	}
	util.Log.Infof("Server gracefully stopped")

}

func init() {
	RootCmd.AddCommand(ServerCmd)
}
