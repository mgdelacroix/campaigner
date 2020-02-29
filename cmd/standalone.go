package cmd

func StandaloneCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use: "standalone",
        Short: "Standalone fire-and-forget commands",
    }
    
    cmd.AddCommand(
        CreateJiraTicketStandaloneCmd(),
    )
    
    return cmd
}

func CreateJiraTicketStandaloneCmd() *cobra.Command{
    cmd := &cobra.Command{
        Use: "create-jira-ticket",
        Short: "Creates a jira ticket from a template",
        Args: cobra.NoArgs,
        Run: createJiraTicketStandaloneCmdF,
    }
    
    cmd.Flags().StringP("username", "u", "", "the jira username")
    _ = cmd.MarkFlagRequired("username")
    cmd.Flags().StringP("token", "t", "", "the jira token")
    _ = cmd.MarkFlagRequired("token")
    cmd.Flags().StringP("summary", "s", "", "the summary of the ticket")
    _ = cmd.MarkFlagRequired("summary")
    cmd.Flags().StringP("template", "m", "", "the template to render the description of the ticket")
    _ = cmd.MarkFlagRequired("template")
    cmd.Flags().StringSliceP("vars", "v", "", "the variables to use in the template")
    
    return cmd
}

func createJiraTicketStandaloneCmdF(cmd *cobra.Command, _ []string) {
    
}