/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/mitchellh/colorstring"
	"github.com/namku/aws-ssm/cmd/dialog"
	"github.com/namku/aws-ssm/pkg"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "The aws-ssm get command gets parameter values from SSM.",
	Long: `Usage: aws-ssm get [args]

The aws-ssm get command, gets parameter values from SSM. Different
flags can be agreed to obtain more extensive searches. For example:

aws-ssm get -p /vars/envs1/param1 -p /vars/envs2/param1 -b /vars/envs3

According to the search it can take a long time.`,
	Run: func(cmd *cobra.Command, args []string) {
		// flags for custom aws config
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")

		bypath, _ := cmd.Flags().GetStringArray("bypath")
		param, _ := cmd.Flags().GetStringArray("param")

		if len(bypath) > 0 {
			getParameterByPath(bypath, profile, region, cmd)
		}
		if len(param) > 0 {
			getParameters(param, profile, region, cmd)
		}
		if len(bypath) == 0 && len(param) == 0 {
			cmd.Help()
		}
	},
}

func getParameterByPath(params []string, profile string, region string, cmd *cobra.Command) {
	ssmClient := pkg.NewSSM(profile, region)

	for k, _ := range params {
		results, err := ssmClient.SSM.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
			Path: &params[k],
		})
		if err != nil {
			dialog.Log("Error", err.Error(), cmd)
			os.Exit(1)
			return
		}

		for _, n := range results.Parameters {
			envVar := strings.Split(*n.Name, "/")
			envVarLast := len(envVar)
			colorstring.Println("[blue]" + envVar[envVarLast-1] + "=[reset]" + *n.Value)
		}
	}

}

func getParameters(params []string, profile string, region string, cmd *cobra.Command) {
	ssmClient := pkg.NewSSM(profile, region)

	results, err := ssmClient.SSM.GetParameters(context.TODO(), &ssm.GetParametersInput{
		Names: params,
	})
	if err != nil {
		dialog.Log("Error", err.Error(), cmd)
		os.Exit(1)
		return
	}

	for _, n := range results.Parameters {
		envVar := strings.Split(*n.Name, "/")
		envVarLast := len(envVar)
		colorstring.Println("[blue]" + envVar[envVarLast-1] + "=[reset]" + *n.Value)
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
