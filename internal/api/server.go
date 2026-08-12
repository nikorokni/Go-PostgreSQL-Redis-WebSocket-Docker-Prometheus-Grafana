package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/nikorokni/Go-PostgreSQL-Redis-WebSocket-Docker-Prometheus-Grafana/internal/hub"
	"github.com/nikorokni/Go-PostgreSQL-Redis-WebSocket-Docker-Prometheus-Grafana/internal/risk"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct { mu sync.RWMutex; latest map[string]risk.Assessment; hub *hub.Hub; assessed prometheus.Counter; critical prometheus.Counter }
func New() *Server {
	s:=&Server{latest:map[string]risk.Assessment{},hub:hub.New(), assessed:prometheus.NewCounter(prometheus.CounterOpts{Name:"risk_assessments_total",Help:"Risk assessments processed."}),critical:prometheus.NewCounter(prometheus.CounterOpts{Name:"critical_assessments_total",Help:"Critical assessments observed."})}
	prometheus.MustRegister(s.assessed,s.critical); return s
}
func (s *Server) Routes() http.Handler { m:=http.NewServeMux(); m.HandleFunc("GET /healthz",func(w http.ResponseWriter,r *http.Request){write(w,http.StatusOK,map[string]string{"status":"ok"})}); m.HandleFunc("POST /v1/assess",s.assess); m.HandleFunc("GET /v1/risks",s.risks); m.HandleFunc("GET /v1/stream",s.stream); m.Handle("GET /metrics",promhttp.Handler()); return withCORS(m) }
func (s *Server) assess(w http.ResponseWriter,r *http.Request){ var in risk.Snapshot; if err:=json.NewDecoder(http.MaxBytesReader(w,r.Body,1<<20)).Decode(&in); err!=nil || in.Symbol=="" || in.Price<=0 { write(w,400,map[string]string{"error":"valid symbol and positive price are required"}); return }; a:=risk.Assess(in); s.mu.Lock(); s.latest[in.Symbol]=a; s.mu.Unlock(); s.assessed.Inc(); if a.Level=="CRITICAL"{s.critical.Inc()}; b,_:=json.Marshal(a); s.hub.Broadcast(b); write(w,200,a) }
func (s *Server) risks(w http.ResponseWriter,r *http.Request){ s.mu.RLock(); out:=make([]risk.Assessment,0,len(s.latest)); for _,a:=range s.latest{out=append(out,a)}; s.mu.RUnlock(); write(w,200,out) }
var upgrader=websocket.Upgrader{CheckOrigin:func(r *http.Request)bool{return true}}
func (s *Server) stream(w http.ResponseWriter,r *http.Request){ c,err:=upgrader.Upgrade(w,r,nil); if err!=nil{return}; defer c.Close(); ch:=s.hub.Subscribe(); defer s.hub.Unsubscribe(ch); for b:=range ch { if err:=c.WriteMessage(websocket.TextMessage,b); err!=nil{return} } }
func write(w http.ResponseWriter,status int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);if err:=json.NewEncoder(w).Encode(v);err!=nil{slog.Error("encode response","error",err)}}
func withCORS(next http.Handler) http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){w.Header().Set("Access-Control-Allow-Origin","*");w.Header().Set("Access-Control-Allow-Headers","Content-Type");if r.Method==http.MethodOptions{w.WriteHeader(204);return};next.ServeHTTP(w,r)})}
