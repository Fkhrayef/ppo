package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/example/ppo/internal/client/lms"
	"github.com/example/ppo/internal/client/product"
	"github.com/example/ppo/internal/client/psp"
	"github.com/example/ppo/internal/config"
	mw "github.com/example/ppo/internal/middleware"
	"github.com/example/ppo/internal/order"
	"github.com/example/ppo/internal/postpurchase"
	"github.com/example/ppo/internal/scheduler"
)

type Server struct {
	Router    *gin.Engine
	Scheduler *scheduler.Scheduler
}

func New(cfg *config.Config, db *gorm.DB, logger *slog.Logger) *Server {
	// --- external clients ---
	var (
		lmsClient  lms.Client
		pspClient  psp.Client
		prodClient product.Client
	)

	if cfg.UseFakeClients {
		logger.Info("using FAKE external clients — responses are static contract stubs")
		lmsClient = lms.NewFake(logger)
		pspClient = psp.NewFake(logger)
		prodClient = product.NewFake(logger)
	} else {
		httpClient := &http.Client{Timeout: 10 * time.Second}
		lmsClient = lms.NewHTTPClient(cfg.LMSBaseURL, httpClient)
		pspClient = psp.NewHTTPClient(cfg.PSPBaseURL, httpClient)
		prodClient = product.NewHTTPClient(cfg.ProductBaseURL, httpClient)
	}

	// --- repositories ---
	orderRepo := order.NewRepository(db)

	// --- services ---
	orderSvc := order.NewService(orderRepo, lmsClient, pspClient, prodClient, logger)
	postPurchaseSvc := postpurchase.NewService(lmsClient, pspClient, logger)

	// --- handlers ---
	orderHandler := order.NewHandler(orderSvc)
	postPurchaseHandler := postpurchase.NewHandler(postPurchaseSvc)

	// --- gin router ---
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestLogger(logger))
	r.Use(mw.ErrorHandler())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	orderHandler.RegisterRoutes(v1)
	postPurchaseHandler.RegisterRoutes(v1)

	// --- scheduler ---
	sched := scheduler.New(lmsClient, pspClient, orderRepo, logger)

	return &Server{
		Router:    r,
		Scheduler: sched,
	}
}

func requestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", time.Since(start).String(),
		)
	}
}
