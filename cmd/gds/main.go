package main

import (
	"context"
	"log"
	"os"
	"text/tabwriter"

	"github.com/joho/godotenv"
	"github.com/trisacrypto/directory/pkg"
	"github.com/trisacrypto/directory/pkg/config"
	"github.com/trisacrypto/directory/pkg/server"
	"github.com/urfave/cli/v3"
	confire "go.rtnl.ai/confire/usage"
)

func main() {
	godotenv.Load()
	app := cli.Command{
		Name:    "gds",
		Usage:   "TRISA Global Directory Service for peer-to-peer matching and discovery",
		Version: pkg.Version(false),
		Commands: []*cli.Command{
			{
				Name:     "serve",
				Usage:    "Start the GDS server",
				Action:   serve,
				Category: "Server",
			},
			{
				Name:     "config",
				Usage:    "Print the gds configuration guide",
				Action:   usage,
				Category: "server",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "print in list mode instead of table mode",
					},
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func serve(ctx context.Context, c *cli.Command) (err error) {
	var srv *server.Server
	if srv, err = server.New(); err != nil {
		return cli.Exit(err, 1)
	}

	if err = srv.Serve(); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func usage(ctx context.Context, c *cli.Command) error {
	tabs := tabwriter.NewWriter(os.Stdout, 1, 0, 4, ' ', 0)
	format := confire.DefaultTableFormat
	if c.Bool("list") {
		format = confire.DefaultListFormat
	}

	var conf config.Config
	if err := confire.Usagef(config.Prefix, &conf, tabs, format); err != nil {
		return cli.Exit(err, 1)
	}
	tabs.Flush()
	return nil
}
