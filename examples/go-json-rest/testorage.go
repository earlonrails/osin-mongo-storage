package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/ant0ine/go-json-rest/rest"
)

//TestStorage struct
type TestStorage struct {
	clients   map[string]osin.Client
	authorize map[string]*osin.AuthorizeData
	access    map[string]*osin.AccessData
	refresh   map[string]string
}

//NewTestStorage new
func NewTestStorage() *TestStorage {
	r := &TestStorage{
		clients:   make(map[string]osin.Client),
		authorize: make(map[string]*osin.AuthorizeData),
		access:    make(map[string]*osin.AccessData),
		refresh:   make(map[string]string),
	}

	r.clients["1234"] = &osin.DefaultClient{
		Id:          "1234",
		Secret:      "aabbccdd",
		RedirectUri: "http://localhost:14000/appauth",
	}

	return r
}

//Clone storage
func (s *TestStorage) Clone() osin.Storage {
	return s
}

//Close storage
func (s *TestStorage) Close() {
}

//GetClient by id and return (client,err)
func (s *TestStorage) GetClient(id string) (osin.Client, error) {
	log.Printf("GetClient: %s\n", id)
	if c, ok := s.clients[id]; ok {
		return c, nil
	}
	return nil, errors.New("Client not found")
}

//SetClient by id  and client
func (s *TestStorage) SetClient(id string, client osin.Client) error {
	log.Printf("SetClient: %s\n", id)
	s.clients[id] = client
	return nil
}

//SaveAuthorize data
func (s *TestStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	log.Printf("SaveAuthorize: %s\n", data.Code)
	s.authorize[data.Code] = data
	return nil
}

//LoadAuthorize code and return ()
func (s *TestStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	log.Printf("LoadAuthorize: %s\n", code)
	if d, ok := s.authorize[code]; ok {
		return d, nil
	}
	return nil, errors.New("Authorize not found")
}

//RemoveAuthorize code
func (s *TestStorage) RemoveAuthorize(code string) error {
	log.Printf("RemoveAuthorize: %s\n", code)
	delete(s.authorize, code)
	return nil
}

//SaveAccess data
func (s *TestStorage) SaveAccess(data *osin.AccessData) error {
	log.Printf("SaveAccess:%s\n", data.AccessToken)
	s.access[data.AccessToken] = data
	if data.RefreshToken != "" {
		s.refresh[data.RefreshToken] = data.AccessToken
		log.Printf("SaveRefreshAccess:%s\n", data.RefreshToken)
	}
	return nil
}

//LoadAccess code
func (s *TestStorage) LoadAccess(code string) (*osin.AccessData, error) {
	log.Printf("LoadAccess:%s\n", code)
	if d, ok := s.access[code]; ok {
		log.Printf("LoadRefresh:%s\n", d.RefreshToken)
		return d, nil
	}

	log.Printf("no found Access: %s\n", code)
	return nil, errors.New("Access not found")
}

//RemoveAccess code
func (s *TestStorage) RemoveAccess(code string) error {
	log.Printf("RemoveAccess: %s\n", code)
	delete(s.access, code)
	return nil
}

//LoadRefresh code
func (s *TestStorage) LoadRefresh(code string) (*osin.AccessData, error) {
	log.Printf("LoadRefresh: %s\n", code)
	if d, ok := s.refresh[code]; ok {
		return s.LoadAccess(d)
	}
	return nil, errors.New("Refresh not found")
}

//RemoveRefresh code
func (s *TestStorage) RemoveRefresh(code string) error {
	log.Printf("RemoveRefresh: %s\n", code)
	delete(s.refresh, code)
	return nil
}

//HandleLoginPage with ar resp req
func HandleLoginPage(ar *osin.AuthorizeRequest, w rest.ResponseWriter, r *http.Request, debug bool) bool {
	r.ParseForm()
	if r.Method == "POST" && r.Form.Get("login") == "test" && r.Form.Get("password") == "test" {
		return true
	}

	return debug
}
