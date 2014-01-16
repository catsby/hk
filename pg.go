package main

import (
	"os"
	"text/tabwriter"

	"github.com/heroku/hk/postgresql"
)

var cmdPgInfo = &Command{
	Run:      runPgInfo,
	Usage:    "pg-info <dbname>",
	NeedsApp: true,
	Category: "app",
	Short:    "show heroku postgres database info",
	Long: `
Pg-info shows general information about a Heroku PostgreSQL
database.
`,
}

func runPgInfo(cmd *Command, args []string) {
	if len(args) != 1 {
		cmd.printUsage()
		os.Exit(2)
	}
	addon, err := client.AddonInfo(mustApp(), args[0])
	must(err)

	// fetch app's config concurrently in case we need to resolve DB names
	confch := make(chan map[string]string, 1)
	errch := make(chan error, 1)
	go func(appname string) {
		if app, err := client.AppInfo(appname); err != nil {
			errch <- err
		} else {
			appch <- app
		}
	}(name)
	for _ = range names {
		select {
		case err := <-errch:
			fmt.Fprintln(os.Stderr, err)
		case app := <-appch:
			if app != nil {
				apps = append(apps, *app)
			}
		}
	}

	db := pgclient.NewDB(addon.ProviderId, addon.Plan.Name)
	info, err := db.Info()
	must(err)

	printPgInfo(info.Info)
}

func printPgInfo(infos []postgresql.InfoEntry) {
	w := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)
	defer w.Flush()

	for _, ie := range infos {
		if len(ie.Values) == 0 {
			listRec(w, ie.Name+":", "none")
		} else {
			for n, val := range ie.Values {
				// TODO(bgentry): Resolve DB names from URLs if ResolveDBName=true
				switch n {
				case 0:
					listRec(w, ie.Name+":", val)
				default:
					listRec(w, "", val)
				}
			}
		}
	}
}

func DatabaseNameFromURL(url string, config map[string]string) string {
}
