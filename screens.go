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
			m.loadExistingFeatures(m.config.Services[serviceName])
			m.screen = featureSelectionScreen
			return m, nil
		} else {
			// Add new service
			m.editingService = false
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
		if len(m.config.Services) == 0 {
			return m, tea.Quit
		}
		return m, nil
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
	
	if features.FirebaseAuth != nil {
		m.selectedFeatures["firebaseauth"] = true
	}
	if features.FirebaseCloudMessagingSender {
		m.selectedFeatures["firebase_cloudmessaging_sender"] = true
	}
	if features.FirebaseCloudMessagingViewer {
		m.selectedFeatures["firebase_cloudmessaging_viewer"] = true
	}
	if features.CloudRunInvoker {
		m.selectedFeatures["cloudrun_invoker"] = true
	}
	if features.EventarcSubrole {
		m.selectedFeatures["eventarc_subrole"] = true
	}
	if features.CIDR != "" {
		m.selectedFeatures["cidr"] = true
	}
	if len(features.BucketWriter) > 0 {
		m.selectedFeatures["bucket_writer"] = true
	}
	if len(features.BucketCreator) > 0 {
		m.selectedFeatures["bucket_creator"] = true
	}
	if len(features.BucketReader) > 0 {
		m.selectedFeatures["bucket_reader"] = true
	}
	if len(features.SubscriptionSubscriber) > 0 {
		m.selectedFeatures["subscription_subscriber"] = true
	}
	if len(features.SubscriptionViewer) > 0 {
		m.selectedFeatures["subscription_viewer"] = true
	}
	if len(features.SubscriptionEditor) > 0 {
		m.selectedFeatures["subscription_editor"] = true
	}
	if len(features.TopicPublisher) > 0 {
		m.selectedFeatures["topic_publisher"] = true
	}
	if len(features.TopicViewer) > 0 {
		m.selectedFeatures["topic_viewer"] = true
	}
	if len(features.TopicEditor) > 0 {
		m.selectedFeatures["topic_editor"] = true
	}
	if features.FirestoreReader {
		m.selectedFeatures["firestore_reader"] = true
	}
	if features.FirestoreWriter {
		m.selectedFeatures["firestore_writer"] = true
	}
	if features.MysqlAccess {
		m.selectedFeatures["mysql_access"] = true
	}
	if features.PostgresAccess {
		m.selectedFeatures["postgres_access"] = true
	}
	if features.EnableProfiling {
		m.selectedFeatures["enable_profiling"] = true
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

// Get all features with descriptions
func getAllFeatures() []FeatureItem {
	return []FeatureItem{
		{"firebaseauth", "Firebase Authentication (viewer/admin)"},
		{"firebase_cloudmessaging_sender", "Firebase Cloud Messaging sender"},
		{"firebase_cloudmessaging_viewer", "Firebase Cloud Messaging viewer"},
		{"cloudrun_invoker", "Cloud Run Invoker permissions"},
		{"eventarc_subrole", "Eventarc subscription role"},
		{"cidr", "CIDR block access (requires CIDR)"},
		{"bucket_writer", "Storage bucket writer (requires bucket names)"},
		{"bucket_creator", "Storage bucket creator (requires bucket names)"},
		{"bucket_reader", "Storage bucket reader (requires bucket names)"},
		{"subscription_subscriber", "Pub/Sub subscription subscriber (requires names)"},
		{"subscription_viewer", "Pub/Sub subscription viewer (requires names)"},
		{"subscription_editor", "Pub/Sub subscription editor (requires names)"},
		{"topic_publisher", "Pub/Sub topic publisher (requires names)"},
		{"topic_viewer", "Pub/Sub topic viewer (requires names)"},
		{"topic_editor", "Pub/Sub topic editor (requires names)"},
		{"firestore_reader", "Firestore database reader"},
		{"firestore_writer", "Firestore database writer"},
		{"mysql_access", "MySQL database access"},
		{"postgres_access", "PostgreSQL database access"},
		{"enable_profiling", "Enable application profiling"},
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
	
	// If not in search mode and no filtered features, use all features
	if !m.searchMode && len(m.filteredFeatures) == 0 {
		m.filteredFeatures = make([]string, len(allFeatures))
		for i, f := range allFeatures {
			m.filteredFeatures[i] = f.Name
		}
	}
	
	switch msg.String() {
	case "/":
		// Enter search mode
		m.searchMode = true
		m.searchInput = ""
		return m, nil
		
	case "esc":
		if m.searchMode {
			// Exit search mode
			m.searchMode = false
			m.searchInput = ""
			m.filteredFeatures = make([]string, len(allFeatures))
			for i, f := range allFeatures {
				m.filteredFeatures[i] = f.Name
			}
			m.featuresCursor = 0
			return m, nil
		} else {
			// Exit feature selection
			if m.editingService {
				m.screen = serviceListScreen
			} else {
				m.screen = serviceNameScreen
			}
			return m, nil
		}
	}
	
	if m.searchMode {
		switch msg.String() {
		case "enter":
			// Exit search mode
			m.searchMode = false
			return m, nil
			
		case "backspace":
			if len(m.searchInput) > 0 {
				m.searchInput = m.searchInput[:len(m.searchInput)-1]
				m.updateFilteredFeatures()
				if m.featuresCursor >= len(m.filteredFeatures) {
					m.featuresCursor = len(m.filteredFeatures) - 1
				}
				if m.featuresCursor < 0 {
					m.featuresCursor = 0
				}
			}
			return m, nil
			
		default:
			if len(msg.String()) == 1 {
				m.searchInput += msg.String()
				m.updateFilteredFeatures()
				m.featuresCursor = 0 // Reset cursor to top after search
			}
			return m, nil
		}
	}
	
	// Normal navigation mode
	switch msg.String() {
	case "up", "k":
		if m.featuresCursor > 0 {
			m.featuresCursor--
		}
		return m, nil
		
	case "down", "j":
		if m.featuresCursor < len(m.filteredFeatures)-1 {
			m.featuresCursor++
		}
		return m, nil
		
	case " ":
		if len(m.filteredFeatures) == 0 {
			return m, nil
		}
		
		feature := m.filteredFeatures[m.featuresCursor]
		// Features that need additional input
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
			// These features need additional input
			m.currentFeature = feature
			m.detailInput = ""
			m.screen = featureDetailsScreen
			return m, nil
		} else {
			// Toggle boolean features
			m.selectedFeatures[feature] = !m.selectedFeatures[feature]
		}
		return m, nil
		
	case "enter":
		// Save service and go back to service list screen
		service := Service{Features: m.buildFeatures()}
		m.config.Services[m.currentService] = service
		m.screen = serviceListScreen
		return m, nil
	}
	
	return m, nil
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
	var content strings.Builder
	
	// Clean title without complex styling
	var title string
	if m.editingService {
		title = "Edit Features for: " + m.currentService
	} else {
		title = "Configure Features for: " + m.currentService
	}
	content.WriteString(headerStyle.Render(title))
	content.WriteString("\n")
	
	// Search input if in search mode
	if m.searchMode {
		searchPrompt := "Search features: " + m.searchInput
		content.WriteString(searchStyle.Render(searchPrompt))
		content.WriteString("\n")
	}
	
	// Feature count and search info
	allFeatures := getAllFeatures()
	descriptions := make(map[string]string)
	for _, f := range allFeatures {
		descriptions[f.Name] = f.Description
	}
	
	selectedCount := 0
	for _, selected := range m.selectedFeatures {
		if selected {
			selectedCount++
		}
	}
	
	statusLine := fmt.Sprintf("%d/%d features selected", selectedCount, len(allFeatures))
	if m.searchMode {
		statusLine += fmt.Sprintf(" | %d matches", len(m.filteredFeatures))
	}
	content.WriteString(accentStyle.Render(statusLine))
	content.WriteString("\n\n")
	
	// Features list (filtered or all)
	features := m.filteredFeatures
	if len(features) == 0 {
		features = make([]string, len(allFeatures))
		for i, f := range allFeatures {
			features[i] = f.Name
		}
	}
	
	if len(features) == 0 {
		content.WriteString(warningStyle.Render("No features match your search"))
		content.WriteString("\n")
	} else {
		for i, feature := range features {
			isSelected := m.selectedFeatures[feature]
			
			// Create a clean line without mixing styles
			var prefix string
			var featureLine string
			
			if i == m.featuresCursor {
				// Selected line
				if isSelected {
					prefix = "[x]"
				} else {
					prefix = "[ ]"
				}
				featureLine = selectedStyle.Render(prefix + " " + feature + " - " + descriptions[feature])
			} else {
				// Normal line  
				if isSelected {
					prefix = successStyle.Render("[x]")
				} else {
					prefix = "[.]"
				}
				featureLine = normalStyle.Render(prefix + " " + feature) + helpStyle.Render(" - " + descriptions[feature])
			}
			
			content.WriteString(featureLine)
			content.WriteString("\n")
		}
	}
	
	content.WriteString("\n")
	
	// Clean help text
	var helpText string
	if m.searchMode {
		helpText = "Search Mode: Type to search, Enter to exit search, Esc to cancel search"
	} else {
		helpText = "Navigation: ↑/↓ navigate, Space select, / search, Enter save, Esc cancel"
	}
	content.WriteString(helpStyle.Render(helpText))
	
	return content.String()
}

// Feature Details Screen
func (m Model) updateFeatureDetailsScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		input := strings.TrimSpace(m.detailInput)
		if input != "" {
			m.selectedFeatures[m.currentFeature] = true
			// Store the detail input (we'll handle this in buildFeatures)
			if m.currentFeature == "firebaseauth" {
				m.selectedFeatures["firebaseauth_value"] = true
			}
		}
		m.screen = featureSelectionScreen
		return m, nil
		
	case "esc":
		m.screen = featureSelectionScreen
		return m, nil
		
	case "backspace":
		if len(m.detailInput) > 0 {
			m.detailInput = m.detailInput[:len(m.detailInput)-1]
		}
		return m, nil
		
	default:
		if len(msg.String()) == 1 {
			m.detailInput += msg.String()
		}
		return m, nil
	}
}

func (m Model) viewFeatureDetailsScreen() string {
	var content strings.Builder
	
	var prompt string
	switch m.currentFeature {
	case "firebaseauth":
		prompt = "Enter Firebase Auth role ('viewer' or 'admin'):"
	case "cidr":
		prompt = "Enter CIDR block (e.g., '10.0.0.0/24'):"
	case "bucket_creator", "bucket_reader", "bucket_writer":
		prompt = "Enter bucket names (comma-separated):"
	case "subscription_subscriber", "subscription_viewer", "subscription_editor":
		prompt = "Enter Pub/Sub subscription names (comma-separated):"
	case "topic_publisher", "topic_viewer", "topic_editor":
		prompt = "Enter Pub/Sub topic names (comma-separated):"
	}
	
	title := "Configure " + m.currentFeature
	content.WriteString(headerStyle.Render(title))
	content.WriteString("\n\n")
	
	content.WriteString(accentStyle.Render(prompt))
	content.WriteString("\n\n")
	
	inputText := m.detailInput
	content.WriteString(inputStyle.Render(inputText))
	content.WriteString("\n")
	
	helpText := "Enter to confirm | Esc to cancel"
	content.WriteString(helpStyle.Render(helpText))
	
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
	// Start with existing features if editing
	features := Features{}
	if m.editingService {
		if existingService, exists := m.config.Services[m.currentService]; exists {
			features = existingService.Features
		}
	}
	
	// Firebase Auth
	if m.selectedFeatures["firebaseauth"] {
		role := m.detailInput
		if role == "" {
			if m.editingService && features.FirebaseAuth != nil {
				// Keep existing value
			} else {
				role = "viewer"
			}
		}
		if role != "" {
			features.FirebaseAuth = role
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
		if m.detailInput != "" {
			features.CIDR = strings.TrimSpace(m.detailInput)
		} else if m.editingService {
			// Keep existing value
		} else {
			features.CIDR = ""
		}
	} else {
		features.CIDR = ""
	}
	
	// Bucket features
	features.BucketWriter = m.buildStringList("bucket_writer", features.BucketWriter)
	features.BucketCreator = m.buildStringList("bucket_creator", features.BucketCreator)
	features.BucketReader = m.buildStringList("bucket_reader", features.BucketReader)
	
	// Pub/Sub Subscription features
	features.SubscriptionSubscriber = m.buildStringList("subscription_subscriber", features.SubscriptionSubscriber)
	features.SubscriptionViewer = m.buildStringList("subscription_viewer", features.SubscriptionViewer)
	features.SubscriptionEditor = m.buildStringList("subscription_editor", features.SubscriptionEditor)
	
	// Pub/Sub Topic features
	features.TopicPublisher = m.buildStringList("topic_publisher", features.TopicPublisher)
	features.TopicViewer = m.buildStringList("topic_viewer", features.TopicViewer)
	features.TopicEditor = m.buildStringList("topic_editor", features.TopicEditor)
	
	return features
}

// Helper to build string lists for features
func (m Model) buildStringList(featureName string, existingList []string) []string {
	if m.selectedFeatures[featureName] {
		if m.detailInput != "" {
			items := strings.Split(m.detailInput, ",")
			for i, item := range items {
				items[i] = strings.TrimSpace(item)
			}
			return items
		} else if m.editingService {
			// Keep existing value
			return existingList
		}
	}
	return []string{}
}