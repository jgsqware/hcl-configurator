package main

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sahilm/fuzzy"
)

// Service List Screen
func (m Model) updateServiceListScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	services := make([]string, 0, len(m.config.Services))
	for name := range m.config.Services {
		services = append(services, name)
	}
	sort.Strings(services)
	
	switch msg.String() {
	case "up", "k":
		if m.servicesCursor > 0 {
			m.servicesCursor--
		}
		return m, nil
		
	case "down", "j":
		if m.servicesCursor < len(services) {
			m.servicesCursor++
		}
		return m, nil
		
	case "enter":
		if m.servicesCursor < len(services) {
			// Edit existing service
			serviceName := services[m.servicesCursor]
			m.currentService = serviceName
			m.editingService = true
			m.selectedFeatures = make(map[string]bool)
			m.featureValues = make(map[string]string)
			m.originalFeatures = make(map[string]bool)
			m.originalValues = make(map[string]string)
			m.loadExistingFeatures(m.config.Services[serviceName])
			m.screen = featureSelectionScreen
			return m, nil
		} else {
			// Add new service
			m.editingService = false
			m.selectedFeatures = make(map[string]bool)
			m.featureValues = make(map[string]string)
			m.originalFeatures = make(map[string]bool)
			m.originalValues = make(map[string]string)
			m.screen = serviceNameScreen
			return m, nil
		}
		
	case "d":
		if m.servicesCursor < len(services) {
			// Delete service
			serviceName := services[m.servicesCursor]
			delete(m.config.Services, serviceName)
			if m.servicesCursor >= len(services)-1 {
				m.servicesCursor = len(services) - 2
				if m.servicesCursor < 0 {
					m.servicesCursor = 0
				}
			}
			return m, nil
		}
		return m, nil
		
	case "s":
		// Go to summary
		m.screen = summaryScreen
		return m, nil
		
	case "esc":
		return m, tea.Quit
	}
	
	return m, nil
}

func (m Model) viewServiceListScreen() string {
	var content strings.Builder
	
	content.WriteString(headerStyle.Render("Cloud Run Services Configuration"))
	content.WriteString("\n\n")
	
	services := make([]string, 0, len(m.config.Services))
	for name := range m.config.Services {
		services = append(services, name)
	}
	sort.Strings(services)
	
	if len(services) == 0 {
		content.WriteString(warningStyle.Render("No services configured yet"))
		content.WriteString("\n\n")
		content.WriteString(selectedStyle.Render("> Add new service"))
	} else {
		content.WriteString(accentStyle.Render(fmt.Sprintf("%d services configured:", len(services))))
		content.WriteString("\n\n")
		
		for i, serviceName := range services {
			service := m.config.Services[serviceName]
			featureCount := m.countFeatures(service.Features)
			
			var line string
			if i == m.servicesCursor {
				line = selectedStyle.Render(fmt.Sprintf("> %s (%d features)", serviceName, featureCount))
			} else {
				line = normalStyle.Render(fmt.Sprintf("  %s (%d features)", serviceName, featureCount))
			}
			
			content.WriteString(line)
			content.WriteString("\n")
		}
		
		content.WriteString("\n")
		// Add "Add new service" option
		if m.servicesCursor == len(services) {
			content.WriteString(selectedStyle.Render("> Add new service"))
		} else {
			content.WriteString(normalStyle.Render("  Add new service"))
		}
	}
	
	content.WriteString("\n\n")
	
	helpText := helpStyle.Render("Navigation: ↑/↓ navigate, Enter edit/add, d delete, s summary, Esc quit")
	content.WriteString(helpText)
	
	return content.String()
}

func (m Model) countFeatures(features Features) int {
	count := 0
	if features.FirebaseAuth != nil {
		count++
	}
	if features.FirebaseCloudMessagingSender {
		count++
	}
	if features.FirebaseCloudMessagingViewer {
		count++
	}
	if features.CloudRunInvoker {
		count++
	}
	if features.EventarcSubrole {
		count++
	}
	if features.CIDR != "" {
		count++
	}
	if len(features.BucketWriter) > 0 {
		count++
	}
	if len(features.BucketCreator) > 0 {
		count++
	}
	if len(features.BucketReader) > 0 {
		count++
	}
	if len(features.SubscriptionSubscriber) > 0 {
		count++
	}
	if len(features.SubscriptionViewer) > 0 {
		count++
	}
	if len(features.SubscriptionEditor) > 0 {
		count++
	}
	if len(features.TopicPublisher) > 0 {
		count++
	}
	if len(features.TopicViewer) > 0 {
		count++
	}
	if len(features.TopicEditor) > 0 {
		count++
	}
	if features.FirestoreReader {
		count++
	}
	if features.FirestoreWriter {
		count++
	}
	if features.MysqlAccess {
		count++
	}
	if features.PostgresAccess {
		count++
	}
	if features.EnableProfiling {
		count++
	}
	return count
}

func (m *Model) loadExistingFeatures(service Service) {
	features := service.Features
	
	// Load FirebaseAuth value
	if features.FirebaseAuth != nil {
		m.selectedFeatures["firebaseauth"] = true
		m.originalFeatures["firebaseauth"] = true
		if val, ok := features.FirebaseAuth.(string); ok {
			m.featureValues["firebaseauth"] = val
			m.originalValues["firebaseauth"] = val
		} else if val, ok := features.FirebaseAuth.(bool); ok && val {
			m.featureValues["firebaseauth"] = "viewer" // default for boolean true
			m.originalValues["firebaseauth"] = "viewer"
		}
	}
	
	// Load CIDR value
	if features.CIDR != "" {
		m.selectedFeatures["cidr"] = true
		m.originalFeatures["cidr"] = true
		m.featureValues["cidr"] = features.CIDR
		m.originalValues["cidr"] = features.CIDR
	}
	
	// Load bucket features (store as comma-separated strings)
	if len(features.BucketWriter) > 0 {
		m.selectedFeatures["bucket_writer"] = true
		m.originalFeatures["bucket_writer"] = true
		val := strings.Join(features.BucketWriter, ",")
		m.featureValues["bucket_writer"] = val
		m.originalValues["bucket_writer"] = val
	}
	if len(features.BucketCreator) > 0 {
		m.selectedFeatures["bucket_creator"] = true
		m.originalFeatures["bucket_creator"] = true
		val := strings.Join(features.BucketCreator, ",")
		m.featureValues["bucket_creator"] = val
		m.originalValues["bucket_creator"] = val
	}
	if len(features.BucketReader) > 0 {
		m.selectedFeatures["bucket_reader"] = true
		m.originalFeatures["bucket_reader"] = true
		val := strings.Join(features.BucketReader, ",")
		m.featureValues["bucket_reader"] = val
		m.originalValues["bucket_reader"] = val
	}
	
	// Load subscription features
	if len(features.SubscriptionSubscriber) > 0 {
		m.selectedFeatures["subscription_subscriber"] = true
		m.originalFeatures["subscription_subscriber"] = true
		val := strings.Join(features.SubscriptionSubscriber, ",")
		m.featureValues["subscription_subscriber"] = val
		m.originalValues["subscription_subscriber"] = val
	}
	if len(features.SubscriptionViewer) > 0 {
		m.selectedFeatures["subscription_viewer"] = true
		m.originalFeatures["subscription_viewer"] = true
		val := strings.Join(features.SubscriptionViewer, ",")
		m.featureValues["subscription_viewer"] = val
		m.originalValues["subscription_viewer"] = val
	}
	if len(features.SubscriptionEditor) > 0 {
		m.selectedFeatures["subscription_editor"] = true
		m.originalFeatures["subscription_editor"] = true
		val := strings.Join(features.SubscriptionEditor, ",")
		m.featureValues["subscription_editor"] = val
		m.originalValues["subscription_editor"] = val
	}
	
	// Load topic features
	if len(features.TopicPublisher) > 0 {
		m.selectedFeatures["topic_publisher"] = true
		m.originalFeatures["topic_publisher"] = true
		val := strings.Join(features.TopicPublisher, ",")
		m.featureValues["topic_publisher"] = val
		m.originalValues["topic_publisher"] = val
	}
	if len(features.TopicViewer) > 0 {
		m.selectedFeatures["topic_viewer"] = true
		m.originalFeatures["topic_viewer"] = true
		val := strings.Join(features.TopicViewer, ",")
		m.featureValues["topic_viewer"] = val
		m.originalValues["topic_viewer"] = val
	}
	if len(features.TopicEditor) > 0 {
		m.selectedFeatures["topic_editor"] = true
		m.originalFeatures["topic_editor"] = true
		val := strings.Join(features.TopicEditor, ",")
		m.featureValues["topic_editor"] = val
		m.originalValues["topic_editor"] = val
	}
	
	// Boolean features (no values to store, just track selection)
	if features.FirebaseCloudMessagingSender {
		m.selectedFeatures["firebase_cloudmessaging_sender"] = true
		m.originalFeatures["firebase_cloudmessaging_sender"] = true
	}
	if features.FirebaseCloudMessagingViewer {
		m.selectedFeatures["firebase_cloudmessaging_viewer"] = true
		m.originalFeatures["firebase_cloudmessaging_viewer"] = true
	}
	if features.CloudRunInvoker {
		m.selectedFeatures["cloudrun_invoker"] = true
		m.originalFeatures["cloudrun_invoker"] = true
	}
	if features.EventarcSubrole {
		m.selectedFeatures["eventarc_subrole"] = true
		m.originalFeatures["eventarc_subrole"] = true
	}
	if features.FirestoreReader {
		m.selectedFeatures["firestore_reader"] = true
		m.originalFeatures["firestore_reader"] = true
	}
	if features.FirestoreWriter {
		m.selectedFeatures["firestore_writer"] = true
		m.originalFeatures["firestore_writer"] = true
	}
	if features.MysqlAccess {
		m.selectedFeatures["mysql_access"] = true
		m.originalFeatures["mysql_access"] = true
	}
	if features.PostgresAccess {
		m.selectedFeatures["postgres_access"] = true
		m.originalFeatures["postgres_access"] = true
	}
	if features.EnableProfiling {
		m.selectedFeatures["enable_profiling"] = true
		m.originalFeatures["enable_profiling"] = true
	}
}

// Service Name Screen
func (m Model) updateServiceNameScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if strings.TrimSpace(m.serviceInput) == "" {
			// If empty input, go back to service list
			m.screen = serviceListScreen
			return m, nil
		}
		
		// Set current service and go to feature selection
		m.currentService = strings.TrimSpace(m.serviceInput)
		m.serviceInput = ""
		m.editingService = false
		m.selectedFeatures = make(map[string]bool)
		m.featureValues = make(map[string]string)
		m.originalFeatures = make(map[string]bool)
		m.originalValues = make(map[string]string)
		m.screen = featureSelectionScreen
		return m, nil
		
	case "backspace":
		if len(m.serviceInput) > 0 {
			m.serviceInput = m.serviceInput[:len(m.serviceInput)-1]
		}
		return m, nil
		
	default:
		if len(msg.String()) == 1 {
			m.serviceInput += msg.String()
		}
		return m, nil
	}
}

func (m Model) viewServiceNameScreen() string {
	var content strings.Builder
	
	content.WriteString(headerStyle.Render("Enter Service Name"))
	content.WriteString("\n\n")
	
	if len(m.config.Services) > 0 {
		content.WriteString(accentStyle.Render("Existing services:"))
		content.WriteString("\n")
		for name := range m.config.Services {
			content.WriteString(successStyle.Render("  " + name))
			content.WriteString("\n")
		}
		content.WriteString("\n")
	}
	
	promptText := "Service name: " + m.serviceInput
	content.WriteString(inputStyle.Render(promptText))
	content.WriteString("\n")
	
	var helpText string
	if len(m.config.Services) > 0 {
		helpText = "Press Enter to add service, or Enter with empty name to go back"
	} else {
		helpText = "Enter a service name and press Enter to continue"
	}
	content.WriteString(helpStyle.Render(helpText))
	
	return content.String()
}

// Get all features with compact descriptions
func getAllFeatures() []FeatureItem {
	return []FeatureItem{
		{"firebaseauth", "Firebase Authentication"},
		{"firebase_cloudmessaging_sender", "Firebase Messaging Sender"},
		{"firebase_cloudmessaging_viewer", "Firebase Messaging Viewer"},
		{"cloudrun_invoker", "Cloud Run Invoker"},
		{"eventarc_subrole", "Eventarc Subrole"},
		{"cidr", "CIDR Block Access"},
		{"bucket_writer", "Storage Bucket Writer"},
		{"bucket_creator", "Storage Bucket Creator"},
		{"bucket_reader", "Storage Bucket Reader"},
		{"subscription_subscriber", "Pub/Sub Subscriber"},
		{"subscription_viewer", "Pub/Sub Subscription Viewer"},
		{"subscription_editor", "Pub/Sub Subscription Editor"},
		{"topic_publisher", "Pub/Sub Topic Publisher"},
		{"topic_viewer", "Pub/Sub Topic Viewer"},
		{"topic_editor", "Pub/Sub Topic Editor"},
		{"firestore_reader", "Firestore Reader"},
		{"firestore_writer", "Firestore Writer"},
		{"mysql_access", "MySQL Access"},
		{"postgres_access", "PostgreSQL Access"},
		{"enable_profiling", "Application Profiling"},
	}
}

type FeatureItem struct {
	Name        string
	Description string
}

func (f FeatureItem) String() string {
	return f.Name + " " + f.Description
}

// Feature Selection Screen
func (m Model) updateFeatureSelectionScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	allFeatures := getAllFeatures()
	
	// Initialize filtered features if needed
	if !m.searchMode && len(m.filteredFeatures) == 0 {
		m.filteredFeatures = make([]string, len(allFeatures))
		for i, f := range allFeatures {
			m.filteredFeatures[i] = f.Name
		}
	}
	
	// Handle search mode input
	if m.searchMode {
		switch msg.String() {
		case "enter", "esc":
			m.searchMode = false
			return m, nil
			
		case "backspace":
			if len(m.searchInput) > 0 {
				m.searchInput = m.searchInput[:len(m.searchInput)-1]
				m.updateFilteredFeatures()
				m.adjustViewport()
			}
			return m, nil
			
		default:
			if len(msg.String()) == 1 {
				m.searchInput += msg.String()
				m.updateFilteredFeatures()
				m.featuresCursor = 0
				m.viewportOffset = 0
			}
			return m, nil
		}
	}
	
	// Normal navigation mode
	switch msg.String() {
	case "/":
		m.searchMode = true
		m.searchInput = ""
		return m, nil
		
	case "up", "k":
		if m.featuresCursor > 0 {
			m.featuresCursor--
			m.adjustViewport()
		}
		return m, nil
		
	case "down", "j":
		if m.featuresCursor < len(m.filteredFeatures)-1 {
			m.featuresCursor++
			m.adjustViewport()
		}
		return m, nil
		
	case "ctrl+u": // Page up
		m.featuresCursor -= 10
		if m.featuresCursor < 0 {
			m.featuresCursor = 0
		}
		m.adjustViewport()
		return m, nil
		
	case "ctrl+d": // Page down
		m.featuresCursor += 10
		if m.featuresCursor >= len(m.filteredFeatures) {
			m.featuresCursor = len(m.filteredFeatures) - 1
		}
		m.adjustViewport()
		return m, nil
		
	case " ":
		if len(m.filteredFeatures) == 0 {
			return m, nil
		}
		
		feature := m.filteredFeatures[m.featuresCursor]
		needsInput := []string{
			"firebaseauth", "cidr", "bucket_creator", "bucket_reader", "bucket_writer",
			"subscription_subscriber", "subscription_viewer", "subscription_editor",
			"topic_publisher", "topic_viewer", "topic_editor",
		}
		
		requiresInput := false
		for _, inputFeature := range needsInput {
			if feature == inputFeature {
				requiresInput = true
				break
			}
		}
		
		if requiresInput {
			m.currentFeature = feature
			// Load existing value if available
			if existingValue, exists := m.featureValues[feature]; exists {
				m.detailInput = existingValue
			} else {
				m.detailInput = ""
			}
			m.choiceCursor = 0
			m.setupChoices(feature)
			// Set cursor to existing choice if it exists
			if m.detailInput != "" {
				for i, choice := range m.choices {
					if choice == m.detailInput {
						m.choiceCursor = i
						break
					}
				}
			}
			m.screen = featureDetailsScreen
			return m, nil
		} else {
			m.selectedFeatures[feature] = !m.selectedFeatures[feature]
		}
		return m, nil
		
	case "enter":
		service := Service{Features: m.buildFeatures()}
		m.config.Services[m.currentService] = service
		m.screen = serviceListScreen
		return m, nil
		
	case "esc":
		// Check for unsaved changes before exiting
		if m.hasUnsavedChanges() {
			m.screen = unsavedChangesScreen
		} else {
			// No changes, safe to exit
			if m.editingService {
				m.screen = serviceListScreen
			} else {
				m.screen = serviceNameScreen
			}
		}
		return m, nil
	}
	
	return m, nil
}

// Adjust viewport to keep cursor visible
func (m *Model) adjustViewport() {
	if m.height == 0 {
		return // Terminal size not initialized yet
	}
	
	// Calculate available height for features (total - header - status - search - help)
	availableHeight := m.height - 5
	if m.searchMode {
		availableHeight-- // Account for search line
	}
	
	if availableHeight <= 0 {
		availableHeight = 10 // Minimum
	}
	
	// Adjust viewport offset to keep cursor visible
	if m.featuresCursor < m.viewportOffset {
		m.viewportOffset = m.featuresCursor
	} else if m.featuresCursor >= m.viewportOffset+availableHeight {
		m.viewportOffset = m.featuresCursor - availableHeight + 1
	}
	
	// Ensure viewport doesn't go negative
	if m.viewportOffset < 0 {
		m.viewportOffset = 0
	}
}

func (m *Model) updateFilteredFeatures() {
	allFeatures := getAllFeatures()
	
	if m.searchInput == "" {
		m.filteredFeatures = make([]string, len(allFeatures))
		for i, f := range allFeatures {
			m.filteredFeatures[i] = f.Name
		}
		return
	}
	
	// Create searchable strings
	searchableFeatures := make([]string, len(allFeatures))
	for i, f := range allFeatures {
		searchableFeatures[i] = f.String() // Use the String() method which combines name and description
	}
	
	// Use fuzzy search
	matches := fuzzy.Find(m.searchInput, searchableFeatures)
	m.filteredFeatures = make([]string, len(matches))
	for i, match := range matches {
		m.filteredFeatures[i] = allFeatures[match.Index].Name
	}
}

func (m Model) viewFeatureSelectionScreen() string {
	allFeatures := getAllFeatures()
	descriptions := make(map[string]string)
	for _, f := range allFeatures {
		descriptions[f.Name] = f.Description
	}
	
	// Calculate terminal dimensions
	availableHeight := m.height - 5 // Reserve space for header, status, help
	if m.searchMode {
		availableHeight-- // Reserve space for search
	}
	if availableHeight <= 0 {
		availableHeight = 10
	}
	
	var content strings.Builder
	
	// Compact title
	action := "Configure"
	if m.editingService {
		action = "Edit"
	}
	title := fmt.Sprintf("%s: %s", action, m.currentService)
	content.WriteString(headerStyle.Render(title))
	content.WriteString("\n")
	
	// Status line with counts
	selectedCount := 0
	for _, selected := range m.selectedFeatures {
		if selected {
			selectedCount++
		}
	}
	
	statusLine := fmt.Sprintf("[%d/%d]", selectedCount, len(allFeatures))
	if len(m.filteredFeatures) < len(allFeatures) {
		statusLine += fmt.Sprintf(" (%d matches)", len(m.filteredFeatures))
	}
	content.WriteString(accentStyle.Render(statusLine))
	content.WriteString("\n")
	
	// Features viewport
	features := m.filteredFeatures
	if len(features) == 0 {
		features = make([]string, len(allFeatures))
		for i, f := range allFeatures {
			features[i] = f.Name
		}
	}
	
	if len(features) == 0 {
		content.WriteString(warningStyle.Render("No matches"))
		content.WriteString("\n")
	} else {
		// Show only visible items in viewport
		start := m.viewportOffset
		end := start + availableHeight
		if end > len(features) {
			end = len(features)
		}
		
		for i := start; i < end; i++ {
			feature := features[i]
			isSelected := m.selectedFeatures[feature]
			isCurrent := i == m.featuresCursor
			
			// Compact display: just checkbox and description
			checkbox := " "
			if isSelected {
				checkbox = "x"
			}
			
			line := fmt.Sprintf("[%s] %s", checkbox, descriptions[feature])
			
			if isCurrent {
				content.WriteString(selectedStyle.Render(line))
			} else {
				if isSelected {
					content.WriteString(successStyle.Render(line))
				} else {
					content.WriteString(normalStyle.Render(line))
				}
			}
			content.WriteString("\n")
		}
		
		// Show scroll indicator if needed
		if len(features) > availableHeight {
			scrollInfo := fmt.Sprintf(" [%d-%d/%d]", start+1, end, len(features))
			content.WriteString(helpStyle.Render(scrollInfo))
			content.WriteString("\n")
		}
	}
	
	// Vim-style status line at bottom
	var statusBottom string
	if m.searchMode {
		statusBottom = fmt.Sprintf("/%s", m.searchInput)
		// Add cursor indicator in search
		if len(statusBottom) < m.width-1 {
			statusBottom += "_"
		}
	} else {
		statusBottom = "k/j:nav ␣:select /:search ↵:save esc:back"
	}
	
	content.WriteString("\n")
	content.WriteString(searchStyle.Render(statusBottom))
	
	return content.String()
}

// Feature Details Screen
func (m Model) updateFeatureDetailsScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.choiceCursor > 0 {
			m.choiceCursor--
		}
		return m, nil
		
	case "down", "j":
		if m.choiceCursor < len(m.choices)-1 {
			m.choiceCursor++
		}
		return m, nil
		
	case "enter":
		if len(m.choices) > 0 {
			selectedChoice := m.choices[m.choiceCursor]
			if selectedChoice == "Custom..." {
				// Switch to text input mode for custom values
				m.detailInput = ""
				m.choices = []string{} // Clear choices to enable text input
				return m, nil
			} else {
				// Use the selected choice
				m.detailInput = selectedChoice
			}
		}
		
		// Mark feature as selected and store the input
		if strings.TrimSpace(m.detailInput) != "" {
			m.selectedFeatures[m.currentFeature] = true
			m.featureValues[m.currentFeature] = strings.TrimSpace(m.detailInput)
		}
		m.screen = featureSelectionScreen
		return m, nil
		
	case "esc":
		m.screen = featureSelectionScreen
		return m, nil
		
	case "backspace":
		// Only allow backspace when in custom text input mode
		if len(m.choices) == 0 && len(m.detailInput) > 0 {
			m.detailInput = m.detailInput[:len(m.detailInput)-1]
		}
		return m, nil
		
	default:
		// Only allow typing when in custom text input mode
		if len(m.choices) == 0 && len(msg.String()) == 1 {
			m.detailInput += msg.String()
		}
		return m, nil
	}
}

func (m Model) viewFeatureDetailsScreen() string {
	var content strings.Builder
	
	// Get feature display name
	allFeatures := getAllFeatures()
	var featureName string
	for _, f := range allFeatures {
		if f.Name == m.currentFeature {
			featureName = f.Description
			break
		}
	}
	
	title := "Configure: " + featureName
	content.WriteString(headerStyle.Render(title))
	content.WriteString("\n")
	
	if len(m.choices) > 0 {
		// Show choice list
		content.WriteString(accentStyle.Render("Select an option:"))
		content.WriteString("\n\n")
		
		for i, choice := range m.choices {
			var line string
			if i == m.choiceCursor {
				line = selectedStyle.Render("> " + choice)
			} else {
				line = normalStyle.Render("  " + choice)
			}
			content.WriteString(line)
			content.WriteString("\n")
		}
		
		content.WriteString("\n")
		content.WriteString(helpStyle.Render("↑/↓ navigate | Enter select | Esc cancel"))
	} else {
		// Show text input (custom mode)
		content.WriteString(accentStyle.Render("Enter custom value:"))
		content.WriteString("\n\n")
		
		inputText := m.detailInput + "_" // Show cursor
		content.WriteString(inputStyle.Render(inputText))
		content.WriteString("\n\n")
		
		content.WriteString(helpStyle.Render("Type custom value | Enter confirm | Esc cancel"))
	}
	
	return content.String()
}

// Summary Screen  
func (m Model) updateSummaryScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "y":
		m.completed = true
		return m, tea.Quit
		
	case "n", "esc":
		m.screen = serviceNameScreen
		return m, nil
	}
	
	return m, nil
}

func (m Model) viewSummaryScreen() string {
	var content strings.Builder
	
	content.WriteString(headerStyle.Render("Configuration Summary"))
	content.WriteString("\n\n")
	
	for serviceName, service := range m.config.Services {
		content.WriteString(fmt.Sprintf("Service: %s\n", serviceName))
		content.WriteString("Features:\n")
		
		if service.Features.FirebaseAuth != nil {
			content.WriteString(fmt.Sprintf("  • Firebase Auth: %v\n", service.Features.FirebaseAuth))
		}
		if service.Features.FirebaseCloudMessagingSender {
			content.WriteString("  • Firebase Cloud Messaging Sender: enabled\n")
		}
		if service.Features.FirebaseCloudMessagingViewer {
			content.WriteString("  • Firebase Cloud Messaging Viewer: enabled\n")
		}
		if service.Features.CloudRunInvoker {
			content.WriteString("  • Cloud Run Invoker: enabled\n")
		}
		if service.Features.EventarcSubrole {
			content.WriteString("  • Eventarc Subrole: enabled\n")
		}
		if service.Features.CIDR != "" {
			content.WriteString(fmt.Sprintf("  • CIDR: %s\n", service.Features.CIDR))
		}
		if len(service.Features.BucketWriter) > 0 {
			content.WriteString(fmt.Sprintf("  • Bucket Writer: %v\n", service.Features.BucketWriter))
		}
		if len(service.Features.BucketCreator) > 0 {
			content.WriteString(fmt.Sprintf("  • Bucket Creator: %v\n", service.Features.BucketCreator))
		}
		if len(service.Features.BucketReader) > 0 {
			content.WriteString(fmt.Sprintf("  • Bucket Reader: %v\n", service.Features.BucketReader))
		}
		if len(service.Features.SubscriptionSubscriber) > 0 {
			content.WriteString(fmt.Sprintf("  • Subscription Subscriber: %v\n", service.Features.SubscriptionSubscriber))
		}
		if len(service.Features.TopicPublisher) > 0 {
			content.WriteString(fmt.Sprintf("  • Topic Publisher: %v\n", service.Features.TopicPublisher))
		}
		if service.Features.FirestoreReader {
			content.WriteString("  • Firestore Reader: enabled\n")
		}
		if service.Features.FirestoreWriter {
			content.WriteString("  • Firestore Writer: enabled\n")
		}
		if service.Features.MysqlAccess {
			content.WriteString("  • MySQL Access: enabled\n")
		}
		if service.Features.PostgresAccess {
			content.WriteString("  • PostgreSQL Access: enabled\n")
		}
		if service.Features.EnableProfiling {
			content.WriteString("  • Profiling: enabled\n")
		}
		content.WriteString("\n")
	}
	
	content.WriteString(helpStyle.Render("Generate HCL? (y/n)"))
	
	return content.String()
}

// Helper function to build features from selected options
func (m Model) buildFeatures() Features {
	// Always start with empty features - we'll build based on current selections
	features := Features{}
	
	// Firebase Auth
	if m.selectedFeatures["firebaseauth"] {
		if role, exists := m.featureValues["firebaseauth"]; exists && role != "" {
			features.FirebaseAuth = role
		} else {
			features.FirebaseAuth = "viewer" // Default value if no value stored
		}
	} else {
		features.FirebaseAuth = nil
	}
	
	// Boolean features
	features.FirebaseCloudMessagingSender = m.selectedFeatures["firebase_cloudmessaging_sender"]
	features.FirebaseCloudMessagingViewer = m.selectedFeatures["firebase_cloudmessaging_viewer"]
	features.CloudRunInvoker = m.selectedFeatures["cloudrun_invoker"]
	features.EventarcSubrole = m.selectedFeatures["eventarc_subrole"]
	features.FirestoreReader = m.selectedFeatures["firestore_reader"]
	features.FirestoreWriter = m.selectedFeatures["firestore_writer"]
	features.MysqlAccess = m.selectedFeatures["mysql_access"]
	features.PostgresAccess = m.selectedFeatures["postgres_access"]
	features.EnableProfiling = m.selectedFeatures["enable_profiling"]
	
	// CIDR
	if m.selectedFeatures["cidr"] {
		if cidr, exists := m.featureValues["cidr"]; exists && cidr != "" {
			features.CIDR = strings.TrimSpace(cidr)
		}
	}
	
	// Bucket features
	features.BucketWriter = m.buildStringList("bucket_writer", nil)
	features.BucketCreator = m.buildStringList("bucket_creator", nil)
	features.BucketReader = m.buildStringList("bucket_reader", nil)
	
	// Pub/Sub Subscription features
	features.SubscriptionSubscriber = m.buildStringList("subscription_subscriber", nil)
	features.SubscriptionViewer = m.buildStringList("subscription_viewer", nil)
	features.SubscriptionEditor = m.buildStringList("subscription_editor", nil)
	
	// Pub/Sub Topic features
	features.TopicPublisher = m.buildStringList("topic_publisher", nil)
	features.TopicViewer = m.buildStringList("topic_viewer", nil)
	features.TopicEditor = m.buildStringList("topic_editor", nil)
	
	return features
}

// Helper to build string lists for features
func (m Model) buildStringList(featureName string, existingList []string) []string {
	if m.selectedFeatures[featureName] {
		if value, exists := m.featureValues[featureName]; exists && value != "" {
			items := strings.Split(value, ",")
			for i, item := range items {
				items[i] = strings.TrimSpace(item)
			}
			return items
		}
	}
	return []string{}
}

// Setup choices for different features
func (m *Model) setupChoices(feature string) {
	switch feature {
	case "firebaseauth":
		m.choices = []string{"viewer", "admin"}
	case "cidr":
		// For CIDR, we still need text input, so provide common examples as choices
		m.choices = []string{"10.0.0.0/24", "192.168.0.0/16", "172.16.0.0/12", "Custom..."}
	case "bucket_creator", "bucket_reader", "bucket_writer":
		// Common bucket names as examples
		m.choices = []string{"data-bucket", "images-bucket", "logs-bucket", "backup-bucket", "Custom..."}
	case "subscription_subscriber", "subscription_viewer", "subscription_editor":
		// Common subscription patterns
		m.choices = []string{"user-events", "system-events", "data-updates", "notifications", "Custom..."}
	case "topic_publisher", "topic_viewer", "topic_editor":
		// Common topic patterns  
		m.choices = []string{"user-events", "system-events", "data-updates", "notifications", "Custom..."}
	default:
		m.choices = []string{}
	}
}

// Check if there are unsaved changes
func (m Model) hasUnsavedChanges() bool {
	// Check if selections have changed
	for feature, selected := range m.selectedFeatures {
		if selected != m.originalFeatures[feature] {
			return true
		}
	}
	
	// Check if any originally selected feature is now deselected
	for feature, originalSelected := range m.originalFeatures {
		if originalSelected && !m.selectedFeatures[feature] {
			return true
		}
	}
	
	// Check if values have changed for selected features
	for feature, value := range m.featureValues {
		if m.selectedFeatures[feature] && value != m.originalValues[feature] {
			return true
		}
	}
	
	// Check if originally set values are now different
	for feature, originalValue := range m.originalValues {
		if m.selectedFeatures[feature] && m.featureValues[feature] != originalValue {
			return true
		}
	}
	
	return false
}

// Unsaved Changes Screen
func (m Model) updateUnsavedChangesScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		// Discard changes - go back without saving
		if m.editingService {
			m.screen = serviceListScreen
		} else {
			m.screen = serviceNameScreen
		}
		return m, nil
		
	case "n", "N", "esc":
		// Keep editing - go back to feature selection
		m.screen = featureSelectionScreen
		return m, nil
		
	case "s", "S":
		// Save changes and exit
		service := Service{Features: m.buildFeatures()}
		m.config.Services[m.currentService] = service
		m.screen = serviceListScreen
		return m, nil
	}
	
	return m, nil
}

func (m Model) viewUnsavedChangesScreen() string {
	var content strings.Builder
	
	content.WriteString(headerStyle.Render("Unsaved Changes"))
	content.WriteString("\n\n")
	
	content.WriteString(warningStyle.Render("You have unsaved changes to this service."))
	content.WriteString("\n\n")
	
	content.WriteString(accentStyle.Render("What would you like to do?"))
	content.WriteString("\n\n")
	
	content.WriteString(errorStyle.Render("Y") + normalStyle.Render(" - Discard changes and exit"))
	content.WriteString("\n")
	content.WriteString(successStyle.Render("S") + normalStyle.Render(" - Save changes and exit"))
	content.WriteString("\n")
	content.WriteString(accentStyle.Render("N") + normalStyle.Render(" - Continue editing (ESC)"))
	content.WriteString("\n\n")
	
	content.WriteString(helpStyle.Render("Press Y to discard, S to save, or N/ESC to continue editing"))
	
	return content.String()
}