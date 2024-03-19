package database

import (
	"fmt"
	"log"
	"sample-manager/models"

	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	username = "abhyudaya.a_ftc"
	password = "password"
	dbName   = "samplemanager"
	sslMode  = "disable"
)

func Connection() *gorm.DB {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, username, password, dbName, sslMode)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}

	log.Println("Connected to the database")

	err = db.AutoMigrate(&models.Mapping{})
	if err != nil {
		log.Fatalf("Error migrating Mapping table: %v", err)
	}

	return db
}

func GetSampleID(db *gorm.DB, segments []string, itemID string) (string, error) {
	var mapping models.Mapping

	result := db.Where("segments = ?::text[] AND item_id = ?", pq.StringArray(segments), itemID).First(&mapping)

	if result.Error != nil {
		return "", result.Error
	}
	return mapping.SampleItemID, nil
}
