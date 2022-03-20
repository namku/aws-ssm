/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/briandowns/spinner"
	"github.com/mitchellh/colorstring"
	"github.com/namku/aws-ssm/cmd/dialog"
	"github.com/namku/aws-ssm/pkg"
	"github.com/spf13/cobra"
)

// Struct for json file
type componentSSM struct {
	PathSSM  string
	ParamSSM string
	ValueSSM string
	TypeSSM  types.ParameterType
}

// Struct for json file
type variablesSSM struct {
	VariablesSSM []componentSSM
}

// getParameters params
type flagsGet struct {
	profile    string
	region     string
	param      []string // only needed for getParamters
	fullPath   bool
	decryption bool
	json       string
}

// getParametersByPath params
type flagsGetByPath struct {
	flagsGet
	bypath    string
	parameter string
	value     string
}

type ssmParam struct {
	ssmParam []string
	ssmValue []string
	ssmType  []types.ParameterType
}

var SSMParamSlice []string
var SSMValueSlice []string
var SSMTypeSlice []types.ParameterType

// Indicator channel
var indicatorSpinner *spinner.Spinner

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

		bypath, _ := cmd.Flags().GetString("bypath")
		parameter, _ := cmd.Flags().GetString("parameter")
		value, _ := cmd.Flags().GetString("value")
		fullPath, _ := cmd.Flags().GetBool("fullPath")
		param, _ := cmd.Flags().GetStringArray("param")
		decryption, _ := cmd.Flags().GetBool("decryption")
		json, _ := cmd.Flags().GetString("json")

		flagsPath := flagsGetByPath{flagsGet{profile, region, param, fullPath, decryption, json}, bypath, parameter, value}
		flags := flagsGet{profile, region, param, fullPath, decryption, json}

		// Start indicator
		indicatorSpinner = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		indicatorSpinner.Start()
		indicatorSpinner.Prefix = "  "

		if flagsPath.value != "" || flagsPath.parameter != "" || len(flagsPath.bypath) > 0 {
			if len(flagsPath.bypath) == 0 {
				flagsPath.bypath = "/"
			}
			getParametersByPath(flagsPath, cmd)
			indicatorSpinner.Stop()
		}

		if len(flags.param) > 0 {
			getParameters(flags, cmd)
			indicatorSpinner.Stop()
		}
		if len(flagsPath.bypath) == 0 && len(flags.param) == 0 {
			cmd.Help()
		}
	},
}

// getParamtersByPath retrive values from path without param.
func getParametersByPath(flag flagsGetByPath, cmd *cobra.Command) {
	ssmClient := pkg.NewSSM(flag.profile, flag.region)

	results, err := ssmClient.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
		Path:           &flag.bypath,
		Recursive:      true,
		WithDecryption: flag.decryption,
	})
	if err != nil {
		dialog.Log("Error", err.Error(), cmd)
		os.Exit(1)
		return
	}

	for _, output := range results.Parameters {
		parametersOutput(flag.value, flag.parameter, output, flag.fullPath)
	}

	if results.NextToken != nil {
		getParametersByPathNextToken(flag, results, cmd)
	} else if flag.json != "" {
		ssmP := ssmParam{SSMParamSlice, SSMValueSlice, SSMTypeSlice}
		writeJson(ssmP, flag.fullPath, flag.json)
	}
}

// getParamtersByPathNexToken retrive values from path without param from the token.
func getParametersByPathNextToken(flag flagsGetByPath, results *ssm.GetParametersByPathOutput, cmd *cobra.Command) {
	ssmClient := pkg.NewSSM(flag.profile, flag.region)

	nextToken := *results.NextToken

	results, err := ssmClient.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
		Path:           &flag.bypath,
		Recursive:      true,
		NextToken:      &nextToken,
		WithDecryption: flag.decryption,
	})
	if err != nil {
		dialog.Log("Error", err.Error(), cmd)
		os.Exit(1)
		return
	}

	for _, output := range results.Parameters {
		parametersOutput(flag.value, flag.parameter, output, flag.fullPath)
	}

	if results.NextToken != nil {
		nextPage(flag, results)
	} else if flag.json != "" {
		ssmP := ssmParam{SSMParamSlice, SSMValueSlice, SSMTypeSlice}
		writeJson(ssmP, flag.fullPath, flag.json)
	}

}

// nextPage paginator options for GetParametersByPath
func nextPage(flag flagsGetByPath, results *ssm.GetParametersByPathOutput) {
	ssmClient := pkg.NewSSM(flag.profile, flag.region)

	nextToken := *results.NextToken

	paginator := ssm.NewGetParametersByPathPaginator(ssmClient, &ssm.GetParametersByPathInput{
		Path:           &flag.bypath,
		Recursive:      true,
		NextToken:      &nextToken,
		WithDecryption: flag.decryption,
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

	if flag.json != "" {
		ssmP := ssmParam{SSMParamSlice, SSMValueSlice, SSMTypeSlice}
		writeJson(ssmP, flag.fullPath, flag.json)
	}

}

// getParameters retrives values from path with param.
func getParameters(flag flagsGet, cmd *cobra.Command) {
	ssmClient := pkg.NewSSM(flag.profile, flag.region)

	results, err := ssmClient.GetParameters(context.TODO(), &ssm.GetParametersInput{
		Names:          flag.param,
		WithDecryption: flag.decryption,
	})
	if err != nil {
		dialog.Log("Error", err.Error(), cmd)
		os.Exit(1)
		return
	}

	for _, output := range results.Parameters {
		parametersOutput("", "", output, flag.fullPath)
	}

	if flag.json != "" {
		ssmP := ssmParam{SSMParamSlice, SSMValueSlice, SSMTypeSlice}
		writeJson(ssmP, flag.fullPath, flag.json)
	}
}

// parametersOutput output with fullpath or without and search for value or param.
func parametersOutput(valueFlag string, parameterFlag string, v types.Parameter, fullPathFlag bool) {
	envVar := strings.Split(*v.Name, "/")
	envVarLast := len(envVar)

	SSMTypeSlice = append(SSMTypeSlice, v.Type)
	SSMValueSlice = append(SSMValueSlice, *v.Value)

	// Define prefix and suffix indicator:
	indicatorSpinner.Prefix = "  "
	indicatorSpinner.Suffix = "  " + *v.Name

	if fullPathFlag == false {
		SSMParamSlice = append(SSMParamSlice, envVar[envVarLast-1])

		if valueFlag != "" {
			if valueFlag == *v.Value {
				indicatorSpinner.Stop()
				colorstring.Println("[blue]" + envVar[envVarLast-1] + "=[reset]" + *v.Value)
				indicatorSpinner.Start()
			}
		} else if parameterFlag != "" {
			if parameterFlag == envVar[envVarLast-1] {
				indicatorSpinner.Stop()
				colorstring.Println("[blue]" + envVar[envVarLast-1] + "=[reset]" + *v.Value)
				indicatorSpinner.Start()
			}
		} else {
			indicatorSpinner.Stop()
			colorstring.Println("[blue]" + envVar[envVarLast-1] + "=[reset]" + *v.Value)
			indicatorSpinner.Start()
		}
	} else {
		SSMParamSlice = append(SSMParamSlice, *v.Name)
		if valueFlag != "" {
			if valueFlag == *v.Value {
				indicatorSpinner.Stop()
				colorstring.Println("[blue]" + *v.Name + "=[reset]" + *v.Value)
				indicatorSpinner.Start()
			}
		} else if parameterFlag != "" {
			if parameterFlag == envVar[envVarLast-1] {
				indicatorSpinner.Stop()
				colorstring.Println("[blue]" + *v.Name + "=[reset]" + *v.Value)
				indicatorSpinner.Start()
			}
		} else {
			indicatorSpinner.Stop()
			colorstring.Println("[blue]" + *v.Name + "=[reset]" + *v.Value)
			indicatorSpinner.Start()
		}
	}

}

func writeJson(ssmParam ssmParam, flagFullPath bool, jsonFile string) {
	var jsonData variablesSSM
	var componentsSSM []componentSSM

	pathRegex, err := regexp.Compile(`/(.*)\/`)
	if err != nil {
		log.Fatal(err)
	}

	for k, _ := range ssmParam.ssmValue {
		sliceFullPath := strings.Split(ssmParam.ssmParam[k], "/")
		paramPos := len(sliceFullPath)
		param := sliceFullPath[paramPos-1]
		path := pathRegex.FindStringSubmatch(ssmParam.ssmParam[k])

		// checking if exists parameters in ssm without "/"
		if len(path) == 0 {
			path = append(path, ssmParam.ssmParam[k])
		}

		componentsSSM = append(componentsSSM, componentSSM{PathSSM: path[0], ParamSSM: param, ValueSSM: ssmParam.ssmValue[k], TypeSSM: ssmParam.ssmType[k]})
	}

	jsonData = variablesSSM{componentsSSM}

	content, err := json.MarshalIndent(jsonData, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(jsonFile, content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	getCmd.Flags().StringP("bypath", "b", "", "Search query by path")
	getCmd.Flags().StringP("parameter", "r", "", "Search parameter in all paths")
	getCmd.Flags().StringP("value", "v", "", "Search value in all paths")
	getCmd.Flags().BoolP("fullPath", "f", false, "Output with full path param")
	getCmd.Flags().StringArrayP("param", "p", nil, "Search query by param")
	getCmd.Flags().BoolP("decryption", "d", false, "Return decrypted secure string value")
	getCmd.Flags().StringP("json", "j", "", "Json file name to write query output")

	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
