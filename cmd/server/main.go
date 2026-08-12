package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
	"github.com/nikorokni/Go-PostgreSQL-Redis-WebSocket-Docker-Prometheus-Grafana/internal/api"
)

func main(){ addr:=os.Getenv("HTTP_ADDR");if addr==""{addr=":8080"};srv:=&http.Server{Addr:addr,Handler:api.New().Routes(),ReadHeaderTimeout:5*time.Second,ReadTimeout:10*time.Second,WriteTimeout:15*time.Second,IdleTimeout:60*time.Second};slog.Info("stablecoin risk engine started","address",addr);if err:=srv.ListenAndServe();err!=nil&&err!=http.ErrServerClosed{slog.Error("server failed","error",err);os.Exit(1)} }
