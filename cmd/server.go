package cmd

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"interview-backend/cmd/flags"
	"interview-backend/internal/client"
	"interview-backend/internal/conf"
	"interview-backend/internal/server"
	"interview-backend/util"
	"net/http"
	"os"
	"os/signal"
	"sync"
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
	if !flags.Dev {
		gin.SetMode(gin.ReleaseMode)
		config.AllowMethods = []string{"GET", "POST", "DELETE"}
	} else {
		gin.SetMode(gin.DebugMode)
		config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
		config.AllowAllOrigins = true
	}
	r := gin.New()
	r.Use(
		server.RouterRecovery(),
		gin.LoggerWithWriter(log.StandardLogger().Out),
		gin.RecoveryWithWriter(log.StandardLogger().Out),
		cors.New(config))
	server.Init(r)
	client.Init()

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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second) //创建一个context，用于优雅地关闭gin
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := httpServer.Shutdown(ctx); err != nil {
				util.Log.Fatalf("Failed to shutdown server: %v", err)
			}
		}()
	}()
	wg.Wait()
	util.Log.Infof("Server gracefully stopped")
}

func init() {
	RootCmd.AddCommand(ServerCmd)
}
