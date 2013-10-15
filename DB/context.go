package DB

import (
	"labix.org/v2/mgo"
)

type Context struct {
	Database *mgo.Database
}

var Session *mgo.Session

func (c *Context) Close() {
	c.Database.Session.Close()
}

func NewContext() (*Context, error) {
	return &Context{
		Database: Session.Clone().DB(Config["DBName"]),
	}, nil
}
