package model

import (
	"os"

	mgo "gopkg.in/mgo.v2"
)

//DBInterface inter
type DBInterface interface {
	//setModel(table string) *mgo.Collection
	InitDB(dbName string)
}

//DBImpl with mgo
type DBImpl struct {
	DB     *mgo.Session
	dbName string
	remind *mgo.Collection
}

// some const
const (
	REMIND string = "remind"
)

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
	i.DB.SetMode(mgo.Monotonic, true)
	i.remind = i.DB.DB(dbName).C(REMIND)
}
