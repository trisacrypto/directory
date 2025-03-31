package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

func migrateDirectory(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	iter := db.ListVASPs(ctx)
	defer iter.Release()

	updated := 0
	for iter.Next() {
		vasp, _ := iter.VASP()
		update := false

		switch vasp.RegisteredDirectory {
		case "vaspdirectory.net":
			vasp.RegisteredDirectory = "trisa.directory"
			update = true
		case "trisatest.net":
			vasp.RegisteredDirectory = "testnet.directory"
			update = true
		}

		if vasp.Website != "" {
			if u, err := url.Parse(vasp.Website); err == nil {
				if strings.HasSuffix(u.Hostname(), "vaspbot.net") {
					u.Host = strings.Replace(u.Host, "vaspbot.net", "vaspbot.com", 1)
					vasp.Website = u.String()
					update = true
				}
			}
		}

		if update {
			if err = db.UpdateVASP(ctx, vasp); err != nil {
				return cli.Exit(err, 1)
			}
			updated++
		}
	}

	if err = iter.Error(); err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Printf("updated %d vasp records\n", updated)
	return nil
}
