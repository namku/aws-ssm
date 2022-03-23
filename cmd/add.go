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
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/namku/aws-ssm/pkg"
	"github.com/spf13/cobra"
)

type flagsPut struct {
	name        string
	value       string
	description string
	typeVar     string
	overwrite   bool
	json        string
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
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

		name, _ := cmd.Flags().GetString("name")
		value, _ := cmd.Flags().GetString("value")
		typeVar, _ := cmd.Flags().GetString("type")
		description, _ := cmd.Flags().GetString("description")
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		json, _ := cmd.Flags().GetString("json")

		flag := flagsPut{name, value, description, typeVar, overwrite, json}

		if json != "" {
			startSpinner()
			importFromJson(flag.json, flag.overwrite, profile, region)
		} else {
			putParameter(flag, profile, region)
		}
	},
}

func importFromJson(file string, overwrite bool, profile string, region string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to read file, %v", err)
	}

	data := variablesSSM{}

	json.Unmarshal([]byte(content), &data)

	for i, _ := range data.VariablesSSM {
		// Define suffix spinner
		indicatorSpinner.Suffix = "  " + data.VariablesSSM[i].PathSSM + data.VariablesSSM[i].ValueSSM

		putParameter(flagsPut{name: data.VariablesSSM[i].PathSSM + data.VariablesSSM[i].ParamSSM, value: data.VariablesSSM[i].ValueSSM, description: "", typeVar: string(data.VariablesSSM[i].TypeSSM), overwrite: overwrite}, profile, region)
	}

	// Stop spinner
	indicatorSpinner.Stop()
}

func putParameter(flags flagsPut, profile string, region string) {
	ssmClient := pkg.NewSSM(profile, region)

	_, err := ssmClient.PutParameter(context.TODO(), &ssm.PutParameterInput{
		Name:        &flags.name,
		Value:       &flags.value,
		Description: &flags.description,
		Type:        types.ParameterType(flags.typeVar),
		Overwrite:   flags.overwrite,
	})

	if err != nil {
		fmt.Println(err)
	}
}

func init() {
	addCmd.Flags().StringP("name", "n", "", "Name of the hierarchy.")
	addCmd.Flags().StringP("value", "v", "", "Value of the hierarchy.")
	addCmd.Flags().StringP("type", "t", "", "Type of the value. [ string, stringList, secret ]")
	addCmd.Flags().StringP("description", "d", "", "Description of the hierarchy.")
	addCmd.Flags().BoolP("overwrite", "o", false, "Overwrite the value of the hierarchy.")
	addCmd.Flags().StringP("json", "j", "", "Json file to import in the parameter store.")

	rootCmd.AddCommand(addCmd)
}
