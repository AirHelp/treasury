package backend

import (
	"github.com/AirHelp/treasury/types"
)

type BackendAPI interface {
	PutObject(*types.PutObjectInput) error
	DeleteObject(*types.DeleteObjectInput) error
	GetObject(*types.GetObjectInput) (*types.GetObjectOutput, error)
	GetObjects(*types.GetObjectsInput) (*types.GetObjectsOuput, error)
}
