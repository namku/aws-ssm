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
		value, _ := cmd.Flags().GetString("value")
		fullPath, _ := cmd.Flags().GetBool("fullPath")
		param, _ := cmd.Flags().GetStringArray("param")

		if len(bypath) > 0 || value != "" {
			if value != "" {
				bypath = []string{"/"}
				getParametersByPath(bypath, profile, region, fullPath, value, cmd)
			} else {
				getParametersByPath(bypath, profile, region, fullPath, value, cmd)
			}
		}
		if len(param) > 0 {
			getParameters(param, profile, region, cmd)
		}
		if len(bypath) == 0 && len(param) == 0 {
			cmd.Help()
		}
	},
}

// getParamtersByPath retrive values from path without param.
func getParametersByPath(params []string, profile string, region string, fullPath bool, value string, cmd *cobra.Command) {
	ssmClient := pkg.NewSSM(profile, region)

	for k, _ := range params {
		results, err := ssmClient.SSM.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
			Path:      &params[k],
			Recursive: true,
		})
		if err != nil {
			dialog.Log("Error", err.Error(), cmd)
			os.Exit(1)
			return
		}

		for _, n := range results.Parameters {
			if fullPath == false {
				envVar := strings.Split(*n.Name, "/")
				envVarLast := len(envVar)
				if value != "" {
					if value == *n.Value {
						colorstring.Println("[blue]" + envVar[envVarLast-1] + "=[reset]" + *n.Value)
					}
				} else {
					colorstring.Println("[blue]" + envVar[envVarLast-1] + "=[reset]" + *n.Value)
				}
			} else {
				if value != "" {
					if value == *n.Value {
						colorstring.Println("[blue]" + *n.Name + "=[reset]" + *n.Value)
					}
				}
			}
		}

		if results.NextToken != nil {
			getParametersByPathNextToken(params, profile, region, fullPath, value, results, cmd)
		}
	}
}

// getParamtersByPathNexToken retrive values from path without param from the token.
func getParametersByPathNextToken(params []string, profile string, region string, fullPath bool, value string, results *ssm.GetParametersByPathOutput, cmd *cobra.Command) {
	ssmClient := pkg.NewSSM(profile, region)

	nextToken := *results.NextToken

	for k, _ := range params {
		results, err := ssmClient.SSM.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
			Path:      &params[k],
			Recursive: true,
			NextToken: &nextToken,
		})
		if err != nil {
			dialog.Log("Error", err.Error(), cmd)
			os.Exit(1)
			return
		}

		for _, n := range results.Parameters {
			if fullPath == false {
				envVar := strings.Split(*n.Name, "/")
				envVarLast := len(envVar)
				if value != "" {
					if value == *n.Value {
						colorstring.Println("[blue]" + envVar[envVarLast-1] + "=[reset]" + *n.Value)
					}
				} else {
					colorstring.Println("[blue]" + envVar[envVarLast-1] + "=[reset]" + *n.Value)
				}
			} else {
				if value != "" {
					if value == *n.Value {
						colorstring.Println("[blue]" + *n.Name + "=[reset]" + *n.Value)
					}
				}
			}
		}
	}
}

// getParameters retrives values from path with param.
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
	getCmd.Flags().StringP("value", "v", "", "Search value in all paths")
	getCmd.Flags().BoolP("fullPath", "f", false, "Output with full path param")
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
