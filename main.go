package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"crawshaw.dev/jsonfile"
)

var favsHostPort = flag.String("f", "localhost:8080", "favorites service host:port")
var assetsHostPort = flag.String("a", "localhost:8080", "assets service host:port")
var assetsClientHostPort = flag.String("c", "", "assets service host:port for the fetch pool. Use with different machines")
var favoritesPath = flag.String("d", "./dbfav.json", "favorites file")
var assetSrvDelay = flag.Duration("w", 0, "asset server delay per request. Use for testing.")

// FavoritesStore is a wrapper for a map. Needed to use jsonfile.
type FavoritesStore struct {
	Favorites map[string][]*Favorite `json:"favorites"`
}

func main() {
	flag.Parse()

	if err := parseUsers(); err != nil {
		log.Fatal(err)
	}
	if err := decodeAssets(); err != nil {
		log.Fatal(err)
	}

	_, err := os.Stat(*favoritesPath)
	exists := err == nil || !os.IsNotExist(err)

	var favoritesDB *jsonfile.JSONFile[FavoritesStore]
	if exists {
		favoritesDB, err = jsonfile.Load[FavoritesStore](*favoritesPath)
	} else {
		favoritesDB, err = jsonfile.New[FavoritesStore](*favoritesPath)
	}
	if err != nil {
		log.Fatalf("failed to load favorites: %v", err)
	}

	fs := newFavoritesStorage()
	favoritesDB.Read(func(data *FavoritesStore) {
		if data.Favorites == nil {
			data.Favorites = make(map[string][]*Favorite)
		}
		fs.favorites = data.Favorites
	})

	// handle signals gracefully for shutdown. During shutdown we want
	// - stop the http get workers for assets, see fetcher.go
	// - write the favorites database to file, using jsonfile
	baseContext, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// for auth, used by both favorites and assets servers
	ah := &AuthService{defaultTokenizer()}
	protected := func(h http.HandlerFunc) http.Handler {
		return ah.withAuthorization(http.HandlerFunc(h))
	}

	serversInSameInstance := *favsHostPort == *assetsHostPort

	favoritesMux := http.NewServeMux()
	favoritesSrv := &http.Server{
		Addr:        *favsHostPort,
		BaseContext: func(_ net.Listener) context.Context { return baseContext },
		Handler:     favoritesMux,
	}

	// assetFetcher, see fetcher.go is a worker pool for GETs from assets server
	// It is server wide resource and needs a configured client to control the
	// number of connections to the assets server
	if *assetsClientHostPort == "" {
		*assetsClientHostPort = *assetsHostPort
	}
	assetFetcher := &AssetsFetcher{}
	assetFetcher.cli = &http.Client{}
	assetFetcher.srv, _ = url.Parse(fmt.Sprintf("http://%s/assets", *assetsClientHostPort))
	fssrv := newFavoritesService(fs, assetFetcher)
	favoritesMux.HandleFunc("POST /login", ah.loginHandler)
	favoritesMux.Handle("GET /verify", protected(ah.verifyTokenHandler))
	favoritesMux.Handle("GET /favorites", protected(fssrv.getAll))
	favoritesMux.Handle("POST /favorites", protected(fssrv.add))
	favoritesMux.Handle("GET /favorites/{id}", protected(fssrv.get))
	favoritesMux.Handle("PUT /favorites/{id}", protected(fssrv.update))
	favoritesMux.Handle("DELETE /favorites/{id}", protected(fssrv.del))

	var assetsMux *http.ServeMux
	var assetsSrv *http.Server
	if serversInSameInstance {
		assetsSrv = favoritesSrv
		assetsMux = favoritesMux
	} else {
		assetsMux = http.NewServeMux()
		assetsSrv = &http.Server{
			Addr:        *assetsHostPort,
			BaseContext: func(_ net.Listener) context.Context { return baseContext },
			Handler:     assetsMux,
		}
	}
	as := &AssetsServer{assets, *assetSrvDelay}
	assetsMux.HandleFunc("GET /assets/insights", http.NotFound)
	assetsMux.HandleFunc("GET /assets/insights/{$}", http.NotFound)
	assetsMux.Handle("GET /assets/insights/{key...}", protected(as.insightsHandler))
	assetsMux.HandleFunc("GET /assets/charts", http.NotFound)
	assetsMux.HandleFunc("GET /assets/charts/{$}", http.NotFound)
	assetsMux.Handle("GET /assets/charts/{key...}", protected(as.chartsHandler))
	assetsMux.HandleFunc("GET /assets/audiences", http.NotFound)
	assetsMux.HandleFunc("GET /assets/audiences/{$}", http.NotFound)
	assetsMux.Handle("GET /assets/audiences/{key...}", protected(as.audiencesHandler))

	log.Println("start services")
	go func() {
		err := favoritesSrv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("favorites server failed to start: %v", err)
		}
	}()
	if !serversInSameInstance {
		go func() {
			err := assetsSrv.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("favorites server failed to start: %v", err)
			}
		}()
	}

	<-baseContext.Done()
	log.Println("services shutdown on signal")

	if err := favoritesDB.Write(func(data *FavoritesStore) error {
		data.Favorites = fs.favorites
		return nil
	}); err != nil {
		log.Printf("failed to write favorites db: %v", err)
	}

	if err := favoritesSrv.Shutdown(context.Background()); err != nil {
		log.Printf("favorites server shutdown error: %v", err)
	}
	if !serversInSameInstance {
		if err := assetsSrv.Shutdown(context.Background()); err != nil {
			log.Printf("assets server shutdown error: %v", err)
		}
	}

	log.Println("service stopped")
}
