package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"

	"github.com/libretro/netplay-lobby-server-go/controller"
	"github.com/libretro/netplay-lobby-server-go/domain"
	"github.com/libretro/netplay-lobby-server-go/model"
	"github.com/libretro/netplay-lobby-server-go/model/entity"
	"github.com/libretro/netplay-lobby-server-go/model/repository"
)

func main() {
	var verbose = flag.Bool("v", false, "verbose logging")
	flag.Parse()

	server := echo.New()
	server.HideBanner = true
	server.Server.ReadTimeout = 5 * time.Second
	server.Server.WriteTimeout = 5 * time.Second

	config, err := readConfig()
	if err != nil {
		server.Logger.Fatalf("Can't get configuration values: %v", err)
	}

	// Init domain logic and model
	db, err := initDatabase(config.Database.Type, config.Database.Connection)
	if err != nil {
		server.Logger.Fatalf("Can't initialize database: %v", err)
	}
	db = db.AutoMigrate(&entity.Session{})

	if *verbose {
		server.Logger.SetLevel(log.INFO)
	} else {
		server.Logger.SetLevel(log.WARN)
	}

	sessionDomain, err := initDomain(db, config)
	if err != nil {
		server.Logger.Fatalf("Can't initialize domain logic: %v", err)
	}

	sessionCotroller := controller.NewSessionController(sessionDomain)

	// Start the cleanup job to purge old sessions
	go func() {
		for true {
			err := sessionDomain.PurgeOld()
			if err != nil {
				server.Logger.Fatalf("Can't purge old sessions: %v", err)
			}
			time.Sleep(2 * time.Minute)
		}
	}()

	// Server setup
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(middleware.BodyLimit("64K"))

	// Set the routes and prerender templates
	sessionCotroller.RegisterRoutes(server)
	if err = sessionCotroller.PrerenderTemplates(server, "./web/templates/*.html"); err != nil {
		server.Logger.Fatalf("Can't prerender templates: %v", err)
	}

	// Start serving
	server.Logger.Fatal(server.Start(config.Server.Address))
}

func readConfig() (*Config, error) {
	viper := viper.New()
	viper.SetConfigName("lobby")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/lobby")
	viper.AddConfigPath("$HOME/.lobby")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Can't read configuration file: %w", err)
	}
	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("Can't unmarshal configuration file: %w", err)
	}
	return &conf, nil
}

func initDatabase(databaseType string, connectionString string) (*gorm.DB, error) {
	switch databaseType {
	case "mysql":
		return model.GetMysqlDB(connectionString)
	case "postgres":
		return model.GetPostgreDB(connectionString)
	case "sqlite":
		return model.GetSqliteDB(connectionString)
	}

	return nil, fmt.Errorf("Unknown database type in configuration: %s", databaseType)
}

func initDomain(db *gorm.DB, config *Config) (*domain.SessionDomain, error) {
	repo := repository.NewSessionRepository(db)
	geo2Domain, err := domain.NewGeoIP2Domain(config.Server.GeoLite2Path)
	if err != nil {
		return nil, fmt.Errorf("Can't intialize geolite2 database: %w", err)
	}
	validationDomain, err := domain.NewValidationDomain(config.Blacklist.Strings, config.Blacklist.IPs)
	if err != nil {
		return nil, fmt.Errorf("Can't intialize validation domain: %w", err)
	}
	mitmDomain := domain.NewMitmDomain(config.Relay)
	sessionDomain := domain.NewSessionDomain(repo, geo2Domain, validationDomain, mitmDomain)

	return sessionDomain, nil
}
