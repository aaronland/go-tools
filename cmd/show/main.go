package main

/*

> go run cmd/gh2b/main.go -mode geojson 9q8yy | go run cmd/show/main.go -stdin
Features are viewable at http://localhost:8080

> go run cmd/gh2b/main.go -mode geojson 9q8yy 9q5ct | go run cmd/show/main.go -stdin
Features are viewable at http://localhost:8080

*/

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/paulmach/orb/geojson"
	"github.com/pkg/browser"
	"github.com/tidwall/gjson"
)

//go:embed *.html
var html_FS embed.FS

//go:embed css/*.css
var css_FS embed.FS

//go:embed javascript/*.js javascript/*.map
var js_FS embed.FS

type mapConfig struct {
	Provider string	`json:"provider"`
	TileURL string `json:"tile_url"`
}

func main() {

	var port int
	var stdin bool

	var map_provider string
	var map_tile_url string

	flag.StringVar(&map_provider, "map-provider", "leaflet", "")
	flag.StringVar(&map_tile_url, "map-tile-url", "https://tile.openstreetmap.org/{z}/{x}/{y}.png", "")
	
	flag.IntVar(&port, "port", 8080, "The port number to listen for requests on (on localhost). If 0 then a random port number will be chosen.")
	flag.BoolVar(&stdin, "stdin", false, "")

	flag.Parse()

	ctx := context.Background()

	fc := geojson.NewFeatureCollection()

	append_features := func(r io.Reader) error {

		body, err := io.ReadAll(r)

		if err != nil {
			return fmt.Errorf("Failed to read body, %w", err)
		}

		type_rsp := gjson.GetBytes(body, "type")

		switch type_rsp.String() {
		case "Feature":

			f, err := geojson.UnmarshalFeature(body)

			if err != nil {
				return fmt.Errorf("Failed to unmarshal Feature, %w", err)
			}

			fc.Append(f)

		case "FeatureCollection":

			other_fc, err := geojson.UnmarshalFeatureCollection(body)

			if err != nil {
				return fmt.Errorf("Failed to unmarshal record as FeatureCollection, %w", err)
			}

			for _, f := range other_fc.Features {
				fc.Append(f)
			}

		default:
			return fmt.Errorf("Invalid type, %s", type_rsp.String())
		}

		return nil
	}

	uris := flag.Args()

	if len(uris) == 1 && uris[0] == "-" {
		stdin = true
	}

	if stdin {

		err := append_features(os.Stdin)

		if err != nil {
			log.Fatalf("Failed to append features, %v", err)
		}

	} else {

		for _, path := range uris {

			r, err := os.Open(path)

			if err != nil {
				log.Fatalf("Failed to open %s for reading, %v", path, err)
			}

			defer r.Close()

			err = append_features(r)

			if err != nil {
				log.Fatalf("Failed to append features, %v", err)
			}
		}
	}

	data_handler := dataHandler(fc)

	map_cfg := &mapConfig{
		Provider: map_provider,
		TileURL: map_tile_url,
	}

	map_cfg_handler := mapConfigHandler(map_cfg)	
	
	html_fs := http.FS(html_FS)
	js_fs := http.FS(js_FS)
	css_fs := http.FS(css_FS)

	mux := http.NewServeMux()
	mux.Handle("/map.json", map_cfg_handler)	
	mux.Handle("/features.geojson", data_handler)

	mux.Handle("/css/", http.FileServer(css_fs))
	mux.Handle("/javascript/", http.FileServer(js_fs))
	mux.Handle("/", http.FileServer(html_fs))

	addr := fmt.Sprintf("localhost:%d", port)
	url := fmt.Sprintf("http://%s", addr)

	http_server := http.Server{
		Addr: addr,
	}

	http_server.Handler = mux

	done_ch := make(chan bool)
	err_ch := make(chan error)

	go func() {

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		slog.Info("Shutting server down")
		err := http_server.Shutdown(ctx)

		if err != nil {
			slog.Error("HTTP server shutdown error", "error", err)
		}

		close(done_ch)
	}()

	go func() {

		err := http_server.ListenAndServe()

		if err != nil {
			err_ch <- fmt.Errorf("Failed to start server, %w", err)
		}
	}()

	server_ready := false

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case err := <-err_ch:
			log.Fatalf("Received error starting server, %v", err)
		case <-ticker.C:

			rsp, err := http.Head(url)

			if err != nil {
				slog.Warn("HEAD request failed", "url", url, "error", err)
			} else {

				defer rsp.Body.Close()

				if rsp.StatusCode != 200 {
					slog.Warn("HEAD request did not return expected status code", "url", url, "code", rsp.StatusCode)
				} else {
					slog.Debug("HEAD request succeeded", "url", url)
					server_ready = true
				}
			}
		}

		if server_ready {
			break
		}
	}

	err := browser.OpenURL(url)

	if err != nil {
		log.Fatalf("Failed to open URL %s, %v", url, err)
	}

	fmt.Printf("Features are viewable at %s\n", url)

	<-done_ch
}

func dataHandler(fc *geojson.FeatureCollection) http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		enc_json, err := fc.MarshalJSON()

		if err != nil {
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-type", "application/json")
		rsp.Write(enc_json)
		return
	}

	return http.HandlerFunc(fn)
}

func mapConfigHandler(cfg *mapConfig) http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		rsp.Header().Set("Content-type", "application/json")
		
		enc := json.NewEncoder(rsp)
		err := enc.Encode(cfg)
		
		if err != nil {
			slog.Error("Failed to encode map config", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
		}
		
		return
	}

	return http.HandlerFunc(fn)
}
