package main

import (
	"embed"
	_ "embed"
	"flag"
	"fmt"
	"informatikaas_crud/config"
	"informatikaas_crud/database"
	"informatikaas_crud/endpoints"
	"log"

	"github.com/gin-gonic/gin"
)

//go:embed schema.sql
var schemaSQLFS embed.FS

func startServer(cfg *config.Config, db *database.DB) {
	router := gin.New()
	router.SetTrustedProxies(nil)
	router.Use(gin.Logger())

	endpoints.Attach(router, &endpoints.EndpointData{
		Cfg: cfg,
		DB:  db,
	})

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func main() {
	configPathFlag := flag.String("config-path", "config.yaml", "path to the configuration file")
	sqlReplFlag := flag.Bool("sql-repl", false, "instead of starting the server, start an sql repl")

	flag.Parse()

	cfg, err := config.Load(*configPathFlag)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	schema, err := schemaSQLFS.ReadFile("schema.sql")
	if err != nil {
		log.Fatalf("Failed to read schema file: %v", err)
	}

	db, err := database.New(cfg, schema)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if *sqlReplFlag {
		if err := db.StartREPL(); err != nil {
			log.Fatalf("Failed to start SQL REPL: %v", err)
		}
	} else {
		startServer(cfg, db)
	}
}
