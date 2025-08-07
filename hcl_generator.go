package main

import (
	"fmt"
	"sort"
	"strings"
)

func generateHCL(config *Config) (string, error) {
	var builder strings.Builder
	
	builder.WriteString("locals {\n")
	builder.WriteString("  cloudrun_services = {\n")
	
	// Sort service names for consistent output
	serviceNames := make([]string, 0, len(config.Services))
	for name := range config.Services {
		serviceNames = append(serviceNames, name)
	}
	sort.Strings(serviceNames)
	
	for i, serviceName := range serviceNames {
		service := config.Services[serviceName]
		builder.WriteString(fmt.Sprintf("    %s = {\n", serviceName))
		
		// Generate features block
		builder.WriteString("      features = {\n")
		generateFeatures(&builder, service.Features)
		builder.WriteString("      }\n")
		
		// Generate mysql_grants if present
		if service.MysqlGrants != nil {
			builder.WriteString("      mysql_grants = {\n")
			generateMysqlGrants(&builder, service.MysqlGrants)
			builder.WriteString("      }\n")
		}
		
		// Generate postgres_grants if present
		if service.PostgresGrants != nil {
			builder.WriteString("      postgres_grants = {\n")
			generatePostgresGrants(&builder, service.PostgresGrants)
			builder.WriteString("      }\n")
		}
		
		if i < len(serviceNames)-1 {
			builder.WriteString("    }\n")
		} else {
			builder.WriteString("    }\n")
		}
	}
	
	builder.WriteString("  }\n")
	builder.WriteString("}\n")
	
	return builder.String(), nil
}

func generateFeatures(builder *strings.Builder, features Features) {
	// Handle firebaseauth
	if features.FirebaseAuth != nil {
		switch v := features.FirebaseAuth.(type) {
		case string:
			builder.WriteString(fmt.Sprintf("        firebaseauth = \"%s\"\n", v))
		case bool:
			builder.WriteString(fmt.Sprintf("        firebaseauth = %t\n", v))
		}
	}
	
	// Handle cloudrun_invoker
	if features.CloudRunInvoker {
		builder.WriteString("        cloudrun_invoker = true\n")
	}
	
	// Handle bucket_creator
	if len(features.BucketCreator) > 0 {
		builder.WriteString("        bucket_creator = [")
		for i, bucket := range features.BucketCreator {
			if i > 0 {
				builder.WriteString(",")
			}
			builder.WriteString(fmt.Sprintf("\"%s\"", bucket))
		}
		builder.WriteString("]\n")
	}
	
	// Handle firestore_access
	if features.FirestoreAccess {
		builder.WriteString("        firestore_access = true\n")
	}
	
	// Handle bucket_reader
	if len(features.BucketReader) > 0 {
		builder.WriteString("        bucket_reader = [")
		for i, bucket := range features.BucketReader {
			if i > 0 {
				builder.WriteString(",")
			}
			builder.WriteString(fmt.Sprintf("\"%s\"", bucket))
		}
		builder.WriteString("]\n")
	}
	
	// Handle bucket_writer
	if len(features.BucketWriter) > 0 {
		builder.WriteString("        bucket_writer = [")
		for i, bucket := range features.BucketWriter {
			if i > 0 {
				builder.WriteString(",")
			}
			builder.WriteString(fmt.Sprintf("\"%s\"", bucket))
		}
		builder.WriteString("]\n")
	}
	
	// Handle mysql_access
	if features.MysqlAccess {
		builder.WriteString("        mysql_access = true\n")
	}
	
	// Handle postgres_access
	if features.PostgresAccess {
		builder.WriteString("        postgres_access = true\n")
	}
}

func generateMysqlGrants(builder *strings.Builder, grants *MysqlGrants) {
	if len(grants.Read) > 0 {
		builder.WriteString("        read = [\n")
		for _, table := range grants.Read {
			builder.WriteString(fmt.Sprintf("          \"%s\",\n", table))
		}
		builder.WriteString("        ]\n")
	}
}

func generatePostgresGrants(builder *strings.Builder, grants *PostgresGrants) {
	if len(grants.Read) > 0 {
		builder.WriteString("        read = [\n")
		for _, table := range grants.Read {
			builder.WriteString(fmt.Sprintf("          \"%s\",\n", table))
		}
		builder.WriteString("        ]\n")
	}
	
	if len(grants.ReadWrite) > 0 {
		builder.WriteString("        read_write = [\n")
		for _, table := range grants.ReadWrite {
			builder.WriteString(fmt.Sprintf("          \"%s\",\n", table))
		}
		builder.WriteString("        ]\n")
	}
	
	if len(grants.ReadWriteDelete) > 0 {
		builder.WriteString("        read_write_delete = [\n")
		for _, table := range grants.ReadWriteDelete {
			builder.WriteString(fmt.Sprintf("          \"%s\",\n", table))
		}
		builder.WriteString("        ]\n")
	}
}