package main

import (
	"github.com/urfave/cli/v2"
)

type Config struct {
	SubscriptionID  string
	Services        []string
	Location        string
	ResourceGroup   string
	StorageAccount  string
	ExtraAllowRules []string
	ExtraDenyRules  []string
	DryRun          bool
}

func flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "subscription-id",
			Usage:    "Azure subscription ID",
			EnvVars:  []string{"AZURE_SUBSCRIPTION_ID"},
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:    "services",
			Usage:   "Azure services to retrieve IPs from",
			EnvVars: []string{"AZURE_SERVICES"},
			Value:   cli.NewStringSlice("AzureFrontDoor.Backend"),
		},
		&cli.StringFlag{
			Name:    "location",
			Usage:   "Azure location to retrieve IPs for",
			EnvVars: []string{"AZURE_LOCATION"},
			Value:   "westus",
		},
		&cli.StringFlag{
			Name:    "resource-group",
			Usage:   "storage account resource group",
			EnvVars: []string{"AZURE_RESOURCE_GROUP"},
		},
		&cli.StringFlag{
			Name:    "storage-account",
			Usage:   "storage account name",
			EnvVars: []string{"AZURE_STORAGE_ACCOUNT"},
		},
		&cli.StringSliceFlag{
			Name:    "extra-allow-rules",
			Usage:   "additional allow IP rules",
			EnvVars: []string{"EXTRA_ALLOW_RULES"},
			Value:   cli.NewStringSlice("168.63.129.16", "169.254.169.254"),
		},
		&cli.StringSliceFlag{
			Name:    "extra-deny-rules",
			Usage:   "additional deny IP rules",
			EnvVars: []string{"EXTRA_DENY_RULES"},
		},
		&cli.BoolFlag{
			Name:    "dry-run",
			Usage:   "only print the IP rules that would be applied",
			EnvVars: []string{"DRY_RUN"},
		},
	}
}

func buildConfig(c *cli.Context) Config {
	return Config{
		SubscriptionID:  c.String("subscription-id"),
		Services:        c.StringSlice("services"),
		Location:        c.String("location"),
		ResourceGroup:   c.String("resource-group"),
		StorageAccount:  c.String("storage-account"),
		ExtraAllowRules: c.StringSlice("extra-allow-rules"),
		ExtraDenyRules:  c.StringSlice("extra-deny-rules"),
		DryRun:          c.Bool("dry-run"),
	}
}
