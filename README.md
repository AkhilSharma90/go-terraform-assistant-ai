# Terraform OpenAI 

`Terraform-assistant`, powered by OpenAI ChatGPT, simplifies the process of applying Terraform HCL files by providing an intelligent and interactive assistant.

## Usage

### Prerequisites

`terraform-assistant` requires an [OpenAI API key](https://platform.openai.com/overview) or an [Azure OpenAI Service](https://aka.ms/azure-openai) API key and endpoint.

For both OpenAI and Azure OpenAI, you can use the following environment variables:

```shell
export OPENAI_API_KEY=<your OpenAI key>
export OPENAI_DEPLOYMENT_NAME=<your OpenAI deployment/model name. defaults to "gpt-3.5-turbo">
```

> Following models are supported:
> - `code-davinci-002`
> - `text-davinci-003`
> - `gpt-3.5-turbo-0301` (deployment must be named `gpt-35-turbo-0301` for Azure)
> - `gpt-3.5-turbo`
> - `gpt-35-turbo-0301`
> - `gpt-4-0314`
> - `gpt-4-32k-0314`

For Azure OpenAI Service, you can use the following environment variables:

```shell
export AZURE_OPENAI_ENDPOINT=<your Azure OpenAI endpoint, like "https://my-aoi-endpoint.openai.azure.com">
```

If `AZURE_OPENAI_ENDPOINT` variable is set, then it will use the Azure OpenAI Service. Otherwise, it will use OpenAI API.

### Flags and environment variables

- `--require-confirmation` flag or `REQUIRE_CONFIRMATION` environment varible can be set to prompt the user for confirmation before applying the manifest. Defaults to true.

- `--temperature` flag or `TEMPERATURE` environment variable can be set between 0 and 1. Higher temperature will result in more creative completions. Lower temperature will result in more deterministic completions. Defaults to 0.

- `--working-dir` flag or `WORKING_DIR` environment variable that can be set for the Terraform project path.

- `--exec-dir` flag or `EXEC_DIR` environment variable that can be set for the Terraform executable binary file.

## Examples

### Creating templates
```shell
go run main.go "create micro ec2  ubuntu image 20.04 with name hello-future"

🦄 Attempting to store the following template:
resource "aws_instance" "hello_future" {
  ami           = "ami-0f65671a86f061fcd"
  instance_type = "t2.micro"
  tags = {
    Name = "hello-future"
  }
}
Use the arrow keys to navigate: ↓ ↑ → ←
? Would you like to apply this? [Reprompt/Apply/Don't Apply]:
+   Reprompt
  ▸ Apply
    Don't Apply
```

### Init provider

```shell
go run main.go init "create aws provider in ohio"

🦄 Attempting to apply the following template:

provider "aws" {
  region  = "us-east-2"
  alias   = "Ohio"
}
Use the arrow keys to navigate: ↓ ↑ → ←
? Would you like to apply this? [Reprompt/Apply/Don't Apply]:
+   Reprompt
  ▸ Apply
    Don't Apply
```


### Optional `--require-confirmation` flag
