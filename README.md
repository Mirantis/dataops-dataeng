# Data Engineering - dataeng

![Image](https://github.com/Mirantis/dataeng/blob/main/.assets/logos/dataeng_logo_1.png)

## Status Checks

[![Greetings](https://github.com/Mirantis/dataeng/actions/workflows/greetings.yml/badge.svg)](https://github.com/Mirantis/dataeng/actions/workflows/greetings.yml)
[![Labeler](https://github.com/Mirantis/dataeng/actions/workflows/label.yml/badge.svg)](https://github.com/Mirantis/dataeng/actions/workflows/label.yml)

The dataeng project is dedicated to storing all data interactions across an organization as IaC.
It provides a Command Line Interface to interact witht the various software tools that are used on a daily basis.
It also provides an anlaysis perspective, and offers a wide variety of tooling in terms of analyzing overall operational status of the entire organization for all products.

## Repository Structure

This repository is loosely organized into 11 main categories:

<!-- TODO: REFACTOR -->
1. [.assets](https://github.com/Mirantis/dataeng/tree/main/.assets)
    - An assets directory housing images and assets used throughout the directory (EACH FILE HAS AN MD5SUM TO VERIFY FILE INTEGRITY)

2. [.docker](https://github.com/Mirantis/dataeng/tree/main/.docker)
    - A configuration and secrets directory for all things docker

4. [.github](https://github.com/Mirantis/dataeng/tree/main/.github/)
    - A configuration and secrets directory for all things github

5. [.kube](https://github.com/Mirantis/dataeng/tree/main/.kube)
    - A configuration and secrets directory for all things kubernetes

6. [.secrets](https://github.com/Mirantis/dataeng/tree/main/.secrets)
    - A configuration and secrets directory for all things secret

7. [Analyses](https://github.com/Mirantis/dataeng/tree/main/Analyses)
    - The Analyses directory that houses a variety of analyses across multiple toolings

8. [Applications](https://github.com/Mirantis/dataeng/tree/main/Applications)
    - The Applications directory houses any and all affiliated applications to the housing repository

9. [Automation](https://github.com/Mirantis/dataeng/tree/main/Automation)
    - The Automation directory will contain automation and tooling for deploying and utilizing the repository

10. [Charts](https://github.com/Mirantis/dataeng/tree/main/Charts)
    - The Charts directory will container Helm Charts related to this repository

11. [Infrastructure](https://github.com/Mirantis/dataeng/tree/main/Infrastructure)
    - The Infrastrucutre directory contains information relate to Ansible and Terraform for continuously deployed infrastructure

12. [Lakes](https://github.com/Mirantis/dataeng/tree/main/Lakes)
    - The Lakes Directory will house the information related to any and all data lakes within the organization

13. [Pipelines](https://github.com/Mirantis/dataeng/tree/main/Pipelines)
    - This will contain any pipelines related to the overall repisitory and is tooling dependent

14. [Queries](https://github.com/Mirantis/dataeng/tree/main/Queries)
    - The Queries directory contains SQL Like queries that are stored for IaC purposes

15. [Warehouse](https://github.com/Mirantis/dataeng/tree/main/Warehouse)
    - The Warehouse directory contains information on how to interact with the Data Warehouse for the respective organization

16. [dataengctl](https://github.com/Mirantis/dataeng/tree/main)
    - The dataengctl houses the binary that allows you to work with a variety of tooling at the organization's disposal

## Navigating Repository

The above showcases the directories that are contained within the repository.
All directories labeled with `.` in front, such as `.secrets` are meant to be ignored, and nothing should ever be commited inside of these directories. The main directories are each provided with a `README.md` that tells you how you should interact with the repository in that particular directory.

## Getting Started with dataenctl

To get started with `dataengctl`, you should create a configuration file wihtin your home directory.
It should be named `~/.dataeng` and should be a hidden directory. In this directory you should create a file called `config.yaml`, which contains the secrets and credentials needed to work with `dataengctl`. The config should look something like this:

```bash
jiraConfig:
  token: YOURJIRATOKEN
  url: YOURJIRAURL
  username: YOURUSERNAME
salesForceConfig:
  url: YOURSALESFORCEURL
  username: YOURSALESFORCEUSERNAME
  password: YOURSALESFORCEPASSWORD
  token: YOURSALESFORCETOKEN
  clientID: YOURSALESFORCEORGCLIENTID
  apiVersion: YOURSALESFORCEDEFUALTAPI
```

## Installing dataengctl

To install dataengctl you need to make sure you have two things installed

1. You Want `jq` installed on your local host
2. You want `go` installed on your local host

## Working with dataengctl

To work with dataengctl you can use the `makefile` that is housed within the dataengctl directory, or you can build it directly in that directory. All Binary files are ignored by default within that directory.

1. Build the dataengctl binary
    - `go build -o dataengctl`
2. Afterwards run the binary specifying it with `--help`

```bash
 rbarrett@MacBook-Pro-2  ~/Git/dataeng   DATAENG-19 ●  . --help                                                      1 ↵  10130  11:38:37
gets you data from different sources

Usage:
  dataengctl [command]

Available Commands:
  analyze     Analyze something
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  jira        Interact with jira
  salesforce  Interact with salesforce

Flags:
      --config-file string   path to config file
      --debug                specify debug level
  -h, --help                 help for dataengctl

Use "dataengctl [command] --help" for more information about a command.
```

As a result, you can see that there are several commands that are available

- 1. analyze
- 2. completion
- 3. help
- 4. jira
- 5. salesforce

To interact with Jira and Salesforce or anything that is using the `config.yaml` mappings, you will need to specify the command as follows:

- 1. Intreracting with `analyze` command specifying the config path

```bash
. analyze --issue-type Escalation --project-key FIELD --config-file ${HOME}/.dataeng/config.yaml | jq "."
```

- 2. Interacting with `jira` command specifying the config path

```bash
. jira issue list --issue-type Escalation --project-key FIELD --config-file ${HOME}/.dataeng/config.yaml | jq ".pri"
```

- 3. Interacting with `salesforce` command specifying the config path

```bash
. salesforce query "SELECT+name+from+Account" --config-file ${HOME}/.dataeng/config.yaml | jq "."
```

Without `jq` installed on your machine the default `JSON` output will not look pretty. All output is default to `JSON` format.
