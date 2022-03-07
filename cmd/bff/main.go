package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/bff"
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
