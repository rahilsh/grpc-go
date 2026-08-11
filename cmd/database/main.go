// Command database groups the database utilities: a GORM example, a raw
// database/sql example, and a TiDB PD region splitter. Pick one with -demo:
//
//	go run ./cmd/database -demo orm
//
// These require running MySQL/TiDB services and are provided as references.
package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	sqlmysql "github.com/go-sql-driver/mysql"
	"github.com/rahilsh/golang-lab/internal/demo"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() { demo.Run(demos) }

var demos = map[string]func(){
	"orm":             orm,
	"dao":             dao,
	"region-splitter": regionSplitter,
}

// Product is a GORM-managed model.
type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func orm() {
	db, err := gorm.Open(gormmysql.Open("root:root@tcp(127.0.0.1:3306)/test"), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect database: %s", err)
		return
	}
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxOpenConns(20)
	}
	if err := db.AutoMigrate(&Product{}); err != nil {
		log.Printf("migrate: %s", err)
		return
	}

	db.Create(&Product{Code: "D42", Price: 100})

	var product Product
	db.First(&product, 1)
	fmt.Printf("id: %d\n", product.ID)
	db.Model(&product).Update("Price", 200)
	fmt.Printf("price: %d\n", product.Price)
	db.Delete(&product, 1)
}

func dao() {
	cfg := sqlmysql.Config{
		User:                 "root",
		Passwd:               "root",
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               "braas",
		AllowNativePasswords: true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Print(err)
		return
	}
	defer func() { _ = db.Close() }()
	if err := db.Ping(); err != nil {
		log.Print(err)
		return
	}
	fmt.Println("Connected!")

	rows, err := db.Query("select id from rtos")
	if err != nil {
		log.Print(err)
		return
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Print(err)
			return
		}
		log.Print(id)
	}
}

const pdHost = "http://localhost:4004"
const maxNoOfRegionsToFetch = 20
const maxAllowedRegionSizeInMB = 20000

// Region is a TiKV region as reported by the PD API.
type Region struct {
	ID              int `json:"id"`
	ApproximateSize int `json:"approximate_size"`
}

// GetRegionsResponse is the PD regions response.
type GetRegionsResponse struct {
	Count   int      `json:"count"`
	Regions []Region `json:"regions"`
}

// SplitRegionRequest is a PD split-region operator request.
type SplitRegionRequest struct {
	Name     string `json:"name"`
	Policy   string `json:"policy"`
	RegionID int    `json:"region_id"`
}

func regionSplitter() {
	for _, region := range largestRegionsDescending() {
		log.Printf("Region: %d, size: %d", region.ID, region.ApproximateSize)
		if region.ApproximateSize < maxAllowedRegionSizeInMB {
			log.Printf("No more region greater than %d", maxAllowedRegionSizeInMB)
			break
		}
		splitRegion(region.ID)
	}
}

func largestRegionsDescending() []Region {
	res, err := http.Get(fmt.Sprintf("%s/pd/api/v1/regions/size?limit=%d", pdHost, maxNoOfRegionsToFetch))
	if err != nil {
		log.Print(err)
		return nil
	}
	defer func() { _ = res.Body.Close() }()

	var resp GetRegionsResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		log.Printf("decode: %s", err)
		return nil
	}
	return resp.Regions
}

func splitRegion(id int) {
	log.Printf("Splitting region %d", id)
	jsonData, _ := json.Marshal(SplitRegionRequest{"split-region", "approximate", id})
	response, err := http.Post(fmt.Sprintf("%s/pd/api/v1/operators", pdHost), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Print(err)
		return
	}
	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("%s", body)
}
