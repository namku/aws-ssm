# `aws-ssm`

<b>Aws-ssm</> is a command line tool to manage fast and efficiently AWS parameter store.

These are some things that can be done with <b>aws-ssm</b> CLI:

- Change some hierarchies names.
- Search values in the hierarchies.
- Import a lot of variables to AWS System Manager Parameter Store from json file.
- Create a json file (hierarchies, values, type)
- ...

With <b>aws-ssm</b> you can do that and more cool  things fast and efficiently.

***

## ⚡️ Getting Started
Retrive the <b>aws-ssm</b> binary by downloading a pre-compiled binary from [`Download section`](https://github.com/namku/aws-ssm/tags) or compiling it from source (Go +1.7 required).

## ⚙️  Commands & Flags

### `get`

CLI command to get information from SSM.

```bash
aws-ssm get [flags]
```

| Flag | Description                                                                                                        | Type          | Default | Required |
|------|--------------------------------------------------------------------------------------------------------------------|---------------|---------|----------|
| `-n` | The complete name of the paramter (hierarchy).                                                                     | `stringArray` |         | No       |
| `-p` | The hierarchy for the parameter. Hierarchies start with a forward slash (/) except the last part of the parameter. | `string`      |         | No       |
| `-r` | The last part of the hierarchy (variable).                                                                         | `string`      |         | No       |
| `-v` | The value of the hierarchy.                                                                                        | `string`      |         | No       |
| `-f` | Print hierarchy.                                                                                                   | `bool`        | `false` | No       |
| `-d` | Print decrypted SecureString.                                                                                      | `bool`        | `false` | No       |
| `-c` | Search all values containing the value in -v flag.                                                                 | `bool`        | `false` | No       |
| `-j` | Write a json file with the output.                                                                                 | `string`      |         | No       |


### `add`

CLI command to add information from SSM.

```bash
aws-ssm add [flags]
```

| Flag | Description                                          | Type          | Default | Required |
|------|------------------------------------------------------|---------------|---------|----------|
| `-n` | Name of the hierarchy.                               | `string`      |         | No       |
| `-v` | Value of the hierarchy.                              | `string`      |         | No       |
| `-t` | Type of the value.                                   | `string`      |         | No       |
| `-o` | Overwrite the value of the hierarchy.                | `bool`        | `false` | No       |
| `-d` | Description of the hierarchy.                        | `string`      |         | No       |
| `-j` | Json file to import in the parameter store.          | `string`      |         | No       |


