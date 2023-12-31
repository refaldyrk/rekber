package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rekber/config"
	"rekber/handler"
	"rekber/middleware"
	"rekber/repository"
	"rekber/service"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

	//Set Flag
	typeRunFlag := flag.String("type", "", "dev or docker")
	flag.Parse()

	typeRun := *typeRunFlag
	//Set Time
	startServerTime := time.Now()
	//Set Context
	ctx := context.Background()
	//Init Viper
	if typeRun == "dev" {
		viper.SetConfigFile(".env")
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
	} else if typeRun == "docker" {
		allEnviron := os.Environ()
		for _, v := range allEnviron {
			arrEnv := strings.Split(v, "=")
			viper.Set(strings.ToUpper(arrEnv[0]), arrEnv[1])
		}
	} else {
		panic("no such config")
	}

	//Init Config
	mongodb := config.ConnectMongo(ctx)
	err := mongodb.Ping(1000)
	if err != nil {
		panic(err)
	}

	//Init Database
	DB := mongodb.Database(viper.GetString("DATABASE_NAME"))

	//=================> Repository
	userRepo := repository.NewUserRepository(DB)
	authRepo := repository.NewAuthRepository(DB)
	orderRepo := repository.NewOrderRepository(DB)
	paymentRepo := repository.NewPaymentRepository(DB)
	balanceRepo := repository.NewBalanceRepository(DB)

	//=================> Service
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, authRepo)
	orderService := service.NewOrderService(orderRepo, userRepo)
	balanceService := service.NewBalanceService(balanceRepo)
	paymentService := service.NewPaymentService(userRepo, orderRepo, paymentRepo, balanceService)

	//=================> Handler
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	orderHandler := handler.NewOrderHandler(orderService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	balanceHandler := handler.NewBalanceHandler(balanceService)

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
	app.POST("/v2/login", authHandler.LoginV2Register)
	app.POST("/v2/login/version/2/:codelink", authHandler.LoginV2)

	//======================> Myself Endpoint
	myselfEndpoint := app.Group("/api/myself")
	myselfEndpoint.Use(middleware.JWTMiddleware(DB, authRepo))

	myselfEndpoint.GET("/", userHandler.MySelf)

	//=======================> Logout Endpoint
	logoutEndpoint := app.Group("/api/logout")
	logoutEndpoint.Use(middleware.JWTMiddleware(DB, authRepo))

	logoutEndpoint.DELETE("/", authHandler.Logout)
	logoutEndpoint.DELETE("/remote/:id", authHandler.RemoteLogout)

	//========================> Login Endpoint
	loginEndpoint := app.Group("/api/authy/login")
	loginEndpoint.Use(middleware.JWTMiddleware(DB, authRepo))

	loginEndpoint.GET("/", authHandler.FindAllLogin)
	loginEndpoint.GET("/count", authHandler.CountLoginData)

	//======================> Order Endpoint Group
	orderEndpoint := app.Group("/api/order")
	orderEndpoint.Use(middleware.JWTMiddleware(DB, authRepo))

	orderEndpoint.GET("/", orderHandler.FindAllOrderByRole)
	orderEndpoint.GET("/:id", orderHandler.GetOrderDetailByOrderID)
	orderEndpoint.GET("/status/:status", orderHandler.GetAllOrderByStatus)
	orderEndpoint.POST("/", orderHandler.NewOrder)
	orderEndpoint.PATCH("/cancel/:id", orderHandler.SetCancelStatusByOrderID)
	orderEndpoint.PATCH("/success/:id", orderHandler.SetSuccessByBuyer)

	//=================> Payment Endpoint Group
	paymentEndpoint := app.Group("/api/payment")
	paymentEndpoint.Use(middleware.JWTMiddleware(DB, authRepo))

	paymentEndpoint.POST("/:id", paymentHandler.NewPayment)

	//=================> Notification Payment Endpoint Group
	paymentNotificationEndpoint := app.Group("/notification/3rd/midtrans")

	paymentNotificationEndpoint.POST("/notification", paymentHandler.NotificationPayment)

	//===================> Balance Endpoint
	balanceEndpoint := app.Group("/api/balance")
	balanceEndpoint.Use(middleware.JWTMiddleware(DB, authRepo))

	balanceEndpoint.GET("/", balanceHandler.FindAllOrderByUserID)
	balanceEndpoint.GET("/:id", balanceHandler.FindDetailBalance)

	//===================> Admin Endpoint Group
	adminEnpoint := app.Group("/admin/api/v1")
	adminEnpoint.Use(middleware.JWTMiddleware(DB, authRepo))

	adminEnpoint.GET("/", middleware.IsAdmin(), func(c *gin.Context) {
		c.String(http.StatusOK, "You Are Admin")
	})

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
	log.Println("Shutdown Server ... ", time.Since(startServerTime).Seconds(), " s")

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()

	log.Println("Server exiting")
}
