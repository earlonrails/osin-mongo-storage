package main

import (
    "fmt"
    "net/http"
    "os"

    "github.com/earlonrails/osin-mongo-storage/mgostore"

    "github.com/RangelReale/osin"
    "github.com/ant0ine/go-json-rest/rest"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "gopkg.in/mgo.v2"
)

func main() {
    mainRouter := mux.NewRouter()
    oAuth := setupOAuth(mainRouter)
    setupRestAPI(mainRouter, oAuth)

    port := fmt.Sprintf(":%v", getenvOrDefault("PORT", "3000"))
    fmt.Printf("Listening on port %v\n", port)

    http.ListenAndServe(port, mainRouter)
}

func setupOAuth(router *mux.Router) *oAuthHandler {
    session, err := mgo.Dial(getenvOrDefault("MGOSTORE_MONGO_URL", "localhost"))
    if err != nil {
        panic(err)
    }
    oAuth := NewOAuthHandler(session, "osinmongostoragetest")

    if _, err := oAuth.Storage.GetClient("1234"); err != nil {
        if _, err := setClient1234(oAuth.Storage); err != nil {
            panic(err)
        }
    }

    oauthSub := router.PathPrefix("/oauth2").Subrouter()
    oauthSub.HandleFunc("/authorize", oAuth.AuthorizeClient)
    oauthSub.HandleFunc("/token", oAuth.GenerateToken)
    oauthSub.HandleFunc("/info", oAuth.HandleInfo)

    return oAuth
}

func setupRestAPI(router *mux.Router, oAuth *oAuthHandler) {
    api := rest.NewApi()
    api.Use(rest.DefaultDevStack...)
    restRouter, err := rest.MakeRouter(
        rest.Get("/api/me", func(w rest.ResponseWriter, req *rest.Request) {
            data := context.Get(req.Request, USERDATA)
            w.WriteJson(&data)
        }),
    )
    if err != nil {
        panic(err)
    }
    api.SetApp(restRouter)
    router.Handle("/api/me", api.MakeHandler())
}

func setClient1234(storage *mgostore.MongoStorage) (*osin.DefaultClient, error) {
    client := &osin.DefaultClient{
        Id:          "1234",
        Secret:      "aabbccdd",
        RedirectUri: "http://localhost:14000/appauth"}
    err := storage.SetClient("1234", client)
    return client, err
}

func getenvOrDefault(key, def string) string {
    value := os.Getenv(key)
    if value == "" {
        return def
    }
    return value
}