package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/bff"
	"github.com/trisacrypto/directory/pkg/bff/api/v1"
	"github.com/trisacrypto/directory/pkg/bff/config"
	"github.com/urfave/cli/v2"
)

func main() {
	// Load the dotenv file if it exists
	godotenv.Load()

	// Create the CLI application
	app := &cli.App{
		Name:    "gds-bff",
		Version: pkg.Version(),
		Usage:   "a backend for front-end for the GDS service",
		Flags:   []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:     "serve",
				Usage:    "run the gds bff server",
				Category: "server",
				Action:   serve,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "addr",
						Aliases: []string{"a"},
						Usage:   "the address and port to bind the server on",
						EnvVars: []string{"GDS_BFF_BIND_ADDR"},
					},
				},
			},
			{
				Name:     "status",
				Usage:    "send a status check to the BFF server",
				Category: "client",
				Action:   status,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "url",
						Aliases: []string{"u", "endpoint"},
						Usage:   "specify the URL to connect to the BFF server on",
						EnvVars: []string{"GDS_BFF_CLIENT_URL"},
						Value:   "http://localhost:4437",
					},
					&cli.BoolFlag{
						Name:    "nogds",
						Aliases: []string{"no-gds", "G"},
						Usage:   "health check the BFF without requesting GDS status",
						Value:   false,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

//===========================================================================
// CLI Actions
//===========================================================================

// Serve the GDS BFF service
func serve(c *cli.Context) (err error) {
	var conf config.Config
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}

	if addr := c.String("addr"); addr != "" {
		conf.BindAddr = addr
	}

	var srv *bff.Server
	if srv, err = bff.New(conf); err != nil {
		return cli.Exit(err, 1)
	}

	if err = srv.Serve(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

// Status checks if the GDS BFF is up
func status(c *cli.Context) (err error) {
	var client api.BFFClient
	if client, err = api.New(c.String("url")); err != nil {
		return cli.Exit(err, 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var params *api.StatusParams
	if c.Bool("nogds") {
		params = &api.StatusParams{
			NoGDS: true,
		}
	}

	var rep *api.StatusReply
	if rep, err = client.Status(ctx, params); err != nil {
		return cli.Exit(err, 1)
	}

	return printJSON(rep)
}

//===========================================================================
// Helper Functions
//===========================================================================

func printJSON(msg interface{}) (err error) {
	var data []byte
	if data, err = json.MarshalIndent(msg, "", "  "); err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Println(string(data))
	return nil
}
