/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"

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

		f := flags{profile, region, name, value, description, typeVar, overwrite}

		putParameter(f)
	},
}

func putParameter(flags flags) {
	ssmClient := pkg.NewSSM(flags.profile, flags.region)

	var typeValue types.ParameterType

	// Improve, string to types.ParameterType
	switch flags.typeVar {
	case "string":
		typeValue = "String"
	case "stringList":
		typeValue = "StringList"
	case "secret":
		typeValue = "SecureString"
	default:
		fmt.Println("Valid options for --type [ string, stringList, secret ]")
	}

	_, err := ssmClient.PutParameter(context.TODO(), &ssm.PutParameterInput{
		Name:        &flags.name,
		Value:       &flags.value,
		Description: &flags.description,
		Type:        typeValue,
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
	addCmd.Flags().StringP("type", "t", "", "Type of the value")
	addCmd.Flags().BoolP("overwrite", "o", false, "Type of the value")

	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
