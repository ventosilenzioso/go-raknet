package main

import (
	"os"
	"os/signal"
	"samp-server-go/core/gamemode"
	"samp-server-go/pkg/logger"
	"samp-server-go/source/server"
	"syscall"
	"time"
)

const (
	VERSION = "1.0.0"
	AUTHOR  = "GO:SA-MP"
)

func main() {
	logger.Banner("RakNet Server - Built with Go", VERSION)
	
	// Load configuration
	config := loadConfig()
	
	// Initialize gamemode
	gm := gamemode.NewFreeroamGamemode()
	logger.Success("Gamemode initialized: Freeroam")
	
	// Create server instance
	srv := server.NewServer(config.Host, config.Port, config.MaxPlayers)
	srv.ServerName = config.ServerName
	srv.GameMode = config.GameMode
	srv.Language = config.Language
	srv.Weather = config.Weather
	srv.WorldTime = config.WorldTime
	srv.MapName = config.MapName
	srv.WebURL = config.WebURL
	
	logger.Info("Server Version: %s", VERSION)
	logger.Info("Starting server on %s:%d", srv.Host, srv.Port)
	logger.Info("Max players: %d", srv.MaxPlayers)
	logger.Info("Server name: %s", srv.ServerName)
	logger.Info("Game mode: %s", srv.GameMode)
	logger.Info("Language: %s", srv.Language)
	logger.Info("Weather: %d", srv.Weather)
	logger.Info("World time: %d:00", srv.WorldTime)
	logger.Info("Map name: %s", srv.MapName)
	logger.Info("Web URL: %s", srv.WebURL)
	logger.Success("Configuration loaded successfully")
	
	// Setup event handlers
	setupGamemodeEvents(srv, gm)
	
	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	
	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := srv.Start(); err != nil {
			errChan <- err
		}
	}()
	
	// Wait for shutdown signal or error
	select {
	case err := <-errChan:
		logger.Fatal("Server error: %v", err)
	case sig := <-sigChan:
		logger.Warn("Received signal: %v", sig)
		logger.Info("Shutting down gracefully...")
		
		// Stop server
		srv.Stop()
		
		// Wait a bit for cleanup
		time.Sleep(1 * time.Second)
		
		logger.Success("Server stopped")
		os.Exit(0)
	}
}

type Config struct {
	Host       string
	Port       int
	MaxPlayers int
	ServerName string
	GameMode   string
	Language   string
	Weather    int
	WorldTime  int
	MapName    string
	WebURL     string
}

func loadConfig() Config {
	// Default configuration
	// You can modify these values or load from environment variables
	return Config{
		Host:       "0.0.0.0",
		Port:       7777,
		MaxPlayers: 100,
		ServerName: "RakNet Server [GO]",
		GameMode:   "Freeroam v1.0",
		Language:   "English",
		Weather:    10,
		WorldTime:  12,
		MapName:    "San Andreas",
		WebURL:     "github.com/yourusername/raknet-go",
	}
}

func setupGamemodeEvents(srv *server.Server, gm *gamemode.FreeroamGamemode) {
	// TODO: Wire up gamemode events to server events
	// This will be implemented when server event system is ready
	logger.Success("Gamemode events configured")
}
