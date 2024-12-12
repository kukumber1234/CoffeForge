package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"hot-coffee/internal/dal"
	"hot-coffee/internal/service"

	"hot-coffee/internal/handler"

	_ "github.com/lib/pq"

	flags "hot-coffee/models"
)

const (
	host     = "PSQL Local"
	port     = 5432
	user     = "latte"
	password = "latte"
	dbname   = "frappuccino"
)

func main() {
	flag.Parse()

	if *flags.Help {
		dal.Helper()
		return
	}

	dal.Create()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("latte", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	// menu
	menuDal := dal.NewMenuRepo(*flags.Dir + "/menu.json")
	menuService := service.NewFileMenuService(menuDal)
	menuHandler := handler.NewMenuHandler(menuService)

	mux.HandleFunc("POST /menu", menuHandler.Add)
	mux.HandleFunc("GET /menu", menuHandler.Get)
	mux.HandleFunc("GET /menu/{id}", menuHandler.GetByID)
	mux.HandleFunc("PUT /menu/{id}", menuHandler.Update)
	mux.HandleFunc("DELETE /menu/{id}", menuHandler.Delete)

	inventoryDal := dal.NewInventoryRepo(*flags.Dir + "/inventory.json")
	inventoryService := service.NewInventoryService(inventoryDal)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)

	ordersService := service.NewFileOrderService(*flags.Dir+"/orders.json", menuService, inventoryService)
	ordersHandler := handler.NewOrdersHandler(ordersService)

	// orders
	mux.HandleFunc("POST /orders", ordersHandler.Add)
	mux.HandleFunc("GET /orders", ordersHandler.Get)
	mux.HandleFunc("GET /orders/{id}", ordersHandler.GetByID)
	mux.HandleFunc("PUT /orders/{id}", ordersHandler.Update)
	mux.HandleFunc("DELETE /orders/{id}", ordersHandler.Delete)
	mux.HandleFunc("POST /orders/{id}/close", ordersHandler.Close)

	// inventory
	mux.HandleFunc("POST /inventory", inventoryHandler.Add)
	mux.HandleFunc("GET /inventory", inventoryHandler.Get)
	mux.HandleFunc("GET /inventory/{id}", inventoryHandler.GetByID)
	mux.HandleFunc("PUT /inventory/{id}", inventoryHandler.Update)
	mux.HandleFunc("DELETE /inventory/{id}", inventoryHandler.Delete)

	reportsDal := dal.NewReportsRepo(*flags.Dir + "/orders.json")
	reportsService := service.NewFileReportsService(reportsDal)
	reportsHandler := handler.NewReportsHandler(reportsService)

	// reports
	mux.HandleFunc("GET /reports/total-sales", reportsHandler.GetTotalSales)
	mux.HandleFunc("GET /reports/popular-items", reportsHandler.GetPopularItems)

	flags.Logger.Info("Host is started")
	log.Fatal(http.ListenAndServe(":"+*flags.Port, mux))
}
