package models

import (
	"restoauth"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/ant0ine/go-json-rest/rest"
)

//Reminder struct
type Reminder struct {
	ID        int64     `bson:"id"`
	Message   string    `sql:"size:1024" bson:"message"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
	DeletedAt time.Time `bson:"-"`
}

//ReminderModel  struct
type ReminderModel struct {
	table *mgo.Collection
}

// some const
const (
	REMINDCOL string = "remind"
)

//NewReminderModel tableName
func NewReminderModel(db *restoauth.DBImpl) *ReminderModel {
	model := &ReminderModel{db.DB.DB(db.DBName).C(REMINDCOL)}
	return model
}

//GetAllRes w,r
func (i *ReminderModel) GetAllRes(w rest.ResponseWriter, r *rest.Request) {
	reminders := []Reminder{}
	iter := i.table.Find(nil).Limit(100).Iter()
	if iter.All(&reminders) != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&reminders)
}

//GetOneRes w,r
func (i *ReminderModel) GetOneRes(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	reminder := Reminder{}
	if i.table.FindId(id).One(&reminder) != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&reminder)
}

//CreateRes w,r
func (i *ReminderModel) CreateRes(w rest.ResponseWriter, r *rest.Request) {

}

//UpdateRes w,r
func (i *ReminderModel) UpdateRes(w rest.ResponseWriter, r *rest.Request) {

}

//RemoveRes w,r
func (i *ReminderModel) RemoveRes(w rest.ResponseWriter, r *rest.Request) {
}
