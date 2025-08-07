package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screen int

const (
	serviceListScreen screen = iota
	serviceNameScreen
	featureSelectionScreen
	featureDetailsScreen
	summaryScreen
)

type Model struct {
	screen         screen
	outputFile     string
	config         *Config
	
	// Service list
	servicesCursor int
	editingService bool
	
	// Service input
	serviceInput   string
	currentService string
	
	// Feature selection
	featuresCursor int
	selectedFeatures map[string]bool
	
	// Feature details
	detailInput    string
	currentFeature string
	
	// State
	completed      bool
	err            error
}

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)
		
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		MarginBottom(1)
		
	selectedStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Background(lipgloss.Color("#E5E5E5"))
		
	normalStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666"))
		
	inputStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(0, 1)
		
	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		MarginTop(1)
)

func NewModel(outputFile string) Model {
	return NewModelWithConfig(outputFile, &Config{Services: make(map[string]Service)})
}

func NewModelWithConfig(outputFile string, config *Config) Model {
	// Start with service list if we have existing services, otherwise start with service name input
	startScreen := serviceListScreen
	if len(config.Services) == 0 {
		startScreen = serviceNameScreen
	}
	
	return Model{
		screen:           startScreen,
		outputFile:       outputFile,
		config:           config,
		selectedFeatures: make(map[string]bool),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
		
		switch m.screen {
		case serviceListScreen:
			return m.updateServiceListScreen(msg)
		case serviceNameScreen:
			return m.updateServiceNameScreen(msg)
		case featureSelectionScreen:
			return m.updateFeatureSelectionScreen(msg)
		case featureDetailsScreen:
			return m.updateFeatureDetailsScreen(msg)
		case summaryScreen:
			return m.updateSummaryScreen(msg)
		}
	}
	
	return m, nil
}

func (m Model) View() string {
	var content string
	
	switch m.screen {
	case serviceListScreen:
		content = m.viewServiceListScreen()
	case serviceNameScreen:
		content = m.viewServiceNameScreen()
	case featureSelectionScreen:
		content = m.viewFeatureSelectionScreen()
	case featureDetailsScreen:
		content = m.viewFeatureDetailsScreen()
	case summaryScreen:
		content = m.viewSummaryScreen()
	}
	
	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("HCL Configuration Builder"),
		"",
		content,
	)
}