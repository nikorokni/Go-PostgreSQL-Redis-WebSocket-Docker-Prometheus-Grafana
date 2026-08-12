package hub

import "sync"

type Hub struct { mu sync.RWMutex; clients map[chan []byte]struct{} }
func New() *Hub { return &Hub{clients:map[chan []byte]struct{}{}} }
func (h *Hub) Subscribe() chan []byte { c:=make(chan []byte,16); h.mu.Lock(); h.clients[c]=struct{}{}; h.mu.Unlock(); return c }
func (h *Hub) Unsubscribe(c chan []byte) { h.mu.Lock(); if _,ok:=h.clients[c]; ok { delete(h.clients,c); close(c) }; h.mu.Unlock() }
func (h *Hub) Broadcast(b []byte) { h.mu.RLock(); defer h.mu.RUnlock(); for c:=range h.clients { select { case c<-b: default: } } }
