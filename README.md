# HCL Generator for Cloud Run Services

A modern TUI application built with Bubble Tea that generates HCL configuration files for Terraform with interactive Cloud Run service configuration.

## Features

- **Interactive Terminal UI**: Beautiful, intuitive interface powered by Bubble Tea
- **Multi-service configuration**: Add multiple Cloud Run services in one session
- **Feature selection**: Easy checkbox-style feature selection with descriptions
- **Dynamic input**: Context-sensitive prompts for feature details
- **Real-time preview**: Configuration summary before generation
- **Professional styling**: Clean, modern terminal interface with colors and formatting

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

## Navigation Controls

### Service List Screen:
- **↑/↓ or j/k**: Navigate through services
- **Enter**: Edit selected service or add new service
- **d**: Delete selected service
- **s**: Go to summary/save screen
- **Esc**: Quit application

### Feature Selection:
- **↑/↓ or j/k**: Navigate through features
- **Space**: Select/deselect features
- **Enter**: Save service configuration
- **Esc**: Cancel and return to service list

### General:
- **Ctrl+C or q**: Quit application at any time

## Example Session Flow

### Service List (when editing existing file):
```
┌─ HCL Configuration Builder ─┐
│                             │
│ Services Configuration      │
│                             │
│ Existing services:          │
│ > mobile-interface (3 features) │
│   datalog-business (2 features) │
│   Add new service           │
│                             │
│ ↑/↓ navigate, Enter edit/add, d delete, s summary │
└─────────────────────────────┘
```

### Feature Selection:
```
┌─ Edit Features for: mobile-interface ─┐
│                                       │
│ > [✓] Firebase Auth (requires role)   │
│   [✓] Cloud Run Invoker permissions   │
│   [✓] Storage bucket creator access   │
│   [ ] Firestore database access       │
│   [ ] Storage bucket reader access    │
│   [ ] Storage bucket writer access    │
│   [ ] MySQL database access           │
│   [ ] PostgreSQL database access      │
│                                       │
│ ↑/↓ navigate, Space select, Enter save │
└───────────────────────────────────────┘
```

## Output Format

Generates clean HCL with the standard Terraform locals structure:

```hcl
locals {
  cloudrun_services = {
    mobile-api = {
      features = {
        firebaseauth = "viewer"
        cloudrun_invoker = true
        bucket_creator = ["datalog-dmz"]
        firestore_access = true
      }
    }
    data-processor = {
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

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling