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

### Creating a campaign

To create a campaign, first go to the directory where you want to create the campaign files.

First you need to create a template for the tickets of the campaign. This is an example of a template:

```
This ticket is for removing the usage of "// ToDo" in the source code of my project. Please go to "{{.filename}}", line {{.lineNo}} and remove the corresponding comment.
```

The template will be filled by go, so you can use the `{{}}` placeholders with the properties that your campaign tickets have been created with. This properties will depend on how you create the tickets: if using `govet` or `grep`, the tickets will be created with the `filename`, `lineNo` and `text` properties. If importing a CSV, they will have the same properties that columns has the CSV.

The template body can use as well the [JIRA text formatting notation](https://jira.atlassian.com/secure/WikiRendererHelpAction.jspa?section=all) to add rich formatting to the ticket. This format will be later transformed automatically to Markdown when publishing the tickets to GitHub.

Lastly, before creating the campaign, we can add a template to act as a footer for the tickets when being created in GitHub. In this case, the template will receive the ticket struct, so we can use any of its properties. This footer template supports GitHub Markdown.

Once we have both files, we can run `campaigner init` to create the campaign:

```sh
$ campaigner init --help
Creates a new campaign in the current directory

Usage:
  campaigner init [flags]

Examples:
  campaigner init \
    --jira-username johndoe \
    --jira-token secret \
    --github-token TOKEN \
    --url http://my-jira-instance.com \
    --epic ASD-27 \
    --issue-type Story \
    --repository johndoe/awesomeproject \
    -l 'Area/API' -l 'Tech/Go' \
    --summary 'Refactor {{.function}} to inject the configuration service' \
    --issue-template ./refactor-config.tmpl \
    --footer-template ./github-footer.tmpl


Flags:
  -e, --epic string              the epic id to associate this campaign with
  -f, --footer-template string   the template path to append to the github issues as a footer
      --github-token string      the github token
  -h, --help                     help for init
  -t, --issue-template string    the template path for the description of the tickets
  -i, --issue-type string        the issue type to create the tickets as (default "Story")
      --jira-token string        the jira token or password
      --jira-username string     the jira username
  -l, --label strings            the labels to add to the Github issues
  -r, --repository string        the github repository
  -s, --summary string           the summary of the tickets
  -u, --url string               the jira server URL
```

If there is any mandatory flag that we don't use, we will be asked for it interactively by the tool.

The `summary` of the campaign can be a go template as well, that will receive the same properties described above for the issue template.

### Adding tickets to the campaign

There are currently three ways to add tickets to a campaign: using `govet`, using `grep` or importing a `csv`. We can add tickets at any point during the campaign lifecycle, and we can use different methods in the same campaign, just remember that the summary and the templates will be filled with whatever properties we have in each ticket, so all of them should have at least those used in the templates.

You can see more information on how each method works using `campaigner add --help` and the help for each of the subcommands.

### Publishing tickets

Once the campaign is ready, we can see its status running `campaigner status`:

```sh
Current campaign for johndoe/testrepo with summary
Remove the ToDo comment in {{.filename}}:{{.lineNo}}

        351     total tickets
        0/351   published in Jira
        0/0     published in Github

```

The status shows the total of tickets that we have created, how many of those have been published in Jira and of that last amount, how many have been published in GitHub too.

To publish tickets, you can use the `campaigner publish jira` and `campaigner publish github` commands, and you can publish tickets in batches or just publish all of them. When running a publish command, `campaigner` will search the first unpublished tickets and will use the provider APIs to publish them.
