package app

import (
	"comment_tree/internal/handlers"
	"comment_tree/internal/middleware"
	"comment_tree/internal/service"
	"comment_tree/internal/storage"
	"comment_tree/internal/storage/postgres"
	"log"

	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/ginext"
)

func Run() {
	cfg := config.New()
	err := cfg.LoadEnvFiles(".env")
	if err != nil {
		log.Fatalf("[app] error of loading cfg: %v", err)
	}
	cfg.EnableEnv("")

	databaseURI := cfg.GetString("DATABASE_URI")
	serverAddr := cfg.GetString("SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = ":8080" //дефолтный порт на случай непрочтения
	}

	postgresStore, err := postgres.New(databaseURI)
	if err != nil {
		log.Fatalf("[app]failed to connect to PG DB: %v", err)
	}
	defer postgresStore.Close()
	log.Println("[app] Connected to Postgres successfully")

	store, err := storage.New(postgresStore)
	if err != nil {
		log.Fatalf("[app] failed to init unified storage: %v", err)
	}
	log.Println("[app] Unified storage initialized successfully")

	serviceLayer := service.New(store)

	engine := ginext.New("release")

	engine.Use(middleware.LoggerMiddleware())
	engine.Use(middleware.CORSMiddleware())

	router := handlers.New(engine, serviceLayer, serviceLayer, serviceLayer)
	router.Routes()

	log.Printf("[app] starting server on %s", serverAddr)
	err = engine.Run(serverAddr)
	if err != nil {
		log.Fatalf("[app] server failed: %v", err)
	}
}
