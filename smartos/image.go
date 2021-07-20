package smartos

import (
	"github.com/google/uuid"
)

type Image struct {
	NodeName string
	ID       *uuid.UUID `json:"uuid,omitempty"`
	Name     string     `json:"name,omitempty"`
	Version  string     `json:"version,omitempty"`
}
