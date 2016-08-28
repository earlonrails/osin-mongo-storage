package restoauth

import (
	"os"

	"github.com/ant0ine/go-json-rest/rest"

	mgo "gopkg.in/mgo.v2"
)

//TableInterface  inter
type TableInterface interface {
	GetOneRes(w rest.ResponseWriter, r *rest.Request)
	GetAllRes(w rest.ResponseWriter, r *rest.Request)
	RemoveRes(w rest.ResponseWriter, r *rest.Request)
	CreateRes(w rest.ResponseWriter, r *rest.Request)
	UpdateRes(w rest.ResponseWriter, r *rest.Request)
}

//DBImpl with mgo
type DBImpl struct {
	DB     *mgo.Session
	DBName string
}

//GetenvOrDefault key,def
func GetenvOrDefault(key, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	return value
}

//InitDB with mgo
func (i *DBImpl) InitDB(dbName string) {
	var err error
	i.DB, err = mgo.Dial(GetenvOrDefault("MGOSTORE_MONGO_URL", "localhost"))
	if err != nil {
		panic(err)
	}
	i.DBName = dbName
	i.DB.SetMode(mgo.Monotonic, true)
}
