package main

type Config struct {
	Services map[string]Service `json:"services"`
}

type Service struct {
	Features     Features      `json:"features"`
	MysqlGrants  *MysqlGrants  `json:"mysql_grants,omitempty"`
	PostgresGrants *PostgresGrants `json:"postgres_grants,omitempty"`
}

type Features struct {
	FirebaseAuth                    interface{} `json:"firebaseauth,omitempty"`
	FirebaseCloudMessagingSender    bool        `json:"firebase_cloudmessaging_sender,omitempty"`
	FirebaseCloudMessagingViewer    bool        `json:"firebase_cloudmessaging_viewer,omitempty"`
	CloudRunInvoker                 bool        `json:"cloudrun_invoker,omitempty"`
	EventarcSubrole                 bool        `json:"eventarc_subrole,omitempty"`
	CIDR                           string      `json:"cidr,omitempty"`
	BucketWriter                   []string    `json:"bucket_writer,omitempty"`
	BucketCreator                  []string    `json:"bucket_creator,omitempty"`
	BucketReader                   []string    `json:"bucket_reader,omitempty"`
	SubscriptionSubscriber         []string    `json:"subscription_subscriber,omitempty"`
	SubscriptionViewer             []string    `json:"subscription_viewer,omitempty"`
	SubscriptionEditor             []string    `json:"subscription_editor,omitempty"`
	TopicPublisher                 []string    `json:"topic_publisher,omitempty"`
	TopicViewer                    []string    `json:"topic_viewer,omitempty"`
	TopicEditor                    []string    `json:"topic_editor,omitempty"`
	FirestoreReader                bool        `json:"firestore_reader,omitempty"`
	FirestoreWriter                bool        `json:"firestore_writer,omitempty"`
	MysqlAccess                    bool        `json:"mysql_access,omitempty"`
	PostgresAccess                 bool        `json:"postgres_access,omitempty"`
	EnableProfiling                bool        `json:"enable_profiling,omitempty"`
}

type MysqlGrants struct {
	Read            []string `json:"read,omitempty"`
	ReadWrite       []string `json:"read_write,omitempty"`
	ReadWriteDelete []string `json:"read_write_delete,omitempty"`
}

type PostgresGrants struct {
	Read            []string `json:"read,omitempty"`
	ReadWrite       []string `json:"read_write,omitempty"`
	ReadWriteDelete []string `json:"read_write_delete,omitempty"`
}