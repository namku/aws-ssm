/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

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

		bypath, _ := cmd.Flags().GetStringArray("bypath")
		param, _ := cmd.Flags().GetStringArray("param")

		if bypath != nil {
			getParameterByPath(bypath, profile, region)
		}
		if param != nil {
			getParameters(param, profile, region)
		}
	},
}

func getParameterByPath(params []string, profile string, region string) {
	ssmClient := pkg.NewSSM(profile, region)

	for k, _ := range params {
		results, err := ssmClient.SSM.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
			Path: &params[k],
		})
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}

		for _, n := range results.Parameters {
			envVar := strings.Split(*n.Name, "/")
			envVarLast := len(envVar)
			fmt.Println(envVar[envVarLast-1] + "=" + *n.Value)
		}
	}

}

func getParameters(params []string, profile string, region string) {
	ssmClient := pkg.NewSSM(profile, region)

	results, err := ssmClient.SSM.GetParameters(context.TODO(), &ssm.GetParametersInput{
		Names: params,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}

	for _, n := range results.Parameters {
		envVar := strings.Split(*n.Name, "/")
		envVarLast := len(envVar)
		fmt.Println(envVar[envVarLast-1] + "=" + *n.Value)
	}
}

func init() {
	getCmd.Flags().StringP("profile", "P", "default", "AWS configuration profile")
	getCmd.Flags().StringP("region", "R", "", "AWS configuration region")

	getCmd.Flags().StringArrayP("bypath", "b", nil, "Search query by path")
	getCmd.Flags().StringArrayP("param", "p", nil, "Search query by param")

	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
