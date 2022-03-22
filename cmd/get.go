/*
Copyright © 2022 Isaac Lopez syak7771@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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

// Struct json file
type componentSSM struct {
	PathSSM  string              `json:"path"`
	ParamSSM string              `json:"param"`
	ValueSSM string              `json:"value"`
	TypeSSM  types.ParameterType `json:"type"`
}

// Struct json file
type variablesSSM struct {
	VariablesSSM []componentSSM `json:"variables"`
}

// getParameters params
type flagsGet struct {
	names      []string // only needed for getParamters
	showPath   bool
	decryption bool
	json       string
}

// getParametersByPath params
type flagsGetByPath struct {
	flagsGet
	path     string
	variable string
	value    string
	contains bool
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

		path, _ := cmd.Flags().GetString("path")
		variable, _ := cmd.Flags().GetString("variable")
		value, _ := cmd.Flags().GetString("value")
		contains, _ := cmd.Flags().GetBool("contains")
		showPath, _ := cmd.Flags().GetBool("show-path")
		names, _ := cmd.Flags().GetStringArray("names")
		decryption, _ := cmd.Flags().GetBool("decryption")
		json, _ := cmd.Flags().GetString("json")

		flagsPath := flagsGetByPath{flagsGet{names, showPath, decryption, json}, path, variable, value, contains}
		flags := flagsGet{names, showPath, decryption, json}

		// Start indicator
		indicatorSpinner = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		indicatorSpinner.Start()
		indicatorSpinner.Prefix = "  "

		if flagsPath.value != "" || flagsPath.variable != "" || len(flagsPath.path) > 0 {
			if len(flagsPath.path) == 0 {
				flagsPath.path = "/"
			}
			getParametersByPath(flagsPath, profile, region, cmd)
			indicatorSpinner.Stop()
		}

		if len(flags.names) > 0 {
			getParameters(flags, profile, region, cmd)
			indicatorSpinner.Stop()
		}
		if len(flagsPath.path) == 0 && len(flags.names) == 0 {
			cmd.Help()
		}
	},
}

// getParamtersByPath retrive values from path without param.
func getParametersByPath(flag flagsGetByPath, profile string, region string, cmd *cobra.Command) {
	ssmClient := pkg.NewSSM(profile, region)

	results, err := ssmClient.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
		Path:           &flag.path,
		Recursive:      true,
		WithDecryption: flag.decryption,
	})
	if err != nil {
		dialog.Log("Error", err.Error(), cmd)
		os.Exit(1)
		return
	}

	for _, output := range results.Parameters {
		parametersOutput(flag.value, flag.variable, output, flag.contains, flag.showPath)
	}

	if results.NextToken != nil {
		getParametersByPathNextToken(flag, profile, region, results, cmd)
	} else if flag.json != "" {
		ssmP := ssmParam{SSMParamSlice, SSMValueSlice, SSMTypeSlice}
		writeJson(ssmP, flag.showPath, flag.json)
	}
}

// getParamtersByPathNexToken retrive values from path without param from the token.
func getParametersByPathNextToken(flag flagsGetByPath, profile string, region string, results *ssm.GetParametersByPathOutput, cmd *cobra.Command) {
	ssmClient := pkg.NewSSM(profile, region)

	nextToken := *results.NextToken

	results, err := ssmClient.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
		Path:           &flag.path,
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
		parametersOutput(flag.value, flag.variable, output, flag.contains, flag.showPath)
	}

	if results.NextToken != nil {
		nextPage(flag, profile, region, results)
	} else if flag.json != "" {
		ssmP := ssmParam{SSMParamSlice, SSMValueSlice, SSMTypeSlice}
		writeJson(ssmP, flag.showPath, flag.json)
	}

}

// nextPage paginator options for GetParametersByPath
func nextPage(flag flagsGetByPath, profile string, region string, results *ssm.GetParametersByPathOutput) {
	ssmClient := pkg.NewSSM(profile, region)

	nextToken := *results.NextToken

	paginator := ssm.NewGetParametersByPathPaginator(ssmClient, &ssm.GetParametersByPathInput{
		Path:           &flag.path,
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
			parametersOutput(flag.value, flag.variable, output, flag.contains, flag.showPath)
		}
	}

	if flag.json != "" {
		ssmP := ssmParam{SSMParamSlice, SSMValueSlice, SSMTypeSlice}
		writeJson(ssmP, flag.showPath, flag.json)
	}

}

// getParameters retrives values from path with param.
func getParameters(flag flagsGet, profile string, region string, cmd *cobra.Command) {
	ssmClient := pkg.NewSSM(profile, region)

	results, err := ssmClient.GetParameters(context.TODO(), &ssm.GetParametersInput{
		Names:          flag.names,
		WithDecryption: flag.decryption,
	})
	if err != nil {
		dialog.Log("Error", err.Error(), cmd)
		os.Exit(1)
		return
	}

	for _, output := range results.Parameters {
		parametersOutput("", "", output, false, flag.showPath)
	}

	if flag.json != "" {
		ssmP := ssmParam{SSMParamSlice, SSMValueSlice, SSMTypeSlice}
		writeJson(ssmP, flag.showPath, flag.json)
	}
}

// parametersOutput output with fullpath or without and search for value or param.
func parametersOutput(valueFlag string, variableFlag string, v types.Parameter, contains bool, showPathFlag bool) {
	envVar := strings.Split(*v.Name, "/")
	envVarLast := len(envVar)

	// Define suffix indicator
	indicatorSpinner.Suffix = "  " + *v.Name

	if showPathFlag == false {
		ouputWithWithoutFlag(valueFlag, variableFlag, v, contains, envVar[envVarLast-1])
	} else {
		ouputWithWithoutFlag(valueFlag, variableFlag, v, contains, *v.Name)
	}
}

func ouputWithWithoutFlag(valueFlag string, variableFlag string, v types.Parameter, contains bool, name string) {
	envVar := strings.Split(*v.Name, "/")
	envVarLast := len(envVar)

	SSMTypeSlice = append(SSMTypeSlice, v.Type)
	SSMValueSlice = append(SSMValueSlice, *v.Value)
	SSMParamSlice = append(SSMParamSlice, name)

	if valueFlag != "" {
		if contains {
			if strings.Contains(*v.Value, valueFlag) {
				outputColor(name, *v.Value)
			}
		} else {
			if valueFlag == *v.Value {
				outputColor(name, *v.Value)
			}
		}
	} else if variableFlag != "" {
		if contains {
			if strings.Contains(envVar[envVarLast-1], variableFlag) {
				outputColor(name, *v.Value)
			}
		} else {
			if variableFlag == envVar[envVarLast-1] {
				outputColor(name, *v.Value)
			}
		}
	} else {
		outputColor(name, *v.Value)
	}

}

func outputColor(name, value string) {
	indicatorSpinner.Stop()
	colorstring.Println("[blue]" + name + "=[reset]" + value)
	indicatorSpinner.Start()
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

		// checking if exists names in ssm without "/"
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
	getCmd.Flags().StringP("path", "p", "", "Search path recursively")
	getCmd.Flags().StringP("variable", "r", "", "Search variable in paths")
	getCmd.Flags().StringP("value", "v", "", "Search value in paths")
	getCmd.Flags().BoolP("contains", "c", false, "Search contains value in paths")
	getCmd.Flags().BoolP("show-path", "f", false, "return with path")
	getCmd.Flags().StringArrayP("names", "n", nil, "return specific names")
	getCmd.Flags().BoolP("decryption", "d", false, "Return decrypted secure string value")
	getCmd.Flags().StringP("json", "j", "", "Json path/name to write results")

	rootCmd.AddCommand(getCmd)
}
