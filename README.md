# Campaigner

Command line tool to create and manage community campaigns. `campaigner` takes care of all the campaign lifecycle, starting with getting information to generate the tickets, then publishing them in jira and github and finally tracking their status, the campaign progress and generating reports.

 - `campaigner init` generates the campaign file, linking it to the jira instance and github repository and to the epic issue that will host each campaign ticket.
 - `campaigner add` parses information from different sources and uses it to generate tickets for the campaign.
 - `campaigner publish` builds the tickets information and publishes it both to jira and github.
 - `campaigner sync` downloads updated information of the campaign progress.
 - `campaigner status` shows the current campaign data and progression.
 - `campaigner report` generates reports from the campaign data.

## Install

To install `campaigner`, if you have the golang environment set up, you just have to run:

```sh
go get git.ctrlz.es/mgdelacroix/campaigner
```

## Usage

```sh
$ campaigner --help
Create and manage Open Source campaigns

Usage:
  campaigner [command]

Available Commands:
  add         Adds tickets to the campaign from the output of grep/ag/govet
  filter      Interactively filters the current ticket list
  help        Help about any command
  init        Creates a new campaign in the current directory
  publish     Publishes the campaign tickets in different providers
  standalone  Standalone fire-and-forget commands
  status      Prints the current status of the campaign
  sync        Synchronizes the status of the tickets with remote providers
  token       Subcommands related to tokens

Flags:
  -h, --help   help for campaigner

Use "campaigner [command] --help" for more information about a command.
```
