package models

import (
	"flag"
	"log/slog"
)

var (
	Help = flag.Bool("help", false, "Show help messege")
	Dir  = flag.String("dir", "data", "Path to the data directory")
	Port = flag.String("port", "8080", "Port number")

	PriceList = make(map[string]float64)

	Logger *slog.Logger
)
