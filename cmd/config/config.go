package config

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"
)

// Env var keys
const (
	SNIPPETBOX_ADDR        = "SNIPPETBOX_ADDR"
	SNIPPETBOX_STATIC      = "SNIPPETBOX_STATIC"
	SNIPPETBOX_DB_PASSWORD = "SNIPPETBOX_DB_PASSWORD"
	SNIPPETBOX_DB_USERNAME = "SNIPPETBOX_DB_USERNAME"
	SNIPPETBOX_DB_NAME     = "SNIPPETBOX_DB_NAME"
)

// Env var fallbacks
const (
	SNIPPETBOX_ADDR_FALLBACK        = ":4000"
	SNIPPETBOX_STATIC_FALLBACK      = "./ui/static/"
	SNIPPETBOX_DB_PASSWORD_FALLBACK = ""
	SNIPPETBOX_DB_USERNAME_FALLBACK = ""
	SNIPPETBOX_DB_NAME_FALLBACK     = "snippetbox"
)

// Config consts
const (
	MAX_OPEN_DATABASE_CONNECTIONS = 5
	MAX_IDLE_DATABASE_CONNECTIONS = 5
	DATABASE_TIMEOUT              = 5 * time.Second
)

type Config struct {
	Addr      string
	StaticDir string
	DB        DBConfig
	SetLimits bool
}

type DBConfig struct {
	DBUser        string
	DBPass        string
	DBSSLDisabled bool
}

func (c *Config) InitConfig() (*sql.DB, error) {
	// Parse all flags from command line
	flag.StringVar(&c.Addr, "addr", "", "HTTP network address")
	flag.StringVar(&c.StaticDir, "static-dir", "", "Path to static assets")
	flag.StringVar(&c.DB.DBPass, "db-pass", "", "Database password")
	flag.StringVar(&c.DB.DBUser, "db-user", "", "Database username")
	flag.StringVar(&c.DB.DBUser, "db-name", "", "Database name")
	flag.BoolVar(&c.DB.DBSSLDisabled, "db-ssl-disabled", true, "Database ssl config")
	flag.BoolVar(&c.SetLimits, "set-limits", false, "Set database pool limits")
	flag.Parse()

	db, err := c.OpenDB()
	return db, err

}
func (c *Config) GetAddr() string {
	return getVar(c.Addr, SNIPPETBOX_ADDR, SNIPPETBOX_ADDR_FALLBACK)
}

func (c *Config) GetStaticDir() string {
	return getVar(c.StaticDir, SNIPPETBOX_STATIC, SNIPPETBOX_STATIC_FALLBACK)
}

func (c *Config) GetDBUser() (string, error) {
	dbUser := getVar(c.DB.DBUser, SNIPPETBOX_DB_USERNAME, SNIPPETBOX_DB_USERNAME_FALLBACK)
	if dbUser == "" {
		return "", fmt.Errorf("Username not set for database connection")
	}
	return dbUser, nil
}

func (c *Config) GetDBPass() (string, error) {
	dbPass := getVar(c.DB.DBPass, SNIPPETBOX_DB_PASSWORD, SNIPPETBOX_DB_PASSWORD_FALLBACK)
	if dbPass == "" {
		return "", fmt.Errorf("Password not set for database connection")
	}
	return dbPass, nil
}

func (c *Config) GetDBName() string {
	return getVar(c.DB.DBUser, SNIPPETBOX_DB_NAME, SNIPPETBOX_DB_NAME_FALLBACK)
}

func (c *Config) GetDBSSLMode() string {
	sslDisabled := "disable"
	if !c.DB.DBSSLDisabled {
		sslDisabled = "enable"
	}
	return sslDisabled
}

func (c *Config) GetDSN() (string, error) {
	dbUser, err := c.GetDBUser()
	if err != nil {
		return "", err
	}

	dbPass, err := c.GetDBPass()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", dbUser, dbPass, c.GetDBName(), c.GetDBSSLMode()), nil
}

func (c *Config) OpenDB() (*sql.DB, error) {
	dsn, err := c.GetDSN()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if c.SetLimits {
		db.SetMaxOpenConns(MAX_OPEN_DATABASE_CONNECTIONS)
		db.SetMaxIdleConns(MAX_IDLE_DATABASE_CONNECTIONS)
	}

	ctx, cancel := context.WithTimeout(context.Background(), DATABASE_TIMEOUT)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getVar(value string, envVar string, fallback string) string {
	if len(value) == 0 {
		return getEnv(envVar, fallback)
	}
	return value
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
