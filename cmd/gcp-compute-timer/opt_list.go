package main

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	log "github.com/sirupsen/logrus"

	gcpt "github.com/jvzantvoort/gcp-compute-timer"
	"github.com/jvzantvoort/gcp-compute-timer/config"
)

type ListSubCmd struct {
	printfull bool
	verbose   bool
}

func (*ListSubCmd) Name() string {
	return "list"
}

func (*ListSubCmd) Synopsis() string {
	return "Check instances"
}

func (*ListSubCmd) Usage() string {
	return `
Check instances
`
	// msgstr, err := gcpt.Asset("messages/usage_check")
	// if err != nil {
	//         log.Error(err)
	//         msgstr = []byte("undefined")
	// }
	// return string(msgstr)
}

func (c *ListSubCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.printfull, "f", false, "Print full")
	f.BoolVar(&c.verbose, "v", false, "Verbose logging")
}

func (c *ListSubCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	configuration := config.NewConfiguration()

	gcp_project := configuration.GCP.Project
	gcp_zone := configuration.GCP.Zone
	gcp_bucket := configuration.GCP.Bucket

	if c.verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugln("Start")

	gcpt.PrintCheck(gcp_project, gcp_zone, gcp_bucket, c.printfull)

	log.Debugln("End")

	return subcommands.ExitSuccess
}
