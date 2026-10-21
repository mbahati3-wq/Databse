package database

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func get_Database_defalt_config() DbConfig {
	port, available := rtx_port(5432)
	homeDir, _ := os.UserHomeDir()
	workingDir := ""
	if homeDir != "" {
		workingDir = filepath.Join(homeDir, ".pgdb")
	}

	if !available {
		port = 0
	}

	return DbConfig{
		DbName:       "postgres",
		DbUserName:   os.Getenv("USER"),
		DbPassword:   os.Getenv("PGPASSWORD"),
		DbWorkingDir: workingDir,
		DbPort:       uint32(port),
		DbSSLMode:    "disable",
		DbTimeZone:   "Asia/Kolkata",
	}
}

func rtx_port(target_port uint16) (uint16, bool) {
	if target_port < 1024 {
		if target_port == 0 {
			target_port = 5432
		} else {
			target_port = 1024
		}
	}

	for port := uint32(target_port); port <= 65535; port++ {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			continue
		}

		if err := listener.Close(); err == nil {
			return uint16(port), true
		}
	}

	return 0, false
}

func port_in_use(port uint16) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return true
	}

	return listener.Close() != nil
}

func Validate_db_config(db_config DbConfig) Error {
	switch {
	case strings.TrimSpace(db_config.DbName) == "":
		return invalidConfig("Database name is missing.", "Set DbName to the name of the database to create.")
	case strings.TrimSpace(db_config.DbUserName) == "":
		return invalidConfig("Database username is missing.", "Set DbUserName to a valid database user.")
	case db_config.DbPassword == "":
		return invalidConfig("Database password is missing.", "Set DbPassword or provide it through a secure secret source.")
	case strings.TrimSpace(db_config.DbWorkingDir) == "":
		return invalidConfig("Database working directory is missing.", "Set DbWorkingDir to a writable directory.")
	case db_config.DbPort == 0:
		return invalidConfig("Database port is missing.", "Set DbPort to an available port between 1024 and 65535.")
	case db_config.DbPort > 65535:
		return invalidConfig("Database port is out of range.", "Set DbPort to a value between 1024 and 65535.")
	case port_in_use(uint16(db_config.DbPort)):
		return invalidConfig(fmt.Sprintf("Database port %d is unavailable.", db_config.DbPort), "Choose a port that is not already in use and that this process can bind.")
	case strings.TrimSpace(db_config.DbSSLMode) == "":
		return invalidConfig("Database SSL mode is missing.", "Set DbSSLMode to the SSL mode required by your database.")
	case strings.TrimSpace(db_config.DbTimeZone) == "":
		return invalidConfig("Database time zone is missing.", "Set DbTimeZone to a valid IANA time zone, such as UTC.")
	default:
		if _, err := time.LoadLocation(db_config.DbTimeZone); err != nil {
			return invalidConfig(fmt.Sprintf("Database time zone %q is invalid.", db_config.DbTimeZone), "Set DbTimeZone to a valid IANA time zone, such as UTC.")
		}

		return Error{
			Sucess:        true,
			Message:       "Database configuration is valid.",
			ExitCode:      0,
			Recomendation: "No action required.",
		}
	}
}

func invalidConfig(message, recommendation string) Error {
	return Error{
		Sucess:        false,
		Message:       message,
		ExitCode:      1,
		Recomendation: recommendation,
	}
}
