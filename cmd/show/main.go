package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	// "net"
	"embed"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/paulmach/orb/geojson"
	"github.com/pkg/browser"
)

//go:embed *.html
var html_FS embed.FS

//go:embed css/*.css
var css_FS embed.FS

//go:embed javascript/*.js javascript/*.map
var js_FS embed.FS

func main() {

	var port int
	var stdin bool
	var is_featurecollection bool

	flag.IntVar(&port, "port", 8080, "The port number to listen for requests on (on localhost). If 0 then a random port number will be chosen.")
	flag.BoolVar(&stdin, "stdin", false, "")
	flag.BoolVar(&is_featurecollection, "featurecollection", false, "")

	flag.Parse()

	ctx := context.Background()

	fc := geojson.NewFeatureCollection()

	append_features := func(body []byte) error {

		if is_featurecollection {

			other_fc, err := geojson.UnmarshalFeatureCollection(body)

			if err != nil {
				log.Fatalf("Failed to unmarshal record as FeatureCollection, %w", err)
			}

			for _, f := range other_fc.Features {
				fc.Append(f)
			}

			return nil
		}

		f, err := geojson.UnmarshalFeature(body)

		if err != nil {
			log.Fatalf("Failed to unmarshal '%s' as Feature, %w", err)
		}

		fc.Append(f)
		return nil
	}

	if stdin {

		log.Fatal("Not implemented")
	} else {

		for _, path := range flag.Args() {

			r, err := os.Open(path)

			if err != nil {
				log.Fatalf("Failed to open %s for reading, %v", path, err)
			}

			defer r.Close()

			body, err := io.ReadAll(r)

			if err != nil {
				log.Fatalf("Failed to read '%s', %", path, err)
			}

			err = append_features(body)

			if err != nil {
				log.Fatalf("Failed to append features, %v", err)
			}
		}
	}

	data_handler := dataHandler(fc)

	html_fs := http.FS(html_FS)
	js_fs := http.FS(js_FS)
	css_fs := http.FS(css_FS)

	mux := http.NewServeMux()
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
