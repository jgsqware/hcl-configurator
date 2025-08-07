package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Service Name Screen
func (m Model) updateServiceNameScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if strings.TrimSpace(m.serviceInput) == "" {
			// If empty input and we have services, go to summary
			if len(m.config.Services) > 0 {
				m.screen = summaryScreen
				return m, nil
			}
			// Otherwise stay on service name screen
			return m, nil
		}
		
		// Set current service and go to feature selection
		m.currentService = strings.TrimSpace(m.serviceInput)
		m.serviceInput = ""
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
		// Save service and go back to service name screen
		service := Service{Features: m.buildFeatures()}
		m.config.Services[m.currentService] = service
		m.screen = serviceNameScreen
		return m, nil
		
	case "esc":
		m.screen = serviceNameScreen
		return m, nil
	}
	
	return m, nil
}

func (m Model) viewFeatureSelectionScreen() string {
	var content strings.Builder
	
	content.WriteString(headerStyle.Render(fmt.Sprintf("Configure Features for: %s", m.currentService)))
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
	features := Features{}
	
	if m.selectedFeatures["firebaseauth"] {
		role := m.detailInput
		if role == "" {
			role = "viewer"
		}
		features.FirebaseAuth = role
	}
	
	if m.selectedFeatures["cloudrun_invoker"] {
		features.CloudRunInvoker = true
	}
	
	if m.selectedFeatures["bucket_creator"] {
		if m.detailInput != "" {
			buckets := strings.Split(m.detailInput, ",")
			for i, bucket := range buckets {
				buckets[i] = strings.TrimSpace(bucket)
			}
			features.BucketCreator = buckets
		}
	}
	
	if m.selectedFeatures["firestore_access"] {
		features.FirestoreAccess = true
	}
	
	if m.selectedFeatures["bucket_reader"] {
		if m.detailInput != "" {
			buckets := strings.Split(m.detailInput, ",")
			for i, bucket := range buckets {
				buckets[i] = strings.TrimSpace(bucket)
			}
			features.BucketReader = buckets
		}
	}
	
	if m.selectedFeatures["bucket_writer"] {
		if m.detailInput != "" {
			buckets := strings.Split(m.detailInput, ",")
			for i, bucket := range buckets {
				buckets[i] = strings.TrimSpace(bucket)
			}
			features.BucketWriter = buckets
		}
	}
	
	if m.selectedFeatures["mysql_access"] {
		features.MysqlAccess = true
	}
	
	if m.selectedFeatures["postgres_access"] {
		features.PostgresAccess = true
	}
	
	return features
}