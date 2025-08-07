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
	unsavedChangesScreen
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
	
	// Feature selection with search
	featuresCursor   int
	selectedFeatures map[string]bool
	featureValues    map[string]string // Store actual feature values (like "viewer" for firebaseauth)
	originalFeatures map[string]bool   // Original state for change detection
	originalValues   map[string]string // Original values for change detection
	searchMode       bool
	searchInput      string
	filteredFeatures []string
	viewportOffset   int  // For scrolling
	
	// Feature details
	detailInput    string
	currentFeature string
	choiceCursor   int
	choices        []string
	
	// Terminal dimensions
	width  int
	height int
	
	// State
	completed      bool
	err            error
}

var (
	// Dark terminal-optimized color palette using simple color names
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("0")).  // Black text
		Background(lipgloss.Color("205")). // Bright pink/coral
		Padding(0, 2).
		MarginBottom(1)
		
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("51")). // Bright cyan
		MarginBottom(1).
		Underline(true)
		
	selectedStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("0")).   // Black text
		Background(lipgloss.Color("51")).  // Bright cyan background
		Padding(0, 1)
		
	normalStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")) // Bright white for dark terminals
		
	inputStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("51")). // Bright cyan
		Foreground(lipgloss.Color("255")).      // White text
		Padding(0, 1).
		MarginBottom(1)
		
	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("247")). // Light gray
		Italic(true).
		MarginTop(1)
		
	// Semantic styles using safe color codes
	successStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("46")) // Bright green
		
	warningStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("214")) // Bright orange
		
	errorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("196")) // Bright red
		
	accentStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("141")) // Bright purple
		
	searchStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("214")). // Orange border
		Foreground(lipgloss.Color("0")).         // Black text
		Background(lipgloss.Color("229")).       // Light yellow background
		Padding(0, 1).
		MarginBottom(1)
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
		featureValues:    make(map[string]string),
		originalFeatures: make(map[string]bool),
		originalValues:   make(map[string]string),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
		
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
		case unsavedChangesScreen:
			return m.updateUnsavedChangesScreen(msg)
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
	case unsavedChangesScreen:
		content = m.viewUnsavedChangesScreen()
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