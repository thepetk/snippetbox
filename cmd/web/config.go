package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
)

// Config consts

type Config struct {
	Addr          string
	DBUser        string
	DBPass        string
	DBName        string
	DBSSLDisabled bool
	LogLevel      string
	StaticDir     string
	SetLimits     bool
}

func (a *application) InitConfig() (*sql.DB, error) {
	// Parse all flags from command line
	flag.StringVar(&a.cfg.Addr, "addr", "", "HTTP network address")
	flag.StringVar(&a.cfg.LogLevel, "log-level", "INFO", "Logging level (DEBUG - INFO - ERROR)")
	flag.StringVar(&a.cfg.StaticDir, "static-dir", "", "Path to static assets")
	flag.StringVar(&a.cfg.DBPass, "db-pass", "", "Database password")
	flag.StringVar(&a.cfg.DBUser, "db-user", "", "Database username")
	flag.StringVar(&a.cfg.DBUser, "db-name", "", "Database name")
	flag.BoolVar(&a.cfg.DBSSLDisabled, "db-ssl-disabled", true, "Database ssl config")
	flag.BoolVar(&a.cfg.SetLimits, "set-limits", false, "Set database pool limits")
	flag.Parse()

	db, err := a.OpenDB()
	return db, err

}

func (a *application) GetAddr() string {
	return GetConfigVar(a.cfg.Addr, SnippetboxAddr)
}

func (a *application) GetStaticDir() string {
	return GetConfigVar(a.cfg.StaticDir, SnippetboxStaticFilepath)
}

func (a *application) GetDBUser() (string, error) {
	dbUser := GetConfigVar(a.cfg.DBUser, SnippetboxDBUser)
	if dbUser == "" {
		return "", fmt.Errorf("Username not set for database connection")
	}
	return dbUser, nil
}

func (a *application) GetDBPass() (string, error) {
	dbPass := GetConfigVar(a.cfg.DBPass, SnippetboxDBPassword)
	if dbPass == "" {
		return "", fmt.Errorf("Password not set for database connection")
	}
	return dbPass, nil
}

func (a *application) GetDBName() string {
	return GetConfigVar(a.cfg.DBName, SnippetboxDBName)

}

func (a *application) GetDBSSLMode() string {
	sslDisabled := "disable"
	if !a.cfg.DBSSLDisabled {
		sslDisabled = "enable"
	}
	return sslDisabled
}

func (a *application) GetDSN() (string, error) {
	dbUser, err := a.GetDBUser()
	if err != nil {
		return "", err
	}

	dbPass, err := a.GetDBPass()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", dbUser, dbPass, a.GetDBName(), a.GetDBSSLMode()), nil
}

func (a *application) OpenDB() (*sql.DB, error) {
	dsn, err := a.GetDSN()
	if err != nil {
		return nil, err
	}

	a.Log("DEBUG", "Opening postgres connection")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if a.cfg.SetLimits {
		a.Log("DEBUG", "Setting limits")
		db.SetMaxOpenConns(MaxOpenDatabaseConns)
		db.SetMaxIdleConns(MaxOpenDatabaseConns)
	}

	ctx, cancel := context.WithTimeout(context.Background(), DatabaseTimeout)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
