package main

import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

func main() {

	api := rest.NewApi()
	oauthHand := NewOAuthHandler("session", "dbname")
	// the Middleware stack
	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			return request.Method == "POST"
		},
		IfTrue: &FormMiddleware{},
	})
	api.Use([]rest.Middleware{
		&rest.ContentTypeCheckerMiddleware{},
	}...)

	api.Use(rest.DefaultDevStack...)
	api.Use(oauthHand)

	// build the App, here the rest Router
	router, err := rest.MakeRouter(
		rest.Get("/api/v1/message", func(w rest.ResponseWriter, req *rest.Request) {
			w.WriteJson(map[string]string{"Body": "Hello World!"})
		}),
		rest.Get("/oauth/authorize", func(w rest.ResponseWriter, req *rest.Request) {
			oauthHand.AuthorizeClient(w.(http.ResponseWriter), req.Request)
			w.WriteJson(map[string]string{"msg": "ok!", "code": "200"})
		}),
		rest.Post("/oauth/token", func(w rest.ResponseWriter, req *rest.Request) {
			oauthHand.GenerateToken(w.(http.ResponseWriter), req.Request)
			w.WriteJson(map[string]string{"msg": "ok!", "code": "200"})
		}),
		rest.Get("/oauth/info", func(w rest.ResponseWriter, req *rest.Request) {
			oauthHand.HandleInfo(w.(http.ResponseWriter), req.Request)
			w.WriteJson(map[string]string{"msg": "ok!", "code": "200"})
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	// build and run the handler
	log.Fatal(http.ListenAndServe(
		":8080",
		api.MakeHandler(),
	))
}
