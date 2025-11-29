package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func InitDB() error {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: Could not load .env file:", err)
	}

	config := DBConfig{
		Host:     os.Getenv("PG_HOST"),
		Port:     os.Getenv("PG_PORT"),
		User:     os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASSWORD"),
		Database: os.Getenv("PG_DATABASE"),
		SSLMode:  os.Getenv("PG_SSLMODE"),
	}

	// Prefer DATABASE_URL if provided
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
					config.User,
					config.Password,
					config.Host,
					config.Port,
					config.Database,
					config.SSLMode,
				)
	}
	fmt.Print("value of env files", connStr)
	log.Printf("Connecting to database: %s@%s:%s/%s", 
		config.User, config.Host, config.Port, config.Database)
		
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("\n error parsing connection string: %w", err)
	}

	// Configure connection pool
	poolConfig.MaxConns = 20
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	DB, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return fmt.Errorf("\n error creating connection pool: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := DB.Ping(ctx); err != nil {
		return fmt.Errorf("\n error pinging database: %w", err)
	}

	log.Println("\n Successfully connected to PostgreSQL!")
	return nil
}



func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("\n Database connection closed")
	}
}


