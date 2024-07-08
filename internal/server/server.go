package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

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
	stopCh  chan struct{}
}

func (s *server) Run(ctx context.Context) error{
	// Initialize cities
	err := s.db.InitializeCities()
	if err != nil {
		return fmt.Errorf("failed to initialize cities: %w", err)
	}

	// Update weather forecast
	s.stopCh = make(chan struct{})
	go s.startWeatherUpdateBackgroundProcess()
	log.Println("init complite")
	
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
		close(s.stopCh)
		if err:=ctx.Err();err!=nil{
			return fmt.Errorf("context faild: %w", err)
		}
			
	}
	return nil
}
func (s *server) Close() error{
	s.db.DB.Close(context.Background())
	close(s.stopCh)
	return s.srv.Close()
}
func (s *server) startWeatherUpdateBackgroundProcess() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			log.Println("Stopping background weather update process")
			return
		case <-ticker.C:
			err := s.db.UpdateWeatherForecast()
			if err != nil {
				log.Printf("failed to update weather forecast: %v", err)
			} else {
				log.Println("Weather forecast updated successfully")
			}
		}
	}
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
	sv.setupRoutes()
	return &sv,nil
	
}

func (s *server)setupRoutes(){
	mx :=http.NewServeMux()

	mx.HandleFunc("/allCyties",s.handleAllCyties)
	mx.HandleFunc("/hortInfoCity",s.handleHortInfoCity)
	mx.HandleFunc("/fullInfoCity",s.handleFullInfoCity)
	// mx.HandleFunc("/favoritCityInfo",s.handleFavoritCityInfo)


	s.srv.Handler = mx
	
}

func (s *server) handleAllCyties(w http.ResponseWriter, r *http.Request) {
	cities, err := s.db.GetAllCities()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get cities: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cities); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}
func (s *server) handleHortInfoCity(w http.ResponseWriter, r *http.Request) {
	cityName := r.URL.Query().Get("city")
	if cityName == "" {
		http.Error(w, "Missing city parameter", http.StatusBadRequest)
		return
	}

	info, err := s.db.GetShortInfoCity(cityName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get city info: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}
func (s *server) handleFullInfoCity(w http.ResponseWriter, r *http.Request) {
	cityName := r.URL.Query().Get("city")
	dateStr := r.URL.Query().Get("date")
	if cityName == "" || dateStr == "" {
		http.Error(w, "Missing city or date parameter", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid date format: %v", err), http.StatusBadRequest)
		return
	}

	info, err := s.db.GetFullInfoCity(cityName, date)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get city info: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}
