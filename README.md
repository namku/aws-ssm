# `aws-ssm`

Aws-ssm is a command line tool to manage fast and efficiently AWS parameter store.

How many times has the name of a project changed and it referred to the parameters saved in SSM or have you had to look for some values and then change them in all those parameters, with aws-ssm you can do that and more cool  things fast and efficiently.

***

## ⚡️ Getting Started
Retrive the <b>aws-ssm</b> binary by downloading a pre-compiled binary from [`Download section`](https://github.com/namku/aws-ssm/tags) or compiling it from source (Go +1.7 required).

## ⚙️  Commands & Flags

### `get`

CLI command to get information from SSM.

```bash
aws-ssm get [flags]
```

| Flag | Description                                          | Type          | Default | Required |
|------|------------------------------------------------------|---------------|---------|----------|
| `-n` | Return the value of the searched name.               | `stringArray` |         | No       |
| `-p` | Return the value(s) of the path searched recursivly. | `string`      |         | No       |
| `-r` | Return the value(s) of the variable searched.        | `string`      |         | No       |
| `-v` | Return the path(s) of the value searched.            | `string`      |         | No       |
| `-f` | Output with full name not only variable.             | `bool`        | `false` | No       |
| `-d` | Decrypt SecureString output.                         | `bool`        | `false` | No       |
| `-c` | Return the path(s) contain this value.               | `bool`        | `false` | No       |
| `-j` | Write a json file with the output.                   | `string`      |         | No       |


### `add`

CLI command to add information from SSM.

```bash
aws-ssm add [flags]
```

| Flag | Description                                          | Type          | Default | Required |
|------|------------------------------------------------------|---------------|---------|----------|
| `-n` | Name of the parameter.                               | `string`      |         | No       |
| `-v` | Value of the parameter.                              | `string`      |         | No       |
| `-t` | Type of the value.                                   | `string`      |         | No       |
| `-o` | Overwrite the value of the parameter.                | `bool`        | `false` | No       |
| `-d` | Description of the parameter.                        | `string`      |         | No       |
| `-j` | Json file to import in the parameter store.          | `string`      |         | No       |
