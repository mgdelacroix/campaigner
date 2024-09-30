package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mgdelacroix/campaigner/app"
	"github.com/mgdelacroix/campaigner/model"

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

	cmd.Flags().String("jira-username", "", "the Jira username")
	cmd.Flags().String("jira-token", "", "the Jira token or password")
	cmd.Flags().String("github-token", "", "the GitHub token")
	cmd.Flags().StringP("url", "u", "", "the Jira server URL")
	cmd.Flags().StringP("epic", "e", "", "the epic id to associate this campaign with")
	cmd.Flags().StringP("repository", "r", "", "the GitHub repository")
	cmd.Flags().StringSliceP("label", "l", []string{}, "the labels to add to the Github issues")
	cmd.Flags().StringP("summary", "s", "", "the summary of the tickets")
	cmd.Flags().StringP("issue-template", "t", "", "the template path for the description of the tickets")
	cmd.Flags().StringP("footer-template", "f", "", "the template path to append to the GitHub issues as a footer")
	cmd.Flags().StringP("issue-type", "i", "Story", "the issue type to create the tickets as")

	return cmd
}

func initCmdF(cmd *cobra.Command, _ []string) {
	campaignPath, _ := cmd.Flags().GetString("campaign")

	_, err := os.Stat(campaignPath)
	if err == nil {
		ErrorAndExit(cmd, fmt.Errorf("cannot use %s as campaign file: file already exists", campaignPath))
	} else if !os.IsNotExist(err) {
		ErrorAndExit(cmd, fmt.Errorf("cannot use %s as campaign file: %w", campaignPath, err))
	}

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

	name := getStringFlagOrAskIfEmpty("name", "Campaign name:")
	jiraUsername := getStringFlagOrAskIfEmpty("jira-username", "Jira username:")
	jiraToken := getStringFlagOrAskIfEmpty("jira-token", "Jira password or token:")
	githubToken := getStringFlagOrAskIfEmpty("github-token", "GitHub token:")
	url := getStringFlagOrAskIfEmpty("url", "Jira server URL:")
	epic := getStringFlagOrAskIfEmpty("epic", "Jira epic:")
	repo := getStringFlagOrAskIfEmpty("repository", "GitHub repository:")
	summary := getStringFlagOrAskIfEmpty("summary", "Ticket summary template:")
	issueTemplate := getStringFlagOrAskIfEmpty("issue-template", "Ticket description template path:")
	footerTemplate := getStringFlagOrAskIfEmpty("footer-template", "Github issue footer template path:")
	issueType, _ := cmd.Flags().GetString("issue-type")
	labels, _ := cmd.Flags().GetStringSlice("label")

	project := strings.Split(epic, "-")[0]

	campaign := model.NewCampaign(name)
	campaign.Jira = model.ConfigJira{
		Url:       url,
		Username:  jiraUsername,
		Token:     jiraToken,
		Project:   project,
		Epic:      epic,
		IssueType: issueType,
	}
	campaign.Github = model.ConfigGithub{
		Token:  githubToken,
		Repo:   repo,
		Labels: labels,
	}
	campaign.Summary = summary
	campaign.IssueTemplate = issueTemplate
	campaign.FooterTemplate = footerTemplate
	if err := app.SaveCampaign(campaign, campaignPath); err != nil {
		ErrorAndExit(cmd, err)
	}
}
