# `aws-ssm`

Aws-ssm is a command line tool to manage fast and efficiently AWS parameter store.

How many times has the name of a project changed and it referred to the parameters saved in SSM or have you had to look for some values and then change them in all those parameters, with aws-ssm you can do that and more cool  things fast and efficiently.

***

## ⚡️ Getting Started
Retrive the <b>aws-ssm</b> binary by downloading a pre-compiled binary from [`Download section`](https://github.com/namku/aws-ssm/tags) or compiling it from source (Go +1.7 required).

## ⚙️  Commands & Flags

### `get`

CLI command for get information from SSM.

```bash
aws-ssm get [flags]
```

| Flag | Description                                 | Type          | Default | Required |
|------|---------------------------------------------|---------------|---------|----------|
| `-n` | Return the value from of the searched name. | `stringArray` |   ""    |    No    |
