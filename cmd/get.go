/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
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
		parameter, _ := cmd.Flags().GetString("parameter")
		value, _ := cmd.Flags().GetString("value")
		fullPath, _ := cmd.Flags().GetBool("fullPath")
		param, _ := cmd.Flags().GetStringArray("param")

		if len(bypath) > 0 || value != "" || parameter != "" {
			if value != "" || parameter != "" {
				bypath = []string{"/"}
			} else {
				bypath = bypath
			}
			getParametersByPath(bypath, profile, region, fullPath, parameter, value, cmd)
		}
		if len(param) > 0 {
			getParameters(param, profile, region, fullPath, cmd)
		}
		if len(bypath) == 0 && len(param) == 0 {
			cmd.Help()
		}
	},
}

// getParamtersByPath retrive values from path without param.
func getParametersByPath(params []string, profile string, region string, fullPath bool, parameter string, value string, cmd *cobra.Command) {
	ssmClient := pkg.NewSSM(profile, region)

	for k, _ := range params {
		results, err := ssmClient.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
			Path:      &params[k],
			Recursive: true,
		})
		if err != nil {
			dialog.Log("Error", err.Error(), cmd)
			os.Exit(1)
			return
		}

		for _, n := range results.Parameters {
			parametersOutput(value, parameter, n, fullPath)
		}

		if results.NextToken != nil {
			getParametersByPathNextToken(params, profile, region, fullPath, parameter, value, results, cmd)
		}
	}
}

// getParamtersByPathNexToken retrive values from path without param from the token.
func getParametersByPathNextToken(params []string, profile string, region string, fullPath bool, parameter string, value string, results *ssm.GetParametersByPathOutput, cmd *cobra.Command) {
	ssmClient := pkg.NewSSM(profile, region)

	nextToken := *results.NextToken

	for k, _ := range params {
		results, err := ssmClient.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
			Path:      &params[k],
			Recursive: true,
			NextToken: &nextToken,
		})
		if err != nil {
			dialog.Log("Error", err.Error(), cmd)
			os.Exit(1)
			return
		}

		for _, v := range results.Parameters {
			parametersOutput(value, parameter, v, fullPath)
		}

		if results.NextToken != nil {
			nextPage(params, profile, region, fullPath, parameter, value, results)
		}
	}
}

// nextPage paginator options for GetParametersByPath
func nextPage(params []string, profile string, region string, fullPath bool, parameter string, value string, results *ssm.GetParametersByPathOutput) {
	ssmClient := pkg.NewSSM(profile, region)

	nextToken := *results.NextToken

	for k, _ := range params {
		paginator := ssm.NewGetParametersByPathPaginator(ssmClient, &ssm.GetParametersByPathInput{
			Path:      &params[k],
			Recursive: true,
			NextToken: &nextToken,
		})

		for paginator.HasMorePages() {
			results, err := paginator.NextPage(context.TODO())
			if err != nil {
				log.Fatalf("failed to get page , %v", err)
			}

			for _, v := range results.Parameters {
				parametersOutput(value, parameter, v, fullPath)
			}
		}
	}
}

// parametersOutput output with fullpath or without and search for value or param.
func parametersOutput(value string, parameter string, v types.Parameter, fullPath bool) {
	envVar := strings.Split(*v.Name, "/")
	envVarLast := len(envVar)
	if fullPath == false {
		if value != "" {
			if value == *v.Value {
				colorstring.Println("[blue]" + envVar[envVarLast-1] + "=[reset]" + *v.Value)
			}
		} else if parameter != "" {
			if parameter == envVar[envVarLast-1] {
				colorstring.Println("[blue]" + envVar[envVarLast-1] + "=[reset]" + *v.Value)
			}
		} else {
			colorstring.Println("[blue]" + envVar[envVarLast-1] + "=[reset]" + *v.Value)
		}
	} else {
		if value != "" {
			if value == *v.Value {
				colorstring.Println("[blue]" + *v.Name + "=[reset]" + *v.Value)
			}
		} else if parameter != "" {
			if parameter == envVar[envVarLast-1] {
				colorstring.Println("[blue]" + *v.Name + "=[reset]" + *v.Value)
			}
		} else {
			colorstring.Println("[blue]" + *v.Name + "=[reset]" + *v.Value)
		}
	}

}

// getParameters retrives values from path with param.
func getParameters(params []string, profile string, region string, fullPath bool, cmd *cobra.Command) {
	ssmClient := pkg.NewSSM(profile, region)

	results, err := ssmClient.GetParameters(context.TODO(), &ssm.GetParametersInput{
		Names: params,
	})
	if err != nil {
		dialog.Log("Error", err.Error(), cmd)
		os.Exit(1)
		return
	}

	for _, v := range results.Parameters {
		parametersOutput("", "", v, fullPath)
	}
}

func init() {
	//getCmd.Flags().StringP("profile", "P", "default", "AWS configuration profile")
	//getCmd.Flags().StringP("region", "R", "", "AWS configuration region")

	getCmd.Flags().StringArrayP("bypath", "b", nil, "Search query by path")
	getCmd.Flags().StringP("parameter", "r", "", "Search parameter in all paths")
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
