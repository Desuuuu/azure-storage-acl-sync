package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Desuuuu/azure-storage-acl-sync/azure"
	"github.com/urfave/cli/v2"
)

//go:embed VERSION
var version string

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	app := cli.NewApp()
	app.Name = "azure-storage-acl-sync"
	app.Usage = "Synchronize Azure storage account IP ACL with Azure service IPs"
	app.Version = version
	app.EnableBashCompletion = true

	app.Flags = flags()

	app.Action = func(c *cli.Context) error {
		return Run(c.Context, buildConfig(c))
	}

	err := app.RunContext(ctx, os.Args)
	if err != nil {
		stop()

		log.Fatal(err)
	}
}

func Run(ctx context.Context, config Config) error {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}

	rules, err := azure.GetIPRulesForServices(ctx, credential, config.SubscriptionID, config.Location, config.Services...)
	if err != nil {
		return err
	}

	allow := "Allow"

	for _, rule := range config.ExtraAllowRules {
		rule := rule

		rules = append(rules, &armstorage.IPRule{
			IPAddressOrRange: &rule,
			Action:           &allow,
		})
	}

	deny := "Deny"

	for _, rule := range config.ExtraDenyRules {
		rule := rule

		rules = append(rules, &armstorage.IPRule{
			IPAddressOrRange: &rule,
			Action:           &deny,
		})
	}

	rules = azure.FilterRules(rules)

	if config.DryRun {
		fmt.Fprintln(os.Stderr, "(DRY-RUN) The following IP rules would be applied:")

		for _, rule := range rules {
			fmt.Printf("%s (%s)\n", *rule.IPAddressOrRange, *rule.Action)
		}

		return nil
	}

	return azure.UpdateStorageAccountIPRules(ctx, credential, config.SubscriptionID, config.ResourceGroup, config.StorageAccount, rules)
}
