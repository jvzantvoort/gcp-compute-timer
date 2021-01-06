package main

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	log "github.com/sirupsen/logrus"

	gcpt "github.com/jvzantvoort/gcp-compute-timer"
	"github.com/jvzantvoort/gcp-compute-timer/config"
)

type CheckSubCmd struct {
	noexec    bool
	verbose   bool
}

func (*CheckSubCmd) Name() string {
	return "check"
}

func (*CheckSubCmd) Synopsis() string {
	return "Check instances"
}

func (*CheckSubCmd) Usage() string {
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

func (c *CheckSubCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.noexec, "n", false, "Do not execute only display a list of images and actions")
	f.BoolVar(&c.verbose, "v", false, "Verbose logging")
}

func (c *CheckSubCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	configuration := config.NewConfiguration()

	gcp_project := configuration.GCP.Project
	gcp_zone := configuration.GCP.Zone
	gcp_bucket := configuration.GCP.Bucket

	if c.verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugln("Start")

	gcpt.RunCheck(gcp_project, gcp_zone, gcp_bucket, c.noexec)

	log.Debugln("End")

	return subcommands.ExitSuccess
}
