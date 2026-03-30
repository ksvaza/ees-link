package httpapi

import (
	"github.com/julienschmidt/httprouter"
)

func setupApiEndpoints(router *httprouter.Router) {

	// Test endpoints
	router.GET("/api/test", Handler(test))

	// Auth
	router.POST("/api/login", PointLogin)

	// Users
	router.GET("/api/users", PointGetUsers)

	// Competitors
	router.GET("/api/cars", PointGetCars)
	router.POST("/api/cars", PointPostCars)

	// Races
	router.POST("/api/race/start", PointRaceStart)
	router.POST("/api/car/finish", PointCarFinish)
	router.GET("/api/races", PointGetRaces)
	router.POST("/api/races", PointPostRaces)
	router.GET("/api/race-config", PointGetRaceConfig)
	router.GET("/api/categories", PointGetCategories)

	// Results
	router.GET("/api/results/:raceName", PointGetRaceResults)

	// Points
	router.POST("/api/points", PointPostPoints)
	router.DELETE("/api/points", PointDeletePoints)

	// Leaderboard
	router.GET("/api/leaderboard/:ageGroup", PointGetLeaderboard)

	// Admin

	// News
	router.GET("/api/news", PointGetNews)
	router.POST("/api/news", PointPostNews)
	router.GET("/api/news/:id", PointGetNewsByID)
	router.PUT("/api/news/:id", PointPutNewsByID)
	router.DELETE("/api/news/:id", PointDeleteNewsByID)
	router.POST("/api/news/:id/attachments", PointPostNewsAttachments)

	// Events
	router.GET("/api/events", PointGetEvents)
	router.POST("/api/events", PointPostEvents)
	router.GET("/api/events/:id", PointGetEventByID)
	router.PUT("/api/events/:id", PointPutEventByID)
	router.DELETE("/api/events/:id", PointDeleteEventByID)

	// ir -- Pieteikumi
	router.GET("/api/applications", Handler(PointGetApplications))
	router.POST("/api/applications", Handler(PointPostApplications))
	router.PATCH("/api/applications/:id", BasicAuth(Handler(PointPatchApplicationByID)))

	// Live websocket token endpoint (optional handler if needed)
	router.GET("/ws", PointWebSocket)
}

func SetupHTTPAPI() error {
	router := httprouter.New()

	setupApiEndpoints(router)

	return nil
}
