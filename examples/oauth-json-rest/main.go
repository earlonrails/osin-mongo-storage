package main

import (
	"log"
	"net/http"
	"restoauth"

	mgo "gopkg.in/mgo.v2"

	"github.com/ant0ine/go-json-rest/rest"
)

func main() {

	api := rest.NewApi()
	//oauthHand := restoauth.NewOAuthHandler("session", "dbname")

	//mongo oauth middleware hand
	session, err := mgo.Dial(restoauth.GetenvOrDefault("MGOSTORE_MONGO_URL", "localhost"))
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	if err != nil {
		panic(err)
	}
	oauthHand := restoauth.NewOAuthHandlerByMgo(session, "osinmongostorage")
	mgostore := oauthHand.Storage.(*restoauth.MongoStorage)
	restoauth.SetMgoClient1234(mgostore)

	// the Middleware stack
	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			return request.Method == "POST"
		},
		IfTrue: &restoauth.FormMiddleware{},
	})
	api.Use([]rest.Middleware{
		&rest.ContentTypeCheckerMiddleware{},
	}...)

	api.Use(rest.DefaultDevStack...)
	api.Use(oauthHand)

	// build the App, here the rest Router
	router, err := rest.MakeRouter(
		rest.Get("/api/v1/message", func(w rest.ResponseWriter, req *rest.Request) {
			restoauth.OutJSON(w, "ok", 200, 200)
		}),
		rest.Get("/oauth/authorize", func(w rest.ResponseWriter, req *rest.Request) {
			oauthHand.AuthorizeClient(w.(http.ResponseWriter), req.Request)
			//restoauth.OutJSON(w, "ok", 200, 200)
		}),
		rest.Post("/oauth/token", func(w rest.ResponseWriter, req *rest.Request) {
			oauthHand.GenerateToken(w.(http.ResponseWriter), req.Request)
			//restoauth.OutJSON(w, "ok", 200, 200)
		}),
		rest.Get("/oauth/info", func(w rest.ResponseWriter, req *rest.Request) {
			oauthHand.HandleInfo(w.(http.ResponseWriter), req.Request)
			//restoauth.OutJSON(w, "ok", 200, 200)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	// build and run the handler
	log.Fatal(http.ListenAndServe(
		":3000",
		api.MakeHandler(),
	))
}
