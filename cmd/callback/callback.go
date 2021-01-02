package main

import (
	"github.com/joho/godotenv"
	"github.com/kaanaktas/openbanking-accountinformation/internal/cache"
	"github.com/kaanaktas/openbanking-accountinformation/internal/config"
	"github.com/kaanaktas/openbanking-accountinformation/internal/store"
	"github.com/kaanaktas/openbanking-accountinformation/pkg/callback"
	cfg "github.com/kaanaktas/openbanking-accountinformation/pkg/configmanager"
	"github.com/kaanaktas/openbanking-accountinformation/pkg/consent"
	"github.com/kaanaktas/openbanking-accountinformation/pkg/token"
	"log"
	"os"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	port := os.Getenv("PORT_CALLBACK")
	if port == "" {
		port = "8081"
	}

	e := config.NewEchoEngine()
	dbx := store.LoadDBConnection()
	chInMemory := cache.LoadInMemory()
	chInRedis := cache.LoadRedis()

	configRepository := cfg.NewRepository(dbx)
	configService := cfg.NewService(configRepository, chInMemory)
	tokenService := token.NewService(configService)
	consentRepository := consent.NewRepositoryRead(dbx)
	consentService := consent.NewServiceRead(consentRepository)
	callbackRepository := callback.NewRepository(dbx)
	callbackService := callback.NewService(callbackRepository, consentService, tokenService, chInRedis)

	callback.RegisterHandler(e, callbackService)

	log.Printf("starting server at :%s", port)

	if err := e.Start(":" + port); err != nil {
		log.Fatalf("error while starting server at :%s, %v", port, err)
	}
}
