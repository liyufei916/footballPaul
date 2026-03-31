package main

import (
	"log"

	"github.com/liyufei916/footballPaul/config"
	"github.com/liyufei916/footballPaul/database"
	"github.com/liyufei916/footballPaul/router"
)

func main() {
	cfg := config.LoadConfig()

	if err := database.InitDatabase(cfg); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	if err := database.SeedDefaultScoringRules(); err != nil {
		log.Println("Warning: Failed to seed scoring rules:", err)
	}

	if err := database.SeedDefaultCompetitions(); err != nil {
		log.Println("Warning: Failed to seed competitions:", err)
	}

	r := router.SetupRouter(cfg)

	log.Printf("Server starting on port %s...", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
