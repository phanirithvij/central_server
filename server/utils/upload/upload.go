package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"github.com/tus/tusd/pkg/filestore"
	tusd "github.com/tus/tusd/pkg/handler"
)

// go:generate rm -f pkged.go
// go:generate pkger -o server/utils/upload
// go:generate mv pkged.go pkged_g.go

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	// TODO check if path has .. i.e relative routes and ban IP
	// https://github.com/mrichman/godnsbl
	// https://github.com/jpillora/ipfilter

	// prepend the path with the path to the static directory
	path := filepath.Join(h.staticPath, r.URL.Path)

	// check whether a file exists at the given path
	_, err := pkger.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		// http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		log.Println("File", h.indexPath)
		file, err := pkger.Open(h.indexPath)
		if err != nil {
			http.Error(w, "file "+r.URL.Path+" does not exist", http.StatusNotFound)
			return
		}
		http.FileServer(file).ServeHTTP(w, r)
		return
	} else if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "file "+r.URL.Path+" does not exist", http.StatusNotFound)
			return
		}
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		// log.Println(err)
		// return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(pkger.Dir(h.staticPath)).ServeHTTP(w, r)
}

// https://github.com/tus/tusd/blob/master/docs/usage-package.md
func main() {
	// https://github.com/gorilla/mux/#graceful-shutdown
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	os.Mkdir("uploads", os.ModeDir)
	// Create a new FileStore instance which is responsible for
	// storing the uploaded file on disk in the specified directory.
	// This path _must_ exist before tusd will store uploads in it.
	// If you want to save them on a different medium, for example
	// a remote FTP server, you can implement your own storage backend
	// by implementing the tusd.DataStore interface.
	store := filestore.FileStore{
		Path: "./uploads",
	}

	// A storage backend for tusd may consist of multiple different parts which
	// handle upload creation, locking, termination and so on. The composer is a
	// place where all those separated pieces are joined together. In this example
	// we only use the file store but you may plug in multiple.
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	// Create a new HTTP handler for the tusd server by providing a configuration.
	// The StoreComposer property must be set to allow the handler to function.
	uploadHandler, err := tusd.NewHandler(tusd.Config{
		BasePath:              "/files/",
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})
	if err != nil {
		panic(fmt.Errorf("unable to create handler: %s", err))
	}

	// Start another goroutine for receiving events from the handler whenever
	// an upload is completed. The event will contains details about the upload
	// itself and the relevant HTTP request.
	go func() {
		for {
			event := <-uploadHandler.CompleteUploads
			fmt.Printf("Upload %s finished\n", event.Upload.ID)
		}
	}()

	// Right now, nothing has happened since we need to start the HTTP server on
	// our own. In the end, tusd will start listening on and accept request at
	// http://localhost:8080/files
	r := mux.NewRouter()

	// https://github.com/gorilla/mux#serving-single-page-applications
	spa := &spaHandler{
		staticPath: "/server/utils/upload/assets",
		indexPath:  "/server/utils/upload/assets/index.html",
	}

	pkger.Include("/server/utils/upload/assets")

	// https://stackoverflow.com/a/27946132/8608146
	r.PathPrefix("/upload").Handler(http.StripPrefix("/upload", spa))
	r.PathPrefix("/files/").Handler(http.StripPrefix("/files/", uploadHandler))

	log.Println("Starting on localhost:8080")
	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(http.ErrServerClosed, err) {
				return
			}
			log.Panicln(fmt.Errorf("unable to listen: %s", err))
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

}
