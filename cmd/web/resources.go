package main

import (
	"os"
	"time"
)

const (
	MaxOpenDatabaseConns = 5
	MaxIdleDatabaseConns = 5
	DatabaseTimeout      = 5 * time.Second
)

type ConfigVar struct {
	key      string
	fallback string
}

type LogLevel struct {
	key   string
	value int
}

// Configuration Variables
var SnippetboxAddr = ConfigVar{
	key:      "SNIPPETBOX_ADDR",
	fallback: ":4000",
}

var SnippetboxStaticFilepath = ConfigVar{
	key:      "SNIPPETBOX_STATIC_FILEPATH",
	fallback: "./ui/static/",
}

var SnippetboxDBPassword = ConfigVar{
	key:      "SNIPPETBOX_DB_PASSWORD",
	fallback: "",
}

var SnippetboxDBUser = ConfigVar{
	key:      "SNIPPETBOX_DB_USERNAME",
	fallback: "",
}

var SnippetboxDBName = ConfigVar{
	key:      "SNIPPETBOX_DB_NAME",
	fallback: "snippetbox",
}

var LogLevels = []LogLevel{
	{
		key:   "DEBUG",
		value: 0,
	},
	{
		key:   "INFO",
		value: 1,
	},
	{
		key:   "ERROR",
		value: 2,
	},
}

func GetConfigVar(value string, configVar ConfigVar) string {
	if len(value) == 0 {
		return getEnv(configVar)
	}
	return value
}

func getEnv(configVar ConfigVar) string {
	value := os.Getenv(configVar.key)
	if len(value) == 0 {
		return configVar.fallback
	}
	return value
}
