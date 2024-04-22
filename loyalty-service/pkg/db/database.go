package db

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func Connect(uris []string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(uris[0]), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to MySQL database: %v", err)
		return nil, err
	}

	var sources []gorm.Dialector

	if len(uris) > 1 {
		for _, uri := range uris[1:] {
			sources = append(sources, mysql.Open(uri))
		}

		err = db.Use(dbresolver.Register(dbresolver.Config{
			Sources: sources,
		}))

		if err != nil {
			return nil, err
		}
	}

	log.Println("Connected to MySQL")
	return db, nil
}
