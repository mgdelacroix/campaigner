# Campaigner

Command line tool to create and manage community campaigns.

## Usage

```sh
$ campaigner --help
Create and manage Open Source campaigns

Usage:
  campaigner [command]

Available Commands:
  add         Adds tickets to the campaign
  filter      Interactively filters the current ticket list
  help        Help about any command
  init        Creates a new campaign in the current directory
  token       Subcommands related to tokens

Flags:
  -h, --help   help for campaigner

Use "campaigner [command] --help" for more information about a command.
```

## Application flow

1. Run `campaigner set-token` to set the tokens for `jira/github`.
2. Run `campaigner init` to create a new community campaign.
3. Run `campaigner add` to add new tickets based either on a `grep/ag`
   command or a `govet` check.
4. Run `campaigner filter` to interactively remove false matches.
5. Run `campaigner status` to see the status of each one of the
   tickets. The tickets can be in an `unpublished`, `jira`, `github`
   and `completed` state.
6. Modify the `template.md` file to adjust the ticket templates.
7. Run `campaigner publish` to create the tickets in `jira` based on
   the template.
8. Possible next step to publish the tickets from `jira` to `github`.
