package httpapi

import (
	"mime"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func setupApiEndpoints(router *httprouter.Router) {

	//tīri testi
	router.GET("/api/test", PointTestAPI)
	router.POST("/api/submit-form", PointReceiveForm)
}

func setupHTTPHost(router *httprouter.Router) error {
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")

	fileServer := http.FileServer(http.Dir("public"))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ws" || len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api" {
			router.ServeHTTP(w, r)
		} else {
			fileServer.ServeHTTP(w, r)
		}
	})

	err := http.ListenAndServe(":1884", handler)
	return err
}

func SetupHTTPAPI() error {
	router := httprouter.New()

	setupApiEndpoints(router)

	err := setupHTTPHost(router)

	return err
}
