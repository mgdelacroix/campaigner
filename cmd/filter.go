package cmd

func FilterCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "filter",
		Short: "Interactively filters the current ticket list",
		Args:  cobra.NoArgs,
		Run:   filterCmdF,
	}
}

func filterCmdF(_ *cobra.Command, _ []string) {}
