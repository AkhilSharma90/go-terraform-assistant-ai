package cli

import (
	"flag"
	"log"
	"strconv"

	"github.com/akhilsharma90/terraform-assistant/pkg/terraform"
	"github.com/spf13/cobra"
	"github.com/walles/env"
)

const version = "0.0.2"

//gettig all values from environment variables and setting our variables
var (
	// openAIDeploymentName is the name of the deployment used for the model in the OpenAI service.
	openAIDeploymentName = flag.String("openai-deployment-name", env.GetOr("OPENAI_DEPLOYMENT_NAME", env.String, "text-davinci-003"), "The deployment name used for the model in OpenAI service.")

	// maxTokens is the maximum number of tokens that will be used. It overrides the max tokens in the max tokens map.
	maxTokens = flag.Int("max-tokens", env.GetOr("MAX_TOKENS", strconv.Atoi, 0), "The max token will overwrite the max tokens in the max tokens map.")

	// openAIAPIKey is the API key for the OpenAI service. This is required.
	openAIAPIKey = flag.String("openai-api-key", env.GetOr("OPENAI_API_KEY", env.String, ""), "The API key for the OpenAI service. This is required.")

	// azureOpenAIEndpoint is the endpoint for the Azure OpenAI service. If provided, Azure OpenAI service will be used instead of OpenAI service.
	azureOpenAIEndpoint = flag.String("azure-openai-endpoint", env.GetOr("AZURE_OPENAI_ENDPOINT", env.String, ""), "The endpoint for Azure OpenAI service. If provided, Azure OpenAI service will be used instead of OpenAI service.")

	// requireConfirmation specifies whether to require confirmation before executing the command. Defaults to true.
	requireConfirmation = flag.Bool("require-confirmation", env.GetOr("REQUIRE_CONFIRMATION", strconv.ParseBool, true), "Whether to require confirmation before executing the command. Defaults to true.")

	// temperature is the temperature to use for the model. Range is between 0 and 1. Set closer to 0 if you want output to be more deterministic but less creative. Defaults to 0.0.
	temperature = flag.Float64("temperature", env.GetOr("TEMPERATURE", env.WithBitSize(strconv.ParseFloat, 64), 0.0), "The temperature to use for the model. Range is between 0 and 1. Set closer to 0 if your want output to be more deterministic but less creative. Defaults to 0.0.")

	// workingDir is the path of the project that you want to run.
	workingDir = flag.String("working-dir", env.GetOr("WORKING_DIR", env.String, ""), "The path of project that you want to run.")

	// execDir is the path of Terraform.
	execDir = flag.String("exec-dir", env.GetOr("EXEC_DIR", env.String, ""), "The path of Terraform.")

	// ops is an instance of the terraform.Ops struct.
	ops terraform.Ops

	err error
)

// InitAndExecute initializes the working directory and execution directory,
// parses command line flags, and executes the root command.
func InitAndExecute(workDir string, executionDir string) {
	flag.Parse()

	//we have received workDir and executionDir in this function as args
	//we will check if the variables for workingDir and execdir we have defined above
	//that get values from the environment variables are empty, if yes
	//we set their values with what's received in the args
	// Set the working directory if not provided
	if *workingDir == "" {
		workingDir = &workDir
	}

	// Set the execution directory if not provided
	if *execDir == "" {
		execDir = &executionDir
	}

	// Check if the OpenAI API key is provided
	if *openAIAPIKey == "" {
		log.Fatal("Please provide an OpenAI key.")
	}

	// Execute the root command
	if err := RootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}

// RootCmd returns the root command for the CLI.
func RootCmd() *cobra.Command {
	//creates a new struct for Terraform (struct defined in the terraform.go file of terraform package)
	//the struct requires working directory and exec directory
	ops, err = terraform.NewTerraform(*workingDir, *execDir)
	if err != nil {
		return nil
	}
//use cobra to start and create the CLI to interact with the user
	cmd := &cobra.Command{
		Use:          "terraform-ai",
		Version:      version,
		Args:         cobra.MinimumNArgs(1),
		RunE:         runCommand, //essentially calling the runCommand which calls the run function (both in run.go file)
		SilenceUsage: true,
	}

	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	initCmd := addInit()
	cmd.AddCommand(initCmd)

	return cmd
}
