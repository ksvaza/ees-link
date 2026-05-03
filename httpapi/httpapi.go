package httpapi

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func setupApiEndpoints(router *httprouter.Router) {

	// Test endpoints
	router.GET("/api/test", Handler(test))
	router.GET("/api/testdb", Handler(TestDatabase))
	router.POST("/api/submit-form", TestPointReceiveForm)

	// Auth and account management
	router.POST("/api/register", Handler(PointRegister))
	router.GET("/api/login", Handler(PointLogin))
	router.GET("/api/account-applications", Handler(PointGetAccountApplications))
	router.GET("/api/account-verification-info/:key", Handler(PointGetAccountVerificationInfo))
	router.POST("/api/verify-account-application/:key", Handler(PointVerifyAccountApplication))
	router.GET("/api/accounts", Handler(PointGetAccounts))
	router.PATCH("/api/account/:key", Handler(PointPatchAccountByKey))
	router.GET("/api/admin-accounts", Handler(PointGetAdminAccounts))
	router.PATCH("/api/admin-account/:key", Handler(PointPatchAdminAccount))

	// Competitors
	// superadmins akceptē komandas , atsevisks strukts komandām
	// router.POST("/api/cars", PointPostCars)

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
	router.GET("/api/teams", Handler(PointGetTeams))
	router.POST("/api/team-application", Handler(PointPostTeamApplication))
	router.PATCH("/api/team/:key", Handler(PointPatchTeamDataByKey))
	router.POST("/api/team/:key", Handler(PointPostTeamDataByKey))
	// Vēl vajag PATCH /api/account-applications manuālās verifikācijas ar pieteikuma pamainīšanu, kur visadministrators var verificēt visu, bet komandas līderis var tikai verificēt savas komandas pieteikumus.

	// Live websocket token endpoint (optional handler if needed)
	router.GET("/ws", PointWebSocket)
}

// Neizdzēšanas zona
// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func redirectHTTPToHTTPS(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Forwarded-Proto") != "https" && r.Header.Get("X-Forwarded-Proto") != "" {
		http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
		return
	}
	if r.TLS == nil && r.Header.Get("X-Forwarded-Proto") == "" {
		http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
		return
	}
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		f.Close()
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			f.Close()
			return nil, os.ErrNotExist // returns 404 instead of listing
		}
	}
	return f, nil
}

func setupHTTPHost(router *httprouter.Router) error {
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("public")})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ws" || strings.HasPrefix(r.URL.Path, "/api/") {
			router.ServeHTTP(w, r)
			return
		}

		//redirectHTTPToHTTPS(w, r)

		path := filepath.Join("public", filepath.Clean(r.URL.Path))
		if strings.HasPrefix(path, "../") || strings.Contains(path, "/../") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		info, err := os.Stat(path)
		if os.IsNotExist(err) || (err == nil && info.IsDir()) {
			if _, err := os.Stat("public/index.html"); err == nil {
				http.ServeFile(w, r, "public/index.html")
			}
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	return http.ListenAndServe(":1884", handler)
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func SetupHTTPAPI() error {
	router := httprouter.New()

	// Nedrīks izdzēst
	// ---------------
	setupApiEndpoints(router)

	err := setupHTTPHost(router)
	// ---------------

	return err
}
