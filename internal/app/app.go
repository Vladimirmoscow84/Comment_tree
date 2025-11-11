package app

import (
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

	postgresStore, err := postgres.New(databaseURI)
	if err != nil {
		log.Fatalf("[app]failed to connect to PG DB: %v", err)
	}
	defer postgresStore.Close()

	store, err := storage.New(postgresStore)
	if err != nil {
		log.Fatalf("[app] failed to init unified storage: %v", err)
	}
	log.Println("[app] Unified storage initialized successfully")

	serviceURL := service.New(store)

	engine := ginext.New("release")

}
