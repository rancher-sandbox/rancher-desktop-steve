// Command dashboard runs the Rancher Desktop dashboard server, including both
// steve (the API backend) and the frontend.
// Once started, this listens on a random port on localhost, and prints the port
// number to stdout before immediately closing it.  The caller can then connect
// to the dashboard on that port.
package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"mime"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rancher/steve/pkg/debug"
	stevecli "github.com/rancher/steve/pkg/server/cli"
	"github.com/rancher/steve/pkg/version"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

//go:embed dashboard/*
var dashboardFiles embed.FS

var (
	config      stevecli.Config
	debugconfig debug.Config
)

func main() {
	app := cli.NewApp()
	app.Name = "rancher-desktop-app-dashboard"
	app.Version = version.FriendlyVersion()
	app.Usage = "Rancher Desktop dashboard server"
	app.Description = "This is part of Rancher Desktop 2.x and should not be run manually."
	app.Flags = append(
		stevecli.Flags(&config),
		debug.Flags(&debugconfig)...)
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(_ *cli.Context) error {
	// Set up steve; however, we do not tell it to listen by itself.
	ctx := signals.SetupSignalContext()
	debugconfig.MustSetupDebug()
	s, err := config.ToServer(ctx)
	if err != nil {
		return err
	}
	s.StartAggregation(ctx)

	// Listen on a port on localhost.
	localhost := net.IPv4(127, 0, 0, 1)
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: localhost, Port: 0})
	if err != nil {
		return err
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port

	// Set up the rewrite handler, to fix up source code that has hard-coded the
	// default port number.
	rewriter, err := rewriteHandler(port)
	if err != nil {
		return fmt.Errorf("could not set up rewrite handler: %w", err)
	}

	// Print the port number to stdout, and then close it, so that the caller
	// can safely read it.
	if _, err := fmt.Printf("%d\n", port); err != nil {
		return err
	}
	if err := os.Stdout.Close(); err != nil {
		return err
	}
	os.Stdout = os.Stderr
	logrus.Infof("Steve is listening on port %d", port)

	// Set up the mux.
	mux := http.NewServeMux()
	// Set the the default handler, to use the steve APIs.
	mux.Handle("/", s)
	// Set the a handler for the root, to redirect to the local cluster explorer.
	mux.Handle("/{$}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/c/local/explorer", http.StatusMovedPermanently)
	}))
	// Set up a handler for any cluster explorer, which gets passed to the rewriter.
	mux.Handle("/c/", rewriter)

	// Enumerate all (top-level) entries in the dashboard FS, and handle them
	// individually.  This is required so that all unknown paths get passed to the
	// steve API handler instead.
	files, err := dashboardFiles.ReadDir("dashboard")
	if err != nil {
		return fmt.Errorf("could not read front end file: %w", err)
	}
	for _, f := range files {
		if f.IsDir() {
			mux.Handle("/"+f.Name()+"/", http.HandlerFunc(rewriter))
		} else {
			mux.Handle("/"+f.Name(), http.HandlerFunc(rewriter))
		}
	}
	// Handle the steve-port API, which is no longer used.
	mux.Handle("/api/steve-port", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%d", port)
	}))

	// Actually set up the HTTP server, plus a goroutine to shut down once the
	// context is done.  Afterwards, start the server.
	server := &http.Server{Handler: mux}
	go func() {
		<-ctx.Done()
		// Wait some time for steve stuff to shut down
		<-time.After(10 * time.Millisecond)
		// Shutdown the HTTP server with a small timeout.
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(timeoutCtx); err != nil {
			logrus.Errorf("error shutting down server: %v", err)
		}
	}()
	if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Serve the request from static files, rewriting any instance of the default
// port to the actual port number.
func rewriteHandler(port int) (http.HandlerFunc, error) {
	subFS, err := fs.Sub(dashboardFiles, "dashboard")
	if err != nil {
		return nil, fmt.Errorf("could not create sub filesystem: %w", err)
	}
	if _, ok := subFS.(fs.ReadFileFS); !ok {
		return nil, fmt.Errorf("sub filesystem does not support ReadFile")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// Calculate the path inside the dashboard FS to serve.
		actualPath := strings.TrimPrefix(r.URL.Path, "/")
		if _, err := subFS.Open(actualPath); errors.Is(err, fs.ErrNotExist) {
			// If the file does not exist, serve the index file because this is
			// a single-page application.
			actualPath = "index.html"
		}
		mimeType := mime.TypeByExtension(path.Ext(actualPath))
		if !strings.HasPrefix(mimeType, "text/") {
			// For non-text files, just serve them directly without rewriting.
			http.FileServerFS(subFS).ServeHTTP(w, r)
			return
		}

		// This is a text file that needs rewriting; read it, replace any
		// relevant strings, and serve it.  At the moment, embed.FS just stores
		// the whole thing in memory anyway, so this should not be a problem.
		contents, err := subFS.(fs.ReadFileFS).ReadFile(actualPath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				// Since we munged actualPath above, this should only happen if
				// "index.html" is itself missing.
				http.NotFound(w, r)
			} else {
				http.Error(
					w,
					fmt.Sprintf("could not read %s: %v", actualPath, err),
					http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", mimeType)
		replacer := strings.NewReplacer(
			"http://127.0.0.1:6120/",
			fmt.Sprintf("http://127.0.0.1:%d/", port))
		if _, err := replacer.WriteString(w, string(contents)); err != nil {
			logrus.Errorf("could not write response: %v", err)
		}
	}, nil
}
