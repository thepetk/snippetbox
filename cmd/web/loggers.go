package main

import "fmt"

func (a *application) Log(logLevel string, msg string) {
	cfgLevel := a.cfg.LogLevel

	if !stringInLogLevels(a.cfg.LogLevel) {
		cfgLevel = "INFO"
	}

	if getLogLevelInt(cfgLevel) <= getLogLevelInt(logLevel) {
		a.logger.Printf(fmt.Sprintf("%s\t%s", logLevel, msg))
	}
}

func getLogLevelInt(key string) int {
	for _, logLevel := range LogLevels {
		if logLevel.key == key {
			return logLevel.value
		}
	}
	return 1
}

func stringInLogLevels(item string) bool {
	for _, logLevel := range LogLevels {
		if logLevel.key == item {
			return true
		}
	}
	return false
}
