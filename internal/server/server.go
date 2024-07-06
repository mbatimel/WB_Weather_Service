package server

import (
	"context"
	"errors"
	"fmt"

	"net"
	"net/http"

	"github.com/mbatimel/WB_Weather_Service/internal/config"
	"github.com/mbatimel/WB_Weather_Service/internal/repo"
)
var ErrChannelClosed = errors.New("channel is closed")
type Server interface {
	Run(ctx context.Context) error
	Close() error

}
type server struct {
	srv *http.Server
	db *repo.DataBase
}

func (s *server) Run(ctx context.Context) error{
	ch:=make(chan error, 1)
	defer close(ch)
	go func(){
		ch <- s.srv.ListenAndServe()
	}()
	select  {
	case err, ok := <-ch:
		if !ok{
			return ErrChannelClosed
		}
		if err != nil{
			return fmt.Errorf("failed to listen and serve: %w", err)
		}
	case <-ctx.Done():
		if err:=ctx.Err();err!=nil{
			return fmt.Errorf("context faild: %w", err)
		}
			
	}
	return nil
}
func (s *server) Close() error{
	return nil
}

func NewServerConfig(cfg config.Config) (Server, error){
	srv:= http.Server{
		Addr: net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
	}
	dataBase, err := repo.SetConfigs("config/config.yaml")
	if err != nil {
		return nil, err
	}
	dataBase.ConnectToDataBase()
	
	
	sv:=server{
		srv :&srv,
		db : dataBase,
	}
	// sv.setupRoutes()
	return &sv,nil
	
}

// func (s *server)setupRoutes(){
// 	mx :=http.NewServeMux()

// 	mx.HandleFunc("/cap",s.handleCap)
// 	mx.HandleFunc("/add",s.handleAdd)
// 	mx.HandleFunc("/clear",s.handleClear)
// 	mx.HandleFunc("/addttl",s.handleAddWithTTL)
// 	mx.HandleFunc("/get",s.handleGet)
// 	mx.HandleFunc("/remove",s.handleRemove)

// 	s.srv.Handler = mx
	
// }

// func (s *server) handleCap(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, "Cache capacity: %d", s.cache.Cap())
// }

// func (s *server) handleAdd(w http.ResponseWriter, r *http.Request) {
// 	key := r.URL.Query().Get("key")
// 	value := r.URL.Query().Get("value")
// 	if key == "" || value == "" {
// 		http.Error(w, "Missing key or value", http.StatusBadRequest)
// 		return
// 	}
// 	s.cache.Add(key, value)
// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, "Added key %s with value %s", key, value)
// }

// func (s *server) handleClear(w http.ResponseWriter, r *http.Request) {
// 	if err := s.cache.Clear(); err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to clear cache: %v", err), http.StatusInternalServerError)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintln(w, "Cache cleared")
// }

// func (s *server) handleAddWithTTL(w http.ResponseWriter, r *http.Request) {
// 	key := r.URL.Query().Get("key")
// 	value := r.URL.Query().Get("value")
// 	ttlStr := r.URL.Query().Get("ttl")
// 	if key == "" || value == "" || ttlStr == "" {
// 		http.Error(w, "Missing key, value or ttl", http.StatusBadRequest)
// 		return
// 	}
// 	ttl, err := time.ParseDuration(ttlStr)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Invalid ttl: %v", err), http.StatusBadRequest)
// 		return
// 	}
// 	s.cache.AddWithTTL(key, value, ttl)
// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, "Added key %s with value %s and ttl %s", key, value, ttl)
// }

// func (s *server) handleGet(w http.ResponseWriter, r *http.Request) {
// 	key := r.URL.Query().Get("key")
// 	if key == "" {
// 		http.Error(w, "Missing key", http.StatusBadRequest)
// 		return
// 	}
// 	value, ok := s.cache.Get(key)
// 	if !ok {
// 		http.Error(w, "Key not found", http.StatusNotFound)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, "Key: %s, Value: %s", key, value)
// }

// func (s *server) handleRemove(w http.ResponseWriter, r *http.Request) {
// 	key := r.URL.Query().Get("key")
// 	if key == "" {
// 		http.Error(w, "Missing key", http.StatusBadRequest)
// 		return
// 	}
// 	s.cache.Remove(key)
// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, "Removed key %s", key)
// }