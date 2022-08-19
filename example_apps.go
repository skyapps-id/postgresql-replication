package main

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type Order struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
}

func main() {
	dbMasterURL := "postgres://postgres:postgres@localhost:5000/order"
	dbSlaveURL := "postgres://postgres:postgres@localhost:5001/order"

	db, err := gorm.Open(postgres.Open(dbMasterURL), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.Use(
		dbresolver.Register(dbresolver.Config{
			// Sources:  []gorm.Dialector{postgres.Open(dbMasterURL)},
			Replicas: []gorm.Dialector{postgres.Open(dbSlaveURL)},
		}).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(100).
			SetMaxOpenConns(200))

	db.AutoMigrate(&Order{})

	var (
		order Order
		ipAdd string
	)
	order = Order{Name: "Test"}

	fmt.Println("Create Data :")
	db.Create(&order)
	db.Clauses(dbresolver.Write).Raw("SELECT inet_server_addr()").First(&ipAdd)
	fmt.Println(order)
	fmt.Println(fmt.Sprintf("IP DB Master: %s", ipAdd))

	fmt.Println()

	fmt.Println("Read Data :")
	db.Take(&order)
	db.Clauses(dbresolver.Read).Raw("SELECT inet_server_addr()").First(&ipAdd)
	fmt.Println(order)
	fmt.Println(fmt.Sprintf("IP DB Slave: %s", ipAdd))

	/* tx := db.Clauses(dbresolver.Write).Begin()
	tx.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
	}
	if err := tx.Raw("SELECT inet_server_addr()").First(&ipAdd).Error; err != nil {
		tx.Rollback()
	}
	tx.Commit() */

}
