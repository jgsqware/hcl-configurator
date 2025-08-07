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
	FirebaseAuth    interface{} `json:"firebaseauth,omitempty"`
	CloudRunInvoker bool        `json:"cloudrun_invoker,omitempty"`
	BucketCreator   []string    `json:"bucket_creator,omitempty"`
	FirestoreAccess bool        `json:"firestore_access,omitempty"`
	BucketReader    []string    `json:"bucket_reader,omitempty"`
	BucketWriter    []string    `json:"bucket_writer,omitempty"`
	MysqlAccess     bool        `json:"mysql_access,omitempty"`
	PostgresAccess  bool        `json:"postgres_access,omitempty"`
}

type MysqlGrants struct {
	Read []string `json:"read,omitempty"`
}

type PostgresGrants struct {
	Read            []string `json:"read,omitempty"`
	ReadWrite       []string `json:"read_write,omitempty"`
	ReadWriteDelete []string `json:"read_write_delete,omitempty"`
}