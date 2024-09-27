# Campaigner

Command line tool to create and manage community campaigns. `campaigner` takes care of all the campaign lifecycle, starting with getting information to generate the tickets, then publishing them in jira and github and finally tracking their status, the campaign progress and generating reports.

 - `campaigner init` generates the campaign file, linking it to the jira instance and github repository and to the epic issue that will host each campaign ticket.
 - `campaigner add` parses information from different sources and uses it to generate tickets for the campaign.
 - `campaigner label` lists and modifies the campaign's labels for GitHub.
 - `campaigner publish` builds the tickets information and publishes it both to jira and github.
 - `campaigner sync` downloads updated information of the campaign progress.
 - `campaigner status` shows the current campaign data and progression.
 - `campaigner list` shows the current campaign tickets and their status.
 - `campaigner report` generates reports from the campaign data.

## Install

To install `campaigner`, if you have the golang environment set up, you just have to run:

```sh
go install github.com/mgdelacroix/campaigner@latest
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

Once we have both files, we can run `campaigner init` to create the campaign. The command can be run without any arguments or flags and it will request the mandatory pieces of information interactively:

```sh
$ campaigner init \
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
```

The `summary` of the campaign can be a go template as well, that will receive the same properties described above for the issue template.

### Modifying campaign GitHub labels

`campaigner` can add labels to the GitHub tickets when creating them. Labels can be added to the campaign as part of the `init` command, but they can be managed as well when the campaign is already created with the `campaigner label` commands:

 - `campaigner label list` lists the current campaign labels.
 - `campaigner label remote` lists the labels that exist in the remote GitHub repository.
 - `campaigner label update` opens `$EDITOR` with the campaign labels and allows you to add / remove / edit them, saving them when the editor closes. It will error if a label doesn't exist in the remote repository, but you can skip this check adding the `--skip-check` flag.

### Adding tickets to the campaign

There are currently three ways to add tickets to a campaign: using `govet`, using `grep` or importing a `csv`. We can add tickets at any point during the campaign lifecycle, and we can use different methods in the same campaign, just remember that the summary and the templates will be filled with whatever properties we have in each ticket, so all of them should have at least those used in the templates.

You can see more information on how each method works using `campaigner add --help` and the help for each of the subcommands.

### Publishing tickets

Once the campaign is ready, we can see its status running `campaigner status`:

```sh
$ campaigner status
Current campaign for johndoe/testrepo with summary
Remove the ToDo comment in {{.filename}}:{{.lineNo}}

         67     -         total tickets
         24   35%     published in Jira
         24   35%   published in Github
         17   25%              assigned
          0    0%                closed

```

The status shows the total of tickets that we have created, how many of those have been published in Jira and of that last amount, how many have been published in GitHub too.

To publish tickets, you can use the `campaigner publish jira` and `campaigner publish github` commands, and you can publish tickets in batches or just publish all of them. When running a publish command, `campaigner` will search the first unpublished tickets and will use the provider APIs to publish them.

### Syncing the campaign status

If we want to check the status of our campaign as it progresses and the tickets get assigned and closed, we have first to sync our local state of the campaign with the status of the tickets we've published:

```sh
$ campaigner sync
Updating ticket 9 of 9
Synchronization completed
```

This will fetch the current state of the jira tickets and github issues and update our local campaign. Then we can run `campaigner status` to see the updated progress, or we can generate reports from the campaign's information with `campaigner report`.
