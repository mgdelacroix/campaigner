package cmd

func PullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "pull",
		Short: "Imports tickets from Jira",
		Long: "Imports all tickets from a Jira epic issue",
		RunE: pullCmdF,
	}

	cmd.Flags().BoolP("epic")
}
