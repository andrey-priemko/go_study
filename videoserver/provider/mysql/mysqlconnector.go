package mysql

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

type Connector struct {
	DB *sql.DB
}

const user = "root"
const pass = "1234"
const dbName = "go_dev"

func (c *Connector) Connect() error {
	if c.DB != nil {
		log.Info("already connected")
		return nil
	}

	DB, err := sql.Open("mysql", user + ":" + pass + "@/" + dbName)
	if err != nil {
		log.Error(err)
		return err
	}

	c.DB = DB

	return nil
}

func (c *Connector) Close() error {
	err := c.DB.Close()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	c.DB = nil

	return nil
}