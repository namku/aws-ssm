/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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

type flags struct {
	profile     string
	region      string
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
		description, _ := cmd.Flags().GetString("description")
		typeVar, _ := cmd.Flags().GetString("type")
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		json, _ := cmd.Flags().GetString("json")

		flag := flags{profile, region, name, value, description, typeVar, overwrite, json}

		if json != "" {
			importFromJson(flag.json, flag.profile, flag.region, flag.overwrite)
		} else {
			putParameter(flag)
		}
	},
}

func importFromJson(file string, profile string, region string, overwrite bool) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to read file, %v", err)
	}

	data := variablesSSM{}

	json.Unmarshal([]byte(content), &data)

	for i, _ := range data.VariablesSSM {
		putParameter(flags{profile: profile, region: region, name: data.VariablesSSM[i].PathSSM + data.VariablesSSM[i].ParamSSM, value: data.VariablesSSM[i].ValueSSM, description: "", typeVar: string(data.VariablesSSM[i].TypeSSM), overwrite: overwrite})
	}
}

func putParameter(flags flags) {
	ssmClient := pkg.NewSSM(flags.profile, flags.region)

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
	addCmd.Flags().StringP("name", "n", "", "Parameter name to add")
	addCmd.Flags().StringP("value", "v", "", "Value of the parameter")
	addCmd.Flags().StringP("description", "d", "", "Description of the parameter")
	addCmd.Flags().StringP("type", "t", "", "Type of the value [ string, stringList, secret ]")
	addCmd.Flags().BoolP("overwrite", "o", false, "Type of the value")
	addCmd.Flags().StringP("json", "j", "", "Json file name to Import")

	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
