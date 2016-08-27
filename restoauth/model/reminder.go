package model

import (
	"time"

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

//GetAllReminders w,r
func (i *DBImpl) GetAllReminders(w rest.ResponseWriter, r *rest.Request) {
	reminders := []Reminder{}
	iter := i.remind.Find(nil).Limit(100).Iter()
	if iter.All(&reminders) != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&reminders)
}

//GetReminder w,r
func (i *DBImpl) GetReminder(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	reminder := Reminder{}
	if i.remind.FindId(id).One(&reminder) != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&reminder)
}

//PostReminder w,r
func (i *DBImpl) PostReminder(w rest.ResponseWriter, r *rest.Request) {

}

//PutReminder w,r
func (i *DBImpl) PutReminder(w rest.ResponseWriter, r *rest.Request) {

}

//DeleteReminder w,r
func (i *DBImpl) DeleteReminder(w rest.ResponseWriter, r *rest.Request) {
}
