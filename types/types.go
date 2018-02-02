package types

// PutObjectInput structure for PutObject
type PutObjectInput struct {
	Key         string
	Value       string
	Application string
	Environment string
}

// GetObjectInput structure for GetObject
type GetObjectInput struct {
	Key     string
	Version string
}

// GetObjectsInput structure for ListObjectsInput
type GetObjectsInput struct {
	Prefix string
}

// GetObjectOuput structure for GetObject
type GetObjectOutput struct {
	Value string
}

// GetObjectsOuput structure for ListObjectsOutput
type GetObjectsOuput struct {
	Secrets map[string]string
}
