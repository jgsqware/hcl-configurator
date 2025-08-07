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
./envhcl -output terraform.hcl
```

### Command line options:
- `-output`: Output HCL file path (default: "terraform.hcl")

## Interactive Workflow

1. **Service Name Input**: Enter the name of your Cloud Run service
2. **Feature Selection**: 
   - Navigate with ↑/↓ arrow keys
   - Select features with Space
   - Some features prompt for additional details (Firebase roles, bucket names)
3. **Add More Services**: Repeat for additional services  
4. **Configuration Summary**: Review your complete configuration
5. **Generate HCL**: Confirm to create the Terraform file

## Navigation Controls

- **↑/↓ or j/k**: Navigate through options
- **Space**: Select/deselect features
- **Enter**: Confirm input or proceed to next screen
- **Esc**: Go back to previous screen
- **Ctrl+C or q**: Quit application

## Example Session Flow

```
┌─ HCL Configuration Builder ─┐
│                             │
│ Enter Service Name          │ 
│                             │
│ Service name: [mobile-api_] │
│                             │
│ Press Enter to continue     │
└─────────────────────────────┘
```

```
┌─ Configure Features for: mobile-api ─┐
│                                      │
│ > [✓] Firebase Auth (requires role)  │
│   [ ] Cloud Run Invoker permissions  │
│   [✓] Storage bucket creator access  │
│   [✓] Firestore database access      │
│   [ ] Storage bucket reader access   │
│   [ ] Storage bucket writer access   │
│   [ ] MySQL database access          │
│   [ ] PostgreSQL database access     │
│                                      │
│ ↑/↓ navigate, Space select, Enter save │
└──────────────────────────────────────┘
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