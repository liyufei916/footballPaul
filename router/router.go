package router

import (
	"github.com/gin-gonic/gin"
	"github.com/liyufei916/footballPaul/config"
	"github.com/liyufei916/footballPaul/handlers"
	"github.com/liyufei916/footballPaul/middleware"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	userHandler := handlers.NewUserHandler(cfg)
	matchHandler := handlers.NewMatchHandler()
	predictionHandler := handlers.NewPredictionHandler()
	leaderboardHandler := handlers.NewLeaderboardHandler()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		matches := api.Group("/matches")
		{
			matches.GET("", matchHandler.GetMatches)
			matches.GET("/:id", matchHandler.GetMatch)
			matches.GET("/:matchId/predictions", predictionHandler.GetMatchPredictions)

			matchesAuth := matches.Use(middleware.AuthMiddleware(cfg))
			{
				matchesAuth.POST("", matchHandler.CreateMatch)
				matchesAuth.PUT("/:id/result", matchHandler.UpdateMatchResult)
			}
		}

		predictions := api.Group("/predictions")
		predictions.Use(middleware.AuthMiddleware(cfg))
		{
			predictions.POST("", predictionHandler.CreatePrediction)
			predictions.PUT("/:id", predictionHandler.UpdatePrediction)
			predictions.GET("/my", predictionHandler.GetUserPredictions)
		}

		leaderboard := api.Group("/leaderboard")
		{
			leaderboard.GET("", leaderboardHandler.GetLeaderboard)

			leaderboardAuth := leaderboard.Use(middleware.AuthMiddleware(cfg))
			{
				leaderboardAuth.GET("/my-rank", leaderboardHandler.GetUserRank)
			}
		}

		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(cfg))
		{
			users.GET("/profile", userHandler.GetProfile)
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return r
}
