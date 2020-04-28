package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"git.ctrlz.es/mgdelacroix/campaigner/campaign"
	"git.ctrlz.es/mgdelacroix/campaigner/model"

	"github.com/spf13/cobra"
)

func InitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Creates a campaign",
		Long:  "Creates a new campaign in the current directory",
		Example: `  campaigner init \
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
`,
		Args: cobra.NoArgs,
		Run:  initCmdF,
	}

	cmd.Flags().String("jira-username", "", "the jira username")
	cmd.Flags().String("jira-token", "", "the jira token or password")
	cmd.Flags().String("github-token", "", "the github token")
	cmd.Flags().StringP("url", "u", "", "the jira server URL")
	cmd.Flags().StringP("epic", "e", "", "the epic id to associate this campaign with")
	cmd.Flags().StringP("repository", "r", "", "the github repository")
	cmd.Flags().StringSliceP("label", "l", []string{}, "the labels to add to the Github issues")
	cmd.Flags().StringP("summary", "s", "", "the summary of the tickets")
	cmd.Flags().StringP("issue-template", "t", "", "the template path for the description of the tickets")
	cmd.Flags().StringP("footer-template", "f", "", "the template path to append to the github issues as a footer")
	cmd.Flags().StringP("issue-type", "i", "Story", "the issue type to create the tickets as")

	return cmd
}

func initCmdF(cmd *cobra.Command, _ []string) {
	getStringFlagOrAskIfEmpty := func(name string, question string) string {
		val, _ := cmd.Flags().GetString(name)
		if val == "" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("%s ", question)
			answer, err := reader.ReadString('\n')
			if err != nil {
				ErrorAndExit(cmd, err)
			}
			val = strings.TrimSpace(answer)
		}
		return val
	}

	jiraUsername := getStringFlagOrAskIfEmpty("jira-username", "JIRA username:")
	jiraToken := getStringFlagOrAskIfEmpty("jira-token", "JIRA password or token:")
	githubToken := getStringFlagOrAskIfEmpty("github-token", "GitHub token:")
	url := getStringFlagOrAskIfEmpty("url", "JIRA server URL:")
	epic := getStringFlagOrAskIfEmpty("epic", "JIRA epic:")
	repo := getStringFlagOrAskIfEmpty("repository", "GitHub repository:")
	summary := getStringFlagOrAskIfEmpty("summary", "Ticket summary template:")
	issueTemplate := getStringFlagOrAskIfEmpty("issue-template", "Ticket description template path:")
	footerTemplate := getStringFlagOrAskIfEmpty("footer-template", "Github issue footer template path:")
	issueType, _ := cmd.Flags().GetString("issue-type")
	labels, _ := cmd.Flags().GetStringSlice("label")

	project := strings.Split(epic, "-")[0]

	cmp := &model.Campaign{
		Jira: model.ConfigJira{
			Url:       url,
			Username:  jiraUsername,
			Token:     jiraToken,
			Project:   project,
			Epic:      epic,
			IssueType: issueType,
		},
		Github: model.ConfigGithub{
			Token:  githubToken,
			Repo:   repo,
			Labels: labels,
		},
		Summary:        summary,
		IssueTemplate:  issueTemplate,
		FooterTemplate: footerTemplate,
	}
	if err := campaign.Save(cmp); err != nil {
		ErrorAndExit(cmd, err)
	}
}
