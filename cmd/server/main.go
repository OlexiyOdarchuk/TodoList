package server

import (
	"log/slog"
	"os"
	"time"

	"todolist/internal/repository"
	"todolist/internal/routes"
	"todolist/internal/service"
	"todolist/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env file found, relying on environment variables")
	}
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := "postgres://" + dbUser + ":" + dbPass + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	r := gin.Default()
	port := os.Getenv("PORT")
	if port == "" {
		slog.Error("PORT environment variable not set")
		os.Exit(1)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		slog.Error("JWT_SECRET environment variable not set")
		os.Exit(1)
	}

	jwtManager := utils.NewJWTManager(jwtSecret, time.Hour*24*7)
	userRepo := repository.NewUsersRepository(db)
	userService := service.NewUserService(userRepo, jwtManager)
	userRouter := routes.NewUserHandler(userService)

	r.POST("/register", userRouter.Register)

	slog.Info("Starting server", "port", port)
	err = r.Run(":" + port)
	if err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
