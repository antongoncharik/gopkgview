package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/antongoncharik/gopkgviewer/internal/graph"
	"github.com/carlmjohnson/versioninfo"
	"github.com/urfave/cli/v2"
)

func main() {
	//go:embed frontend/dist/*
	// var frontend embed.FS

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
			gomod := cCtx.String("gomod")
			skipBrowser := cCtx.Bool("skip-browser")

			log.Println("creating graph...")
			pkgGraph, err := graph.New()
			if err != nil {
				return fmt.Errorf("failed to build graph: %w", err)
			}
			log.Println("graph created")
			log.Println(pkgGraph)

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("error:")
		fmt.Printf(" > %+v\n", err)
	}
}
