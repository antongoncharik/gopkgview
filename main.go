package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/antonhancharyk/gopkgviewer/internal/graph"
	"github.com/carlmjohnson/versioninfo"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
)

//go:embed frontend/dist/*
var frontend embed.FS

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	go func() {
		<-ctx.Done()
		time.Sleep(3 * time.Second)

		log.Print("force exit")
		os.Exit(1)
	}()

	app := &cli.App{
		Name:    "gopkgviewer",
		Usage:   "Show dependencies of a Go package",
		Version: versioninfo.Short(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "gomod",
				EnvVars: []string{"GO_PKGVIEWER_GOMOD"},
				Usage:   "Path to go.mod file to detect external dependencies",
			},
			&cli.StringFlag{
				Name:        "addr",
				EnvVars:     []string{"GO_PKGVIEWER_ADDR"},
				Usage:       "Address to listen on",
				DefaultText: ":0",
			},
			&cli.BoolFlag{
				Name:    "skip-browser",
				EnvVars: []string{"GO_PKGVIEW_SKIP_BROWSER"},
				Usage:   "Don't open browser on start",
			},
		},
		Action: func(cCtx *cli.Context) error {
			addr := cCtx.String("addr")
			// gomod := cCtx.String("gomod")
			skipBrowser := cCtx.Bool("skip-browser")

			log.Print("creating graph...")
			pkgGraph := graph.New()
			// if err != nil {
			// 	return fmt.Errorf("failed to build graph: %v", err)
			// }

			graphData := map[string]interface{}{
				"nodes": pkgGraph.Nodes,
				"edges": pkgGraph.Edges,
			}

			graphJSON, err := json.Marshal(graphData)
			if err != nil {
				return fmt.Errorf("failed to marshal to JSON: %v", err)
			}

			handler := func(data []byte) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Content-Type", "application/json")
					if _, err := w.Write(data); err != nil {
						log.Printf("failed to write JSON: %v", err)
					}
				}
			}

			fsys, err := fs.Sub(frontend, "frontend/dist")
			if err != nil {
				return fmt.Errorf("failed to get frontend subdirectory: %v", err)
			}

			mux := http.NewServeMux()
			mux.Handle("/data", handler(graphJSON))
			mux.Handle("/", http.FileServer(http.FS(fsys)))

			listener, err := net.Listen("tcp", addr)
			if err != nil {
				return fmt.Errorf("failed to listen: %v", err)
			}
			defer listener.Close()

			server := &http.Server{Handler: mux}
			go func() {
				log.Print("starting server on ", listener.Addr())

				if !skipBrowser {
					webAddr := "http://" + listener.Addr().String()
					log.Print("opening browser on ", webAddr)
					if err := browser.OpenURL(webAddr); err != nil {
						log.Printf("failed to open browser: %v", err)
					}
				}

				if err := server.Serve(listener); err != http.ErrServerClosed {
					log.Fatalf("serve(): %v", err)
				}
			}()

			<-ctx.Done()
			log.Print("shutting down...")

			return server.Shutdown(ctx)
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("error:")
		fmt.Printf(" > %+v\n", err)
	}
}
