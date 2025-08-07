package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

func parseHCLFile(filename string) (*Config, error) {
	if filename == "" {
		return &Config{Services: make(map[string]Service)}, nil
	}
	
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return &Config{Services: make(map[string]Service)}, nil
	}
	
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(filename)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL: %s", diags.Error())
	}
	
	content, diags := file.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type: "locals",
			},
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to decode HCL content: %s", diags.Error())
	}
	
	config := &Config{Services: make(map[string]Service)}
	
	// Find locals block
	for _, block := range content.Blocks {
		if block.Type == "locals" {
			err := parseLocalsBlock(block.Body, config)
			if err != nil {
				return nil, err
			}
		}
	}
	
	return config, nil
}

func parseLocalsBlock(body hcl.Body, config *Config) error {
	attrs, diags := body.JustAttributes()
	if diags.HasErrors() {
		return fmt.Errorf("failed to get attributes: %s", diags.Error())
	}
	
	for name, attr := range attrs {
		if name == "cloudrun_services" {
			val, diags := attr.Expr.Value(nil)
			if diags.HasErrors() {
				return fmt.Errorf("failed to evaluate cloudrun_services: %s", diags.Error())
			}
			
			err := parseCloudRunServices(val, config)
			if err != nil {
				return err
			}
		}
	}
	
	return nil
}

func parseCloudRunServices(val cty.Value, config *Config) error {
	if !val.Type().IsObjectType() {
		return fmt.Errorf("cloudrun_services must be an object")
	}
	
	for serviceName, serviceVal := range val.AsValueMap() {
		service := Service{}
		
		if !serviceVal.Type().IsObjectType() {
			continue
		}
		
		serviceMap := serviceVal.AsValueMap()
		
		// Parse features
		if featuresVal, exists := serviceMap["features"]; exists && !featuresVal.IsNull() {
			features, err := parseFeatures(featuresVal)
			if err != nil {
				return fmt.Errorf("error parsing features for service %s: %v", serviceName, err)
			}
			service.Features = features
		}
		
		// Parse mysql_grants
		if grantsVal, exists := serviceMap["mysql_grants"]; exists && !grantsVal.IsNull() {
			grants, err := parseMysqlGrants(grantsVal)
			if err != nil {
				return fmt.Errorf("error parsing mysql_grants for service %s: %v", serviceName, err)
			}
			service.MysqlGrants = grants
		}
		
		// Parse postgres_grants
		if grantsVal, exists := serviceMap["postgres_grants"]; exists && !grantsVal.IsNull() {
			grants, err := parsePostgresGrants(grantsVal)
			if err != nil {
				return fmt.Errorf("error parsing postgres_grants for service %s: %v", serviceName, err)
			}
			service.PostgresGrants = grants
		}
		
		config.Services[serviceName] = service
	}
	
	return nil
}

func parseFeatures(val cty.Value) (Features, error) {
	features := Features{}
	
	if !val.Type().IsObjectType() {
		return features, fmt.Errorf("features must be an object")
	}
	
	featuresMap := val.AsValueMap()
	
	// Parse firebaseauth
	if v, exists := featuresMap["firebaseauth"]; exists && !v.IsNull() {
		if v.Type() == cty.String {
			features.FirebaseAuth = v.AsString()
		} else if v.Type() == cty.Bool {
			features.FirebaseAuth = v.True()
		}
	}
	
	// Parse boolean features
	parseBoolFeature := func(name string, target *bool) {
		if v, exists := featuresMap[name]; exists && !v.IsNull() && v.Type() == cty.Bool {
			*target = v.True()
		}
	}
	
	parseBoolFeature("firebase_cloudmessaging_sender", &features.FirebaseCloudMessagingSender)
	parseBoolFeature("firebase_cloudmessaging_viewer", &features.FirebaseCloudMessagingViewer)
	parseBoolFeature("cloudrun_invoker", &features.CloudRunInvoker)
	parseBoolFeature("eventarc_subrole", &features.EventarcSubrole)
	parseBoolFeature("firestore_reader", &features.FirestoreReader)
	parseBoolFeature("firestore_writer", &features.FirestoreWriter)
	parseBoolFeature("mysql_access", &features.MysqlAccess)
	parseBoolFeature("postgres_access", &features.PostgresAccess)
	parseBoolFeature("enable_profiling", &features.EnableProfiling)
	
	// Parse CIDR
	if v, exists := featuresMap["cidr"]; exists && !v.IsNull() && v.Type() == cty.String {
		features.CIDR = v.AsString()
	}
	
	// Parse string list features
	parseStringList := func(name string, target *[]string) {
		if v, exists := featuresMap[name]; exists && !v.IsNull() {
			if v.Type().IsListType() || v.Type().IsTupleType() {
				var items []string
				err := gocty.FromCtyValue(v, &items)
				if err == nil {
					*target = items
				}
			}
		}
	}
	
	parseStringList("bucket_writer", &features.BucketWriter)
	parseStringList("bucket_creator", &features.BucketCreator)
	parseStringList("bucket_reader", &features.BucketReader)
	parseStringList("subscription_subscriber", &features.SubscriptionSubscriber)
	parseStringList("subscription_viewer", &features.SubscriptionViewer)
	parseStringList("subscription_editor", &features.SubscriptionEditor)
	parseStringList("topic_publisher", &features.TopicPublisher)
	parseStringList("topic_viewer", &features.TopicViewer)
	parseStringList("topic_editor", &features.TopicEditor)
	
	return features, nil
}

func parseMysqlGrants(val cty.Value) (*MysqlGrants, error) {
	grants := &MysqlGrants{}
	
	if !val.Type().IsObjectType() {
		return grants, fmt.Errorf("mysql_grants must be an object")
	}
	
	grantsMap := val.AsValueMap()
	
	parseStringList := func(name string, target *[]string) {
		if v, exists := grantsMap[name]; exists && !v.IsNull() {
			if v.Type().IsListType() || v.Type().IsTupleType() {
				var items []string
				err := gocty.FromCtyValue(v, &items)
				if err == nil {
					*target = items
				}
			}
		}
	}
	
	parseStringList("read", &grants.Read)
	parseStringList("read_write", &grants.ReadWrite)
	parseStringList("read_write_delete", &grants.ReadWriteDelete)
	
	return grants, nil
}

func parsePostgresGrants(val cty.Value) (*PostgresGrants, error) {
	grants := &PostgresGrants{}
	
	if !val.Type().IsObjectType() {
		return grants, fmt.Errorf("postgres_grants must be an object")
	}
	
	grantsMap := val.AsValueMap()
	
	if v, exists := grantsMap["read"]; exists && !v.IsNull() {
		if v.Type().IsListType() || v.Type().IsTupleType() {
			var tables []string
			err := gocty.FromCtyValue(v, &tables)
			if err == nil {
				grants.Read = tables
			}
		}
	}
	
	if v, exists := grantsMap["read_write"]; exists && !v.IsNull() {
		if v.Type().IsListType() || v.Type().IsTupleType() {
			var tables []string
			err := gocty.FromCtyValue(v, &tables)
			if err == nil {
				grants.ReadWrite = tables
			}
		}
	}
	
	if v, exists := grantsMap["read_write_delete"]; exists && !v.IsNull() {
		if v.Type().IsListType() || v.Type().IsTupleType() {
			var tables []string
			err := gocty.FromCtyValue(v, &tables)
			if err == nil {
				grants.ReadWriteDelete = tables
			}
		}
	}
	
	return grants, nil
}