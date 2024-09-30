package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mgdelacroix/campaigner/app"
	"github.com/mgdelacroix/campaigner/model"
	"github.com/spf13/cobra"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func UICmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ui",
		Short: "Shows the campaigner UI",
		Args:  cobra.NoArgs,
		RunE:  withAppE(uiCmdF),
	}
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type uiModel struct {
	list list.Model
}

func (m uiModel) Init() tea.Cmd {
	return nil
}

func (m uiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+q" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m uiModel) View() string {
	return docStyle.Render(m.list.View())
}

func renderTicketDescription(ticket *model.Ticket) string {
	var status, bgColor string
	if ticket.IsClosed() {
		status = "closed"
		bgColor = "#DB0202"
	} else if ticket.IsAssigned() {
		status = fmt.Sprintf("assigned [%s]", ticket.GithubAssignee)
		bgColor = "#146300"
	} else if ticket.IsPublishedGithub() {
		status = "GitHub"
		bgColor = "#171515"
	} else if ticket.IsPublishedJira() {
		status = "Jira"
		bgColor = "#0052CC"
	} else {
		status = "unpublished"
		bgColor = "#8C6700"
	}

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color(bgColor)).
		Padding(0, 2)

	return fmt.Sprintf("Status: %s", style.Render(status))
}

func uiCmdF(a *app.App, cmd *cobra.Command, _ []string) error {
	items := make([]list.Item, len(a.Campaign.Tickets))
	for i, ticket := range a.Campaign.Tickets {
		title := ticket.Summary
		if title == "" {
			b, err := json.Marshal(ticket.Data)
			if err != nil {
				return fmt.Errorf("cannot marshal ticket data: %w", err)
			}
			title = string(b)
		}
		items[i] = item{title, renderTicketDescription(ticket)}
	}

	delegate := list.NewDefaultDelegate()

	m := uiModel{list: list.New(items, delegate, 0, 0)}
	m.list.Title = fmt.Sprintf("Campaign: %q", a.Campaign.GetName())

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
