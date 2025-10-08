package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/banggibima/be-assignment/config"
	"github.com/banggibima/be-assignment/docs"
	"github.com/banggibima/be-assignment/internal/handlers"
	"github.com/banggibima/be-assignment/internal/repositories"
	"github.com/banggibima/be-assignment/internal/services"
	"github.com/banggibima/be-assignment/pkg/db"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	pool, err := db.Init(cfg)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	productRepo := repositories.NewDatabaseProductRepository(pool)
	orderRepo := repositories.NewDatabaseOrderRepository(pool)
	jobRepo := repositories.NewDatabaseJobRepository(pool)
	transRepo := repositories.NewDatabaseTransactionRepository(pool)
	settleRepo := repositories.NewDatabaseSettlementRepository(pool)

	orderService := services.NewOrderService(pool, orderRepo, productRepo)
	jobService := services.NewJobService(pool, jobRepo, transRepo, settleRepo)

	ctx, cancel := context.WithCancel(context.Background())
	jobService.StartWorkerPool(ctx)

	router := gin.Default()

	docs.SwaggerInfo.Title = "My API"
	docs.SwaggerInfo.Description = "This is a sample Go Gin API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:" + cfg.HTTP.Port
	docs.SwaggerInfo.Schemes = []string{"http"}

	handlers.Register(router)
	handlers.NewOrderHandler(orderService).Register(router)
	handlers.NewJobHandler(jobService).Register(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := ":" + cfg.HTTP.Port

	go func() {
		if err := router.Run(addr); err != nil {
			panic(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	cancel()
}
