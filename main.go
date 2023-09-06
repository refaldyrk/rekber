package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rekber/config"
	"rekber/handler"
	"rekber/middleware"
	"rekber/repository"
	"rekber/service"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	//Set Context
	ctx := context.Background()
	//Init Viper
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	//Init Config
	mongodb := config.ConnectMongo(ctx)
	err = mongodb.Ping(1000)
	if err != nil {
		panic(err)
	}

	//Init Database
	DB := mongodb.Database(viper.GetString("DATABASE_NAME"))

	//=================> Repository
	userRepo := repository.NewUserRepository(DB)
	orderRepo := repository.NewOrderRepository(DB)

	//=================> Service
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo)
	orderService := service.NewOrderService(orderRepo, userRepo)

	//=================> Handler
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	orderHandler := handler.NewOrderHandler(orderService)

	//Server
	app := gin.Default()

	app.Use(gin.Recovery())
	app.Use(gin.Logger())

	// cors	config
	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = []string{"*"}
	cfg.AllowCredentials = true
	cfg.AllowMethods = []string{"*"}
	cfg.AllowHeaders = []string{"*"}

	app.Use(cors.New(cfg))

	//======================> Route

	app.POST("/register", userHandler.Register)
	app.POST("/login", authHandler.Login)

	//======================> Order Endpoint Group
	orderEndpoint := app.Group("/api/order")
	orderEndpoint.Use(middleware.JWTMiddleware(DB))

	orderEndpoint.POST("/", orderHandler.NewOrder)
	orderEndpoint.GET("/", orderHandler.FindAllOrderByRole)
	orderEndpoint.GET("/:id", orderHandler.GetOrderDetailByOrderID)

	//Init Server
	srv := &http.Server{
		Addr:    ":9090",
		Handler: app,
	}

	// graceful shutdown
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()

	log.Println("Server exiting")
}
