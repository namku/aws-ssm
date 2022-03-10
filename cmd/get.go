/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/namku/aws-ssm/pkg"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files

to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// flags for custom aws config
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")

		ssmClient := pkg.NewSSM(profile, region)
		results, err := ssmClient.SSM.GetParameters(context.TODO(), &ssm.GetParametersInput{
			Names: args,
		})
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}

		for _, n := range results.Parameters {
			fmt.Println(*n.Value)
		}
	},
}

func init() {
	getCmd.PersistentFlags().String("profile", "default", "AWS configuration profile")
	getCmd.PersistentFlags().String("region", "", "AWS configuration profile")

	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
