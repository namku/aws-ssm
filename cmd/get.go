/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
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

type flagsGet struct {
	profile  string
	region   string
	param    []string // only needed for getParamters
	fullPath bool
}

type flagsGetByPath struct {
	flagsGet
	bypath    []string
	parameter string
	value     string
}

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

		flagsPath := flagsGetByPath{flagsGet{profile, region, param, fullPath}, bypath, parameter, value}
		flags := flagsGet{profile, region, param, fullPath}

		if len(flagsPath.bypath) > 0 || flagsPath.value != "" || flagsPath.parameter != "" {
			if flagsPath.value != "" || flagsPath.parameter != "" {
				flagsPath.bypath = []string{"/"}
			} else {
				flagsPath.bypath = flagsPath.bypath
			}

			getParametersByPath(flagsPath, cmd)
		}
		if len(flags.param) > 0 {
			getParameters(flags, cmd)
		}
		if len(flagsPath.bypath) == 0 && len(flags.param) == 0 {
			cmd.Help()
		}
	},
}

// getParamtersByPath retrive values from path without param.
func getParametersByPath(flag flagsGetByPath, cmd *cobra.Command) {
	ssmClient := pkg.NewSSM(flag.profile, flag.region)

	//var res *string
	for k, _ := range flag.bypath {
		results, err := ssmClient.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
			Path:      &flag.bypath[k],
			Recursive: true,
		})
		if err != nil {
			dialog.Log("Error", err.Error(), cmd)
			os.Exit(1)
			return
		}

		for _, output := range results.Parameters {
			parametersOutput(flag.value, flag.parameter, output, flag.fullPath)
		}

		//if res == nil {
		//	res = results.NextToken
		//	//	fmt.Println("adiooooooooooooooooooooooooooooooooooooooooos")
		//} else {
		//	//	fmt.Println("isaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaac")
		//	results.NextToken = res
		//	//	fmt.Println(results.NextToken)
		//}
		//fmt.Println(*results.NextToken)
		//if flag.bypath[k] == "/directj" {
		//if (results.NextToken != nil || k < len(flag.bypath)-1) || (results.NextToken != nil && len(flag.bypath) == 1) {
		//if results.NextToken != nil {
		//	getParametersByPathNextToken(flag, results, cmd)
		//}
	}
}

// getParamtersByPathNexToken retrive values from path without param from the token.
func getParametersByPathNextToken(flag flagsGetByPath, results *ssm.GetParametersByPathOutput, cmd *cobra.Command) {
	ssmClient := pkg.NewSSM(flag.profile, flag.region)

	nextToken := *results.NextToken

	fmt.Println(results.NextToken)
	results, err := ssmClient.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
		Path:      &flag.bypath[0],
		Recursive: true,
		NextToken: &nextToken,
	})
	//fmt.Println(results.NextToken)
	fmt.Println("Isaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaac")
	if err != nil {
		//	fmt.Println(results.NextToken)
		dialog.Log("Error", err.Error(), cmd)
		os.Exit(1)
		return
	}

	for _, output := range results.Parameters {
		parametersOutput(flag.value, flag.parameter, output, flag.fullPath)
	}

	//if results.NextToken != nil {
	//	nextPage(flag, results)
	//}
}

// nextPage paginator options for GetParametersByPath
func nextPage(flag flagsGetByPath, results *ssm.GetParametersByPathOutput) {
	ssmClient := pkg.NewSSM(flag.profile, flag.region)

	nextToken := *results.NextToken

	paginator := ssm.NewGetParametersByPathPaginator(ssmClient, &ssm.GetParametersByPathInput{
		Path:      &flag.bypath[0],
		Recursive: true,
		NextToken: &nextToken,
	})

	for paginator.HasMorePages() {
		results, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page , %v", err)
		}

		for _, output := range results.Parameters {
			parametersOutput(flag.value, flag.parameter, output, flag.fullPath)
		}
	}
	fmt.Println("isaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaac")
	fmt.Println(*results.NextToken)
}

// getParameters retrives values from path with param.
func getParameters(flags flagsGet, cmd *cobra.Command) {
	ssmClient := pkg.NewSSM(flags.profile, flags.region)

	results, err := ssmClient.GetParameters(context.TODO(), &ssm.GetParametersInput{
		Names: flags.param,
	})
	if err != nil {
		dialog.Log("Error", err.Error(), cmd)
		os.Exit(1)
		return
	}

	for _, output := range results.Parameters {
		parametersOutput("", "", output, flags.fullPath)
	}
}

// parametersOutput output with fullpath or without and search for value or param.
func parametersOutput(valueFlag string, parameterFlag string, v types.Parameter, fullPathFlag bool) {
	envVar := strings.Split(*v.Name, "/")
	envVarLast := len(envVar)

	if fullPathFlag == false {
		if valueFlag != "" {
			if valueFlag == *v.Value {
				colorstring.Println("[blue]" + envVar[envVarLast-1] + "=[reset]" + *v.Value)
			}
		} else if parameterFlag != "" {
			if parameterFlag == envVar[envVarLast-1] {
				colorstring.Println("[blue]" + envVar[envVarLast-1] + "=[reset]" + *v.Value)
			}
		} else {
			colorstring.Println("[blue]" + envVar[envVarLast-1] + "=[reset]" + *v.Value)
		}
	} else {
		if valueFlag != "" {
			if valueFlag == *v.Value {
				colorstring.Println("[blue]" + *v.Name + "=[reset]" + *v.Value)
			}
		} else if parameterFlag != "" {
			if parameterFlag == envVar[envVarLast-1] {
				colorstring.Println("[blue]" + *v.Name + "=[reset]" + *v.Value)
			}
		} else {
			colorstring.Println("[blue]" + *v.Name + "=[reset]" + *v.Value)
		}
	}

}

func init() {
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
