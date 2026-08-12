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

// GetObjectsInput structure for GetObjects. Set Prefix to get every secret
// under a path, or Keys to get only the listed ones.
type GetObjectsInput struct {
	Prefix string
	Keys   []string
}

// GetObjectOuput structure for GetObject
type GetObjectOutput struct {
	Value string
}

// GetObjectsOuput structure for GetObjects
type GetObjectsOuput struct {
	Secrets map[string]string
}

type DeleteObjectInput struct {
	Key string
}
