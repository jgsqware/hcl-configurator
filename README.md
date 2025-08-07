# HCL Generator for Cloud Run Services

A vibrant, modern TUI application built with Bubble Tea that generates HCL configuration files for Terraform with interactive Cloud Run service configuration and fuzzy search.

## ✨ Features

- **🎨 Vibrant UI**: Bright, modern color scheme with excellent readability
- **🔍 Fuzzy Search**: Lightning-fast feature discovery with intelligent matching
- **🎯 Interactive Configuration**: Intuitive checkbox-style feature selection
- **📦 Multi-service Support**: Manage multiple Cloud Run services in one session
- **🔄 Live Editing**: Read, modify, and update existing HCL configurations
- **💡 Context-aware Input**: Smart prompts for different feature types
- **📊 Real-time Stats**: Feature counts and search result indicators
- **🎭 Rich Visual Feedback**: Icons, colors, and styled components throughout

## Supported Features

- **Firebase Auth**: Role-based authentication (viewer, editor, etc.)
- **Cloud Run Invoker**: Service-to-service invocation permissions  
- **Storage Buckets**: Creator, reader, and writer access with bucket selection
- **Firestore Access**: NoSQL database permissions
- **MySQL Access**: Relational database access
- **PostgreSQL Access**: Advanced relational database access

## Usage

### Build and run:
```bash
go build -o envhcl

# Create new configuration
./envhcl -output terraform.hcl

# Edit existing configuration
./envhcl -input existing.hcl -output terraform.hcl

# Edit in-place (reads and writes to same file)
./envhcl -output terraform.hcl
```

### Command line options:
- `-input`: Input HCL file to read and edit (optional)
- `-output`: Output HCL file path (default: "terraform.hcl")

### Smart Input Detection:
If no `-input` is specified but the output file exists, the application will automatically read from the output file for editing.

## Interactive Workflow

### New Configuration:
1. **Service Name Input**: Enter the name of your Cloud Run service
2. **Feature Selection**: Choose features with visual checkboxes
3. **Add More Services**: Repeat for additional services
4. **Configuration Summary**: Review and generate HCL

### Edit Existing Configuration:
1. **Service List**: Browse existing services with feature counts
2. **Edit/Add/Delete**: Modify services or add new ones
3. **Live Editing**: Changes are applied to existing configuration
4. **Summary & Save**: Review changes and update the HCL file

## 🎮 Navigation Controls

### Service List Screen:
- **↑/↓ or j/k**: Navigate through services
- **Enter**: Edit selected service or add new service  
- **d**: Delete selected service (🗑️)
- **s**: Go to summary/save screen (📊)
- **Esc**: Quit application

### Feature Selection:
- **🔍 / (slash)**: Enter fuzzy search mode
- **↑/↓ or j/k**: Navigate through features
- **Space**: Select/deselect features (☑️/☐)
- **Enter**: Save service configuration
- **Esc**: Exit search mode or cancel

### Fuzzy Search Mode:
- **Type**: Search features by name or description
- **Enter**: Exit search mode
- **Esc**: Cancel search and show all features
- **Backspace**: Delete search characters

### General:
- **Ctrl+C or q**: Quit application at any time

## Example Session Flow

### 🚀 Service List (vibrant styling):
```
🚀 Cloud Run Services Configuration

📦 2 services configured:

  ⚡ service-a (8 features)
► 🔧 service-b (3 features)  
  ✨ Add new service

Navigation: ↑/↓ navigate, Enter edit/add, d delete, s summary, Esc quit
```

### 🔍 Feature Selection with Search:
```
Edit Features for: mobile-interface

🔍 Search features: fire
📦 3/20 features selected | 🔍 2 matches

☑ ► firebaseauth - Firebase Authentication (viewer/admin)
☐   firestore_reader - Firestore database reader

Navigation: ↑/↓ navigate, Space select, / search, Enter save, Esc cancel
```

### 💡 Search Examples:
- Type `fire` → finds Firebase Auth, Firestore features
- Type `bucket` → finds all bucket-related features  
- Type `pub` → finds Pub/Sub topics and subscriptions
- Type `sql` → finds MySQL and PostgreSQL features

## Output Format

Generates clean HCL with the standard Terraform locals structure:

```hcl
locals {
  cloudrun_services = {
    service-a = {
      features = {
        firebaseauth = "viewer"
        cloudrun_invoker = true
        bucket_creator = ["bucket-a"]
        firestore_access = true
      }
    }
    service-b = {
      features = {
        firestore_access = true
        bucket_reader = ["input-data"]
        mysql_access = true
      }
    }
  }
}
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Modern TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling and colors
- [HashiCorp HCL](https://github.com/hashicorp/hcl) - HCL parsing and generation
- [Fuzzy](https://github.com/sahilm/fuzzy) - Fast fuzzy string matching

## 🎨 Design Philosophy

This application follows modern design principles:

- **High Contrast**: Bright colors on light backgrounds for excellent readability
- **Semantic Colors**: Green for success, orange for warnings, red for errors, purple for accents
- **Visual Hierarchy**: Icons, bold text, and spacing to guide the eye
- **Immediate Feedback**: Real-time search results and feature counters
- **Accessibility**: Clear visual distinctions and intuitive navigation