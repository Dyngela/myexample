package main

import (
	"basicapi/api/entity"
	"basicapi/api/handler"
	"basicapi/api/handler/middleware"
	"basicapi/api/handler/validators"
	"basicapi/config"
	"basicapi/docs"
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/graceful"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os/signal"
	"syscall"
)

func main() {
	config.InitConfiguration()

	// Allow graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	router, err := graceful.Default(graceful.WithAddr(fmt.Sprintf(":%s", config.GetConfig().ServerPort)))
	defer stop()
	defer router.Close()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validators.RegisterCustomValidations(v)
	}

	if config.GetConfig().ServerMode == "dev" {
		docs.SwaggerInfo.BasePath = ""
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		gin.ForceConsoleColor()
		gin.SetMode(gin.DebugMode)
		autoMigrateTable()
	}

	initAPI(router)

	if err = router.RunWithContext(ctx); err != nil && !errors.Is(err, context.Canceled) {
		config.Logger.Fatal().Msg(err.Error())
		panic(err)
	}

}

func autoMigrateTable() {
	err := config.DB.AutoMigrate(&entity.Vehicle{})
	if err != nil {
		panic(err)
	}
}

func initAPI(router *graceful.Graceful) {
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.APIKeyMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.RateLimitMiddleware)
	router.Use(middleware.TimeoutMiddleware(10))

	handler.VehicleHandler(router)
}
