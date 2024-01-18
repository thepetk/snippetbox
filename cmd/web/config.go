package config

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
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

type ConfigManager struct{}

func (c *ConfigManager) GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func (c *ConfigManager) GetConfigVar(flagAddr string, fallback string, envVar string) string {
	if len(flagAddr) == 0 {
		return c.GetEnv(envVar, fallback)
	}
	return flagAddr
}

func (c *ConfigManager) GetDSN(dbUser string, dbPass string, dbName string, dbSSLDisabled bool) string {
	sslDisabled := "disable"
	if !dbSSLDisabled {
		sslDisabled = "enable"
	}
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", dbUser, dbPass, dbName, sslDisabled)
}

func (c *ConfigManager) OpenDB(dsn string, setLimits bool) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if setLimits {
		fmt.Println("setting limits")
		db.SetMaxOpenConns(5)
		db.SetMaxIdleConns(5)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
