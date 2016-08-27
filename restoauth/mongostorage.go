package restoauth

import (
	"log"
	"os"

	"github.com/RangelReale/osin"
	"github.com/mitchellh/mapstructure"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// collection names for the entities
const (
	CLIENTCOL    = "clients"
	AUTHORIZECOL = "authorizations"
	ACCESSCOL    = "accesses"
)

// REFRESHTOKEN  names for the entities
const REFRESHTOKEN = "refreshtoken"

//MongoStorage  config
type MongoStorage struct {
	dbName  string
	session *mgo.Session
}

//GetenvOrDefault key,def
func GetenvOrDefault(key, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	return value
}

//SetMgoClient1234 storage
func SetMgoClient1234(storage *MongoStorage) (osin.Client, error) {
	client := &osin.DefaultClient{
		Id:          "1234",
		Secret:      "aabbccdd",
		RedirectUri: "http://localhost:14000/appauth"}
	err := storage.SetClient("1234", client)
	return client, err
}

//Close storage
func (store *MongoStorage) Close() {
}

//Clone storage
func (store *MongoStorage) Clone() osin.Storage {
	return store
}

//NewMgoStorage the mongo session
func NewMgoStorage(session *mgo.Session, dbName string) *MongoStorage {
	storage := &MongoStorage{dbName, session}
	index := mgo.Index{
		Key:        []string{REFRESHTOKEN},
		Unique:     false, // refreshtoken is sometimes empty
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}
	accesses := storage.session.DB(dbName).C(ACCESSCOL)
	err := accesses.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	return storage
}

//GetClient ID
func (store *MongoStorage) GetClient(id string) (osin.Client, error) {
	session := store.session.Copy()
	defer session.Close()
	clients := session.DB(store.dbName).C(CLIENTCOL)
	client := new(osin.DefaultClient)
	err := clients.FindId(id).One(&client)
	if err != nil {
		log.Printf("GetClient: %s\n", err)
	}
	return client, err
}

//SetClient ID client
func (store *MongoStorage) SetClient(id string, client osin.Client) error {
	session := store.session.Copy()
	defer session.Close()
	clients := session.DB(store.dbName).C(CLIENTCOL)
	_, err := clients.UpsertId(id, client)

	if err != nil {
		log.Printf("SetClient: %s\n", err)
	}
	return err
}

//SaveAuthorize  auth data
func (store *MongoStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	session := store.session.Copy()
	defer session.Close()
	authorizations := session.DB(store.dbName).C(AUTHORIZECOL)
	_, err := authorizations.UpsertId(data.Code, data)
	if err != nil {
		log.Printf("SaveAuthorize: %s\n", err)
	}
	return err
}

//LoadAuthorize code
func (store *MongoStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	session := store.session.Copy()
	c := session.DB(store.dbName).C(AUTHORIZECOL)
	defer session.Close()
	result := bson.M{}

	err := c.Find(bson.M{"code": code}).One(&result)
	if err != nil {
		return nil, err
	}

	clientID := result["client"].(bson.M)["id"].(string)
	client, err := store.GetClient(clientID)
	if err != nil {
		return nil, err
	}
	result["Client"] = client
	var authData osin.AuthorizeData
	err = mapstructure.Decode(result, &authData)
	if err != nil {
		return nil, err
	}
	return &authData, err
}

//RemoveAuthorize code
func (store *MongoStorage) RemoveAuthorize(code string) error {
	session := store.session.Copy()
	defer session.Close()
	authorizations := session.DB(store.dbName).C(AUTHORIZECOL)
	return authorizations.RemoveId(code)
}

//SaveAccess data
func (store *MongoStorage) SaveAccess(data *osin.AccessData) error {
	session := store.session.Copy()
	defer session.Close()
	accesses := session.DB(store.dbName).C(ACCESSCOL)
	_, err := accesses.UpsertId(data.AccessToken, data)

	if err != nil {
		log.Printf("SaveAccess:%s\n", err)
	}
	return err
}

//LoadAccess token
func (store *MongoStorage) LoadAccess(token string) (*osin.AccessData, error) {
	session := store.session.Copy()
	c := session.DB(store.dbName).C(ACCESSCOL)
	defer session.Close()
	accData := new(osin.AccessData)
	result := bson.M{}
	err := c.FindId(token).One(&result)
	if err != nil {
		return nil, err
	}

	clientID := result["client"].(bson.M)["id"].(string)
	client, err := store.GetClient(clientID)
	if err != nil {
		return nil, err
	}
	result["Client"] = client
	result["AuthorizeData"] = nil
	err = mapstructure.Decode(result, &accData)
	if err != nil {
		return nil, err
	}

	if err != nil {
		log.Printf("LoadAccess:%s\n", err)
	}
	return accData, err
}

//RemoveAccess token
func (store *MongoStorage) RemoveAccess(token string) error {
	session := store.session.Copy()
	defer session.Close()
	accesses := session.DB(store.dbName).C(ACCESSCOL)
	return accesses.RemoveId(token)
}

//LoadRefresh token
func (store *MongoStorage) LoadRefresh(token string) (*osin.AccessData, error) {
	session := store.session.Copy()
	defer session.Close()
	accesses := session.DB(store.dbName).C(ACCESSCOL)
	accData := new(osin.AccessData)
	result := bson.M{}
	err := accesses.Find(bson.M{REFRESHTOKEN: token}).One(result)
	if err != nil {
		return nil, err
	}

	clientID := result["client"].(bson.M)["id"].(string)
	client, err := store.GetClient(clientID)
	if err != nil {
		return nil, err
	}
	result["Client"] = client
	result["AuthorizeData"] = nil
	err = mapstructure.Decode(result, &accData)
	if err != nil {
		return nil, err
	}

	if err != nil {
		log.Printf("LoadAccess:%s\n", err)
	}
	return accData, err
}

//RemoveRefresh token
func (store *MongoStorage) RemoveRefresh(token string) error {
	session := store.session.Copy()
	defer session.Close()

	accesses := session.DB(store.dbName).C(ACCESSCOL)
	return accesses.Update(bson.M{REFRESHTOKEN: token}, bson.M{
		"$unset": bson.M{
			REFRESHTOKEN: 1,
		}})
}
