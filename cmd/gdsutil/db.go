package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/trisacrypto/directory/pkg/utils/wire"
	"github.com/urfave/cli/v2"
)

//===========================================================================
// Database Functions
//===========================================================================

func dbUsage(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	counters := []struct {
		namespace string
		count     func(context.Context) (uint64, error)
	}{
		{wire.NamespaceVASPs, db.CountVASPs},
		{wire.NamespaceCertReqs, db.CountCertReqs},
		{wire.NamespaceCerts, db.CountCerts},
		{wire.NamespaceAnnouncements, db.CountAnnouncementMonths},
		{wire.NamespaceActivities, db.CountActivityMonth},
		{wire.NamespaceOrganizations, db.CountOrganizations},
		{wire.NamespaceContacts, db.CountContacts},
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Namespace\tObjects")
	for _, counter := range counters {
		var count uint64
		if count, err = counter.count(ctx); err != nil {
			return cli.Exit(err, 1)
		}
		fmt.Fprintf(w, "%s\t%d\n", counter.namespace, count)
	}
	w.Flush()
	return nil
}
