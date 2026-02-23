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
	"github.com/example/ppo/internal/order"
	"github.com/example/ppo/internal/postpurchase"
	"github.com/example/ppo/internal/scheduler"
)

type Server struct {
	Router    *gin.Engine
	Scheduler *scheduler.Scheduler
}

func New(cfg *config.Config, db *gorm.DB, logger *slog.Logger) *Server {
	httpClient := &http.Client{Timeout: 10 * time.Second}

	// --- external clients ---
	lmsClient := lms.NewClient(cfg.LMSBaseURL, httpClient)
	pspClient := psp.NewClient(cfg.PSPBaseURL, httpClient)
	prodClient := product.NewClient(cfg.ProductBaseURL, httpClient)

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
