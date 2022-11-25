package cli

import (
	"fmt"
	"log"

	"github.com/BishopFox/cloudfox/azure"
	"github.com/BishopFox/cloudfox/globals"
	"github.com/spf13/cobra"
)

var (
	AzTenantID        string
	AzSubscriptionID  string
	AzRGName          string
	AzOutputFormat    string
	AzOutputDirectory string
	AzVerbosity       int
	AzCommands        = &cobra.Command{
		Use:     "azure",
		Aliases: []string{"az"},
		Long:    `See "Available Commands" for Azure Modules`,
		Short:   "See \"Available Commands\" for Azure Modules",

		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			fmt.Println()
			fmt.Println("Your Azure CLI is currently authenticated to the following scopes:")
			azure.PrintAvailableScopeFull(azure.GetAvailableScope())
		},
	}

	AzInstancesCommand = &cobra.Command{
		Use:     "instances",
		Aliases: []string{},
		Short:   "Enumerates Azure Compute Instances",
		Long: `
Select scope from interactive menu:
./cloudfox az instances

Enumerate VMs for all resource groups in a subscription:
./cloudfox az instances --subscription SUBSCRIPTION_NAME

Enumerate VMs from a specific resource group:
./cloudfox az instances -s SUBSCRIPTION_NAME -g RESOURCE_GROUP_NAME`,
		Run: func(cmd *cobra.Command, args []string) {
			err := azure.AzInstancesCommand(AzSubscriptionID, AzRGName, AzOutputFormat, AzVerbosity)
			if err != nil {
				log.Fatalf("[%s] %s", globals.AZ_INTANCES_MODULE_NAME, err)
			}
		},
	}

	AzRbacCommand = &cobra.Command{
		Use:     "rbac",
		Aliases: []string{},
		Short:   "Display all role assignemts for all Azure principals",
		Long: `
Select scope from interactive menu:
./cloudfox az rbac

Enumerate role assignments for a specific subscriptions:
./cloudfox az rbac --subscription SUBSCRIPTION_NAME
`,
		Run: func(cmd *cobra.Command, args []string) {
			err := azure.AzRbacCommand(azure.CloudFoxRBACclient{}, AzTenantID, AzSubscriptionID, AzOutputFormat, AzVerbosity)
			if err != nil {
				log.Fatalf("[%s] %s", globals.AZ_RBAC_MODULE_NAME, err)
			}
		},
	}
)

func init() {
	// Global flags
	AzCommands.PersistentFlags().StringVarP(&AzOutputFormat, "output", "o", "all", "[\"table\" | \"csv\" | \"all\" ]")
	AzCommands.PersistentFlags().IntVarP(&AzVerbosity, "verbosity", "v", 1, "1 = Print control messages only\n2 = Print control messages, module output\n3 = Print control messages, module output, and loot file output\n")
	AzCommands.PersistentFlags().StringVarP(&AzTenantID, "tenant", "t", "", "Tenant name")
	AzCommands.PersistentFlags().StringVarP(&AzSubscriptionID, "subscription", "s", "", "Subscription Name")
	AzCommands.PersistentFlags().StringVarP(&AzRGName, "resource-group", "g", "", "Resource Group name")

	AzCommands.AddCommand(AzInstancesCommand, AzRbacCommand)
}
