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
	
	// Boolean features
	if features.FirebaseCloudMessagingSender {
		builder.WriteString("        firebase_cloudmessaging_sender = true\n")
	}
	if features.FirebaseCloudMessagingViewer {
		builder.WriteString("        firebase_cloudmessaging_viewer = true\n")
	}
	if features.CloudRunInvoker {
		builder.WriteString("        cloudrun_invoker = true\n")
	}
	if features.EventarcSubrole {
		builder.WriteString("        eventarc_subrole = true\n")
	}
	if features.FirestoreReader {
		builder.WriteString("        firestore_reader = true\n")
	}
	if features.FirestoreWriter {
		builder.WriteString("        firestore_writer = true\n")
	}
	if features.MysqlAccess {
		builder.WriteString("        mysql_access = true\n")
	}
	if features.PostgresAccess {
		builder.WriteString("        postgres_access = true\n")
	}
	if features.EnableProfiling {
		builder.WriteString("        enable_profiling = true\n")
	}
	
	// Handle CIDR
	if features.CIDR != "" {
		builder.WriteString(fmt.Sprintf("        cidr = \"%s\"\n", features.CIDR))
	}
	
	// Bucket features
	generateStringList(builder, "bucket_writer", features.BucketWriter)
	generateStringList(builder, "bucket_creator", features.BucketCreator)
	generateStringList(builder, "bucket_reader", features.BucketReader)
	
	// Pub/Sub Subscription features
	generateStringList(builder, "subscription_subscriber", features.SubscriptionSubscriber)
	generateStringList(builder, "subscription_viewer", features.SubscriptionViewer)
	generateStringList(builder, "subscription_editor", features.SubscriptionEditor)
	
	// Pub/Sub Topic features
	generateStringList(builder, "topic_publisher", features.TopicPublisher)
	generateStringList(builder, "topic_viewer", features.TopicViewer)
	generateStringList(builder, "topic_editor", features.TopicEditor)
}

func generateStringList(builder *strings.Builder, name string, items []string) {
	if len(items) > 0 {
		builder.WriteString(fmt.Sprintf("        %s = [", name))
		for i, item := range items {
			if i > 0 {
				builder.WriteString(",")
			}
			builder.WriteString(fmt.Sprintf("\"%s\"", item))
		}
		builder.WriteString("]\n")
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