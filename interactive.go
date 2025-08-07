package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func buildConfigInteractively() (*Config, error) {
	config := &Config{
		Services: make(map[string]Service),
	}
	
	scanner := bufio.NewScanner(os.Stdin)
	
	fmt.Println("=== HCL Configuration Builder ===")
	fmt.Println()
	
	for {
		serviceName, err := promptForServiceName(scanner)
		if err != nil {
			return nil, err
		}
		
		if serviceName == "" {
			break
		}
		
		features, err := promptForFeatures(scanner)
		if err != nil {
			return nil, err
		}
		
		config.Services[serviceName] = Service{
			Features: features,
		}
		
		fmt.Printf("✓ Service '%s' configured\n\n", serviceName)
	}
	
	if len(config.Services) == 0 {
		return nil, fmt.Errorf("no services configured")
	}
	
	return config, nil
}

func promptForServiceName(scanner *bufio.Scanner) (string, error) {
	fmt.Print("Enter service name (or press Enter to finish): ")
	if !scanner.Scan() {
		return "", scanner.Err()
	}
	
	serviceName := strings.TrimSpace(scanner.Text())
	return serviceName, nil
}

func promptForFeatures(scanner *bufio.Scanner) (Features, error) {
	features := Features{}
	
	fmt.Println("\nAvailable features:")
	fmt.Println("1. Firebase Auth")
	fmt.Println("2. Cloud Run Invoker")
	fmt.Println("3. Bucket Creator")
	fmt.Println("4. Firestore Access")
	fmt.Println("5. Bucket Reader")
	fmt.Println("6. Bucket Writer")
	fmt.Println("7. MySQL Access")
	fmt.Println("8. PostgreSQL Access")
	fmt.Println()
	
	for {
		fmt.Print("Select feature by number (or press Enter to finish): ")
		if !scanner.Scan() {
			return features, scanner.Err()
		}
		
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			break
		}
		
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid choice. Please enter a number.")
			continue
		}
		
		switch choice {
		case 1:
			value, err := promptForFirebaseAuth(scanner)
			if err != nil {
				return features, err
			}
			features.FirebaseAuth = value
			fmt.Println("✓ Firebase Auth configured")
			
		case 2:
			features.CloudRunInvoker = true
			fmt.Println("✓ Cloud Run Invoker enabled")
			
		case 3:
			buckets, err := promptForBucketList(scanner, "bucket creator")
			if err != nil {
				return features, err
			}
			features.BucketCreator = buckets
			fmt.Println("✓ Bucket Creator configured")
			
		case 4:
			features.FirestoreAccess = true
			fmt.Println("✓ Firestore Access enabled")
			
		case 5:
			buckets, err := promptForBucketList(scanner, "bucket reader")
			if err != nil {
				return features, err
			}
			features.BucketReader = buckets
			fmt.Println("✓ Bucket Reader configured")
			
		case 6:
			buckets, err := promptForBucketList(scanner, "bucket writer")
			if err != nil {
				return features, err
			}
			features.BucketWriter = buckets
			fmt.Println("✓ Bucket Writer configured")
			
		case 7:
			features.MysqlAccess = true
			fmt.Println("✓ MySQL Access enabled")
			
		case 8:
			features.PostgresAccess = true
			fmt.Println("✓ PostgreSQL Access enabled")
			
		default:
			fmt.Println("Invalid choice. Please select 1-8.")
		}
	}
	
	return features, nil
}

func promptForFirebaseAuth(scanner *bufio.Scanner) (string, error) {
	fmt.Print("Enter Firebase Auth role (e.g., 'viewer', 'editor'): ")
	if !scanner.Scan() {
		return "", scanner.Err()
	}
	
	role := strings.TrimSpace(scanner.Text())
	if role == "" {
		return "viewer", nil // default
	}
	
	return role, nil
}

func promptForBucketList(scanner *bufio.Scanner, bucketType string) ([]string, error) {
	fmt.Printf("Enter %s bucket names (comma-separated): ", bucketType)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}
	
	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return []string{}, nil
	}
	
	buckets := strings.Split(input, ",")
	for i, bucket := range buckets {
		buckets[i] = strings.TrimSpace(bucket)
	}
	
	return buckets, nil
}