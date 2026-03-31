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
	competitionHandler := handlers.NewCompetitionHandler()
	groupHandler := handlers.NewGroupHandler()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		competitions := api.Group("/competitions")
		{
			competitions.GET("", competitionHandler.GetCompetitions)
			competitions.GET("/:id", competitionHandler.GetCompetition)
		}

		matches := api.Group("/matches")
		{
			matches.GET("", matchHandler.GetMatches)
			matches.GET("/:id", matchHandler.GetMatch)

			matchesAuth := matches.Use(middleware.AuthMiddleware(cfg))
			{
				matchesAuth.POST("", matchHandler.CreateMatch)
				matchesAuth.PUT("/:id/result", matchHandler.UpdateMatchResult)
				matchesAuth.GET("/:id/predictions", predictionHandler.GetMatchPredictions)
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

		// Group routes
		groups := api.Group("/groups")
		groups.Use(middleware.AuthMiddleware(cfg))
		{
			groups.POST("", groupHandler.CreateGroup)
			groups.GET("", groupHandler.GetMyGroups)
			groups.POST("/join", groupHandler.JoinGroup)

			groups.GET("/:id", groupHandler.GetGroup)
			groups.DELETE("/:id", groupHandler.DeleteGroup)
			groups.DELETE("/:id/leave", groupHandler.LeaveGroup)

			groups.GET("/:id/members", groupHandler.GetMembers)

			groups.GET("/:id/competitions", groupHandler.GetCompetitions)
			groups.POST("/:id/competitions", groupHandler.AddCompetition)
			groups.DELETE("/:id/competitions/:competitionId", groupHandler.RemoveCompetition)

			groups.GET("/:id/leaderboard/:competitionId", groupHandler.GetLeaderboard)
			groups.PUT("/:id/transfer-owner", groupHandler.TransferOwnership)
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
