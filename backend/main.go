package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	// crypto/tls is necessary for the HTTPS server
	// "crypto/tls"

	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/alexbeattie/golangone/config"
	"github.com/alexbeattie/golangone/handlers"
	"github.com/alexbeattie/golangone/models"
	"github.com/alexbeattie/golangone/services"
	//  "golang.org/x/crypto/acme/autocert"
	// "golang.org/x/crypto/acme/autocert/autocert"
)

func initLogger() (*os.File, error) {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}
	logFile := filepath.Join("logs", fmt.Sprintf("app_%s.log", time.Now().Format("2006-01-02")))
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	log.SetOutput(io.MultiWriter(f, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	return f, nil
}

func loadEnvFile() {
	envFile := ".env.local" // Default to local
	if os.Getenv("APP_ENV") == "production" {
		envFile = ".env.production"
	}
	log.Printf("[LOAD_ENV] Time: %s | File: %s", time.Now().Format("2006-01-02 15:04:05"), envFile)
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("[LOAD_ENV] Warning: Could not load %s file: %v", envFile, err)
	} else {
		log.Printf("[LOAD_ENV] Successfully loaded environment variables from file: %s", envFile)
	}
}

func initDB(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("database connection string (DSN) is empty")
	}
	log.Printf("Attempting to connect to database...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	log.Printf("Successfully connected to database")
	if err := db.AutoMigrate(&models.UserPreferences{}); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Printf("Database migrations completed successfully")
	return db, nil
}

func setupRouter(handler *handlers.Handler) *gin.Engine {
	r := gin.Default()

	// corsOrigins := strings.Split(os.Getenv("CORS_ORIGINS"), ",")
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173", // Vite default dev server
			"http://localhost:4173", // Vite preview
			"http://127.0.0.1:5173", // Alternative localhost
			"http://localhost:8080",
			"http://localhost:8081",
			// "http://34.207.185.237:8080",
			// "http://34.207.185.237",
			// "https://34.207.185.237",
			// "http://34.207.185.237:8080",
			// "https://34.207.185.237:443",
			"https://onestepgpsdemo.com",
			"http://onestepgpsdemo.com",
			"https://www.onestepgpsdemo.com",
			"http://www.onestepgpsdemo.com",

			// Optional if you have another local server
			// Use carefully, preferably only in development
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept",
			"Authorization",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/assets", "./dist/assets")
	r.Static("/js", "./dist/js")
	r.Static("/css", "./dist/css")
	r.StaticFile("/favicon.ico", "./dist/favicon.ico")

	// API routes
	api := r.Group("/api/v1")
	{
		api.GET("/preferences/:userId", handler.GetUserPreferences)
		api.PUT("/preferences/:userId", handler.UpdateUserPreferences)
		api.GET("/devices", handler.GetDevices)
	}

	v3 := r.Group("/v3/api")
	{
		v3.GET("/device-info", handler.GetDeviceInfo)
		v3.GET("/route/drive-stop", handler.GetDriveStopRoute)
	}

	// Handle the root route and any unmatched routes with the Vue app
	r.GET("/", func(c *gin.Context) {
		c.File("./dist/index.html")
	})

	r.NoRoute(func(c *gin.Context) {
		c.File("./dist/index.html")
	})

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[GIN] %s | %d | %s | %s | %s %s | Host: %s | Origin: %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
			param.Request.Host, // This explicitly logs the domain hitting the server
			param.Request.Header.Get("Origin"),
		)
	}))
	r.Use(func(c *gin.Context) {
		log.Printf("Request Headers: %v", c.Request.Header)
		c.Next()
	})
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		log.Printf("Incoming request from origin: %s", origin)
		c.Next()
	})
	return r
}
func getCertificate() (string, string) {
	return "./certs/cert.pem", "./certs/key.pem"
}
func main() {
	logFile, err := initLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logFile.Close()

	log.Println("Starting application...")

	loadEnvFile()

	requiredEnvVars := []string{"ONESTEPGPS_API_KEY", "GOOGLE_MAPS_API_KEY", "DSN"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Required environment variable %s is not set", envVar)
		}
	}
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("DSN is not set")
	}

	cfg := &config.Config{
		OneStepGPSAPIKey: os.Getenv("ONESTEPGPS_API_KEY"),
		GoogleMapsAPIKey: os.Getenv("GOOGLE_MAPS_API_KEY"),
		DSN:              os.Getenv("DSN"),
	}

	db, err := initDB(cfg.DSN)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	service := services.NewService(db, cfg)
	handler := handlers.NewHandler(service, db)

	r := setupRouter(handler)

	// Use autocert for automatic certificate management

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
		log.Printf("APP_PORT not set, defaulting to %s", port)
	}
	address := fmt.Sprintf("0.0.0.0:%s", port)

	server := &http.Server{
		Addr:    address,
		Handler: r,
	}
	// this runs server locally on HTTP
	go func() {
		log.Printf("Starting server on %s", address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	/// These are the changes to the main function that allow the server to run on both HTTP and HTTPS uncomment the code below to run on HTTPS and comment out the code below to run on HTTP only on an EC2 instance
	// if err != nil {
	// 	log.Fatalf("Failed to load SSL/TLS certificate: %v", err)
	// }
	// httpsServer := &http.Server{
	// 	 Addr:    ":443",  // Bind to all interfaces
	//     Handler: r,
	//     TLSConfig: &tls.Config{
	//         MinVersion:               tls.VersionTLS12,
	//         PreferServerCipherSuites: true,
	//     },
	// }

	// Configure HTTP server to redirect to HTTPS

	// httpServer := &http.Server{
	//     Addr: ":80",  // Explicitly bind to all interfaces
	//     Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//         target := "https://" + r.Host + r.RequestURI
	//         http.Redirect(w, r, target, http.StatusMovedPermanently)
	//     }),
	// }

	// Start HTTP server (in a goroutine)

	// go func() {
	//     log.Printf("HTTP Server starting on port 80 (redirecting to HTTPS)")
	//     if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//         log.Printf("HTTP server error: %v", err)
	//     }
	// }()

	// Start HTTPS server (in a goroutine)

	// go func() {
	//     certFile, keyFile := getCertificate()
	//     log.Printf("HTTPS Server starting on port 443")
	//     if err := httpsServer.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
	//         log.Fatalf("Failed to start HTTPS server: %v", err)
	//     }
	// }()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	log.Printf("Received shutdown signal: %v", sig)
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
		log.Println("Database connection closed.")
	}

	log.Println("Server exited properly")

}
