package main

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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
	
	content.WriteString(headerStyle.Render("Services Configuration"))
	content.WriteString("\n\n")
	
	services := make([]string, 0, len(m.config.Services))
	for name := range m.config.Services {
		services = append(services, name)
	}
	sort.Strings(services)
	
	if len(services) == 0 {
		content.WriteString("No services configured yet.\n")
		content.WriteString(selectedStyle.Render("> Add new service"))
	} else {
		content.WriteString("Existing services:\n")
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
		
		// Add "Add new service" option
		if m.servicesCursor == len(services) {
			content.WriteString(selectedStyle.Render("> Add new service"))
		} else {
			content.WriteString(normalStyle.Render("  Add new service"))
		}
	}
	
	content.WriteString("\n\n")
	content.WriteString(helpStyle.Render("Use ↑/↓ to navigate, Enter to edit/add, d to delete, s for summary, Esc to quit"))
	
	return content.String()
}

func (m Model) countFeatures(features Features) int {
	count := 0
	if features.FirebaseAuth != nil {
		count++
	}
	if features.CloudRunInvoker {
		count++
	}
	if len(features.BucketCreator) > 0 {
		count++
	}
	if features.FirestoreAccess {
		count++
	}
	if len(features.BucketReader) > 0 {
		count++
	}
	if len(features.BucketWriter) > 0 {
		count++
	}
	if features.MysqlAccess {
		count++
	}
	if features.PostgresAccess {
		count++
	}
	return count
}

func (m *Model) loadExistingFeatures(service Service) {
	features := service.Features
	
	if features.FirebaseAuth != nil {
		m.selectedFeatures["firebaseauth"] = true
		// Store the existing value for later use in buildFeatures
	}
	if features.CloudRunInvoker {
		m.selectedFeatures["cloudrun_invoker"] = true
	}
	if len(features.BucketCreator) > 0 {
		m.selectedFeatures["bucket_creator"] = true
	}
	if features.FirestoreAccess {
		m.selectedFeatures["firestore_access"] = true
	}
	if len(features.BucketReader) > 0 {
		m.selectedFeatures["bucket_reader"] = true
	}
	if len(features.BucketWriter) > 0 {
		m.selectedFeatures["bucket_writer"] = true
	}
	if features.MysqlAccess {
		m.selectedFeatures["mysql_access"] = true
	}
	if features.PostgresAccess {
		m.selectedFeatures["postgres_access"] = true
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
		content.WriteString("Configured services:\n")
		for name := range m.config.Services {
			content.WriteString(fmt.Sprintf("  • %s\n", name))
		}
		content.WriteString("\n")
	}
	
	content.WriteString(inputStyle.Render(fmt.Sprintf("Service name: %s", m.serviceInput)))
	content.WriteString("\n")
	
	if len(m.config.Services) > 0 {
		content.WriteString(helpStyle.Render("Press Enter to add service, or press Enter with empty name to finish"))
	} else {
		content.WriteString(helpStyle.Render("Press Enter to continue"))
	}
	
	return content.String()
}

// Feature Selection Screen
func (m Model) updateFeatureSelectionScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	features := []string{
		"firebaseauth",
		"cloudrun_invoker", 
		"bucket_creator",
		"firestore_access",
		"bucket_reader",
		"bucket_writer", 
		"mysql_access",
		"postgres_access",
	}
	
	switch msg.String() {
	case "up", "k":
		if m.featuresCursor > 0 {
			m.featuresCursor--
		}
		return m, nil
		
	case "down", "j":
		if m.featuresCursor < len(features)-1 {
			m.featuresCursor++
		}
		return m, nil
		
	case " ":
		feature := features[m.featuresCursor]
		if feature == "firebaseauth" || feature == "bucket_creator" || feature == "bucket_reader" || feature == "bucket_writer" {
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
		
	case "esc":
		if m.editingService {
			m.screen = serviceListScreen
		} else {
			m.screen = serviceNameScreen
		}
		return m, nil
	}
	
	return m, nil
}

func (m Model) viewFeatureSelectionScreen() string {
	var content strings.Builder
	
	title := fmt.Sprintf("Configure Features for: %s", m.currentService)
	if m.editingService {
		title = fmt.Sprintf("Edit Features for: %s", m.currentService)
	}
	content.WriteString(headerStyle.Render(title))
	content.WriteString("\n\n")
	
	features := []string{
		"firebaseauth",
		"cloudrun_invoker",
		"bucket_creator", 
		"firestore_access",
		"bucket_reader",
		"bucket_writer",
		"mysql_access",
		"postgres_access",
	}
	
	descriptions := map[string]string{
		"firebaseauth":     "Firebase Authentication (requires role)",
		"cloudrun_invoker": "Cloud Run Invoker permissions",
		"bucket_creator":   "Storage bucket creator access (requires bucket names)",
		"firestore_access": "Firestore database access",
		"bucket_reader":    "Storage bucket reader access (requires bucket names)",
		"bucket_writer":    "Storage bucket writer access (requires bucket names)",
		"mysql_access":     "MySQL database access",
		"postgres_access":  "PostgreSQL database access",
	}
	
	for i, feature := range features {
		var line string
		
		if i == m.featuresCursor {
			line = selectedStyle.Render(fmt.Sprintf("> [ ] %s", descriptions[feature]))
		} else {
			checkbox := " "
			if m.selectedFeatures[feature] {
				checkbox = "✓"
			}
			line = normalStyle.Render(fmt.Sprintf("  [%s] %s", checkbox, descriptions[feature]))
		}
		
		if m.selectedFeatures[feature] {
			line = strings.Replace(line, "[ ]", "[✓]", 1)
			line = strings.Replace(line, "> [ ]", "> [✓]", 1)
		}
		
		content.WriteString(line)
		content.WriteString("\n")
	}
	
	content.WriteString("\n")
	content.WriteString(helpStyle.Render("Use ↑/↓ to navigate, Space to select, Enter to save, Esc to cancel"))
	
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
		prompt = "Enter Firebase Auth role (e.g., 'viewer', 'editor'):"
	case "bucket_creator", "bucket_reader", "bucket_writer":
		prompt = "Enter bucket names (comma-separated):"
	}
	
	content.WriteString(headerStyle.Render(fmt.Sprintf("Configure %s", m.currentFeature)))
	content.WriteString("\n\n")
	content.WriteString(prompt)
	content.WriteString("\n\n")
	content.WriteString(inputStyle.Render(m.detailInput))
	content.WriteString("\n")
	content.WriteString(helpStyle.Render("Press Enter to confirm, Esc to cancel"))
	
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
		if service.Features.CloudRunInvoker {
			content.WriteString("  • Cloud Run Invoker: enabled\n")
		}
		if len(service.Features.BucketCreator) > 0 {
			content.WriteString(fmt.Sprintf("  • Bucket Creator: %v\n", service.Features.BucketCreator))
		}
		if service.Features.FirestoreAccess {
			content.WriteString("  • Firestore Access: enabled\n")
		}
		if len(service.Features.BucketReader) > 0 {
			content.WriteString(fmt.Sprintf("  • Bucket Reader: %v\n", service.Features.BucketReader))
		}
		if len(service.Features.BucketWriter) > 0 {
			content.WriteString(fmt.Sprintf("  • Bucket Writer: %v\n", service.Features.BucketWriter))
		}
		if service.Features.MysqlAccess {
			content.WriteString("  • MySQL Access: enabled\n")
		}
		if service.Features.PostgresAccess {
			content.WriteString("  • PostgreSQL Access: enabled\n")
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
	
	// Apply selected features (this will override existing ones)
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
		// Feature not selected, remove it
		features.FirebaseAuth = nil
	}
	
	if m.selectedFeatures["cloudrun_invoker"] {
		features.CloudRunInvoker = true
	} else {
		features.CloudRunInvoker = false
	}
	
	if m.selectedFeatures["bucket_creator"] {
		if m.detailInput != "" {
			buckets := strings.Split(m.detailInput, ",")
			for i, bucket := range buckets {
				buckets[i] = strings.TrimSpace(bucket)
			}
			features.BucketCreator = buckets
		} else if !m.editingService {
			features.BucketCreator = []string{}
		}
	} else {
		features.BucketCreator = []string{}
	}
	
	if m.selectedFeatures["firestore_access"] {
		features.FirestoreAccess = true
	} else {
		features.FirestoreAccess = false
	}
	
	if m.selectedFeatures["bucket_reader"] {
		if m.detailInput != "" {
			buckets := strings.Split(m.detailInput, ",")
			for i, bucket := range buckets {
				buckets[i] = strings.TrimSpace(bucket)
			}
			features.BucketReader = buckets
		} else if !m.editingService {
			features.BucketReader = []string{}
		}
	} else {
		features.BucketReader = []string{}
	}
	
	if m.selectedFeatures["bucket_writer"] {
		if m.detailInput != "" {
			buckets := strings.Split(m.detailInput, ",")
			for i, bucket := range buckets {
				buckets[i] = strings.TrimSpace(bucket)
			}
			features.BucketWriter = buckets
		} else if !m.editingService {
			features.BucketWriter = []string{}
		}
	} else {
		features.BucketWriter = []string{}
	}
	
	if m.selectedFeatures["mysql_access"] {
		features.MysqlAccess = true
	} else {
		features.MysqlAccess = false
	}
	
	if m.selectedFeatures["postgres_access"] {
		features.PostgresAccess = true
	} else {
		features.PostgresAccess = false
	}
	
	return features
}