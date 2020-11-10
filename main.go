package main

import (
	"encoding/csv"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "zt45"
	password = "zt45duke"
	dbname   = "web_dev"
)

type exercise struct {
	ID         int
	Met        float64
	Category1  string
	Category2  string
	Activities string
}

func migrate() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,          // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		Logger: newLogger,
	})
	hasTable := db.Migrator().HasTable(&exercise{})
	if hasTable {
		fmt.Printf("Has Table %T\n", exercise{})
	} else {
		fmt.Printf("Does not have Table %T\n", exercise{})
	}
	db.Migrator().DropTable(&exercise{})
	hasTable = db.Migrator().HasTable(&exercise{})
	if hasTable {
		fmt.Printf("Has Table %T\n", exercise{})
	} else {
		fmt.Printf("Table %T Removed\n", exercise{})
	}
	db.AutoMigrate(&exercise{})

	if err != nil {
		println("connect error")
		log.Fatal(err)
	}
	fmt.Println("Database is connected")
	fmt.Println("Getting file")
	csvFile, err := os.Open("PAC2011.csv")
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}

	reader := csv.NewReader(csvFile)
	allRecords, err := reader.ReadAll()
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}
	for index, line := range allRecords {
		if index != 0 {
			i, _ := strconv.Atoi(line[0])
			met, _ := strconv.ParseFloat(line[1], 64)
			temp := exercise{
				ID:         i,
				Met:        met,
				Category1:  string(line[2]),
				Category2:  string(line[3]),
				Activities: string(line[4]),
			}
			result := db.Create(&temp)
			if result.Error != nil {
				fmt.Println("An error encountered ::", err)
				return
			}
		}

	}

	err = csvFile.Close()
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return
	}

}
func main() {
	migrate()
	fmt.Println("Finish Migrating!")
}
