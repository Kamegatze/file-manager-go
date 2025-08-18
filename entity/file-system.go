package entity

import (
	"time"

	"github.com/google/uuid"
)

type FileSystem struct {
	Id        uuid.UUID `col:"id" id:"Id"`
	OwnerId   string    `col:"owner_id"`
	ParentId  string    `col:"parent_id"`
	Rights    string    `col:"rights"`
	IsFile    bool      `col:"is_file"`
	Name      string    `col:"name"`
	Path      string    `col:"path"`
	CreatedAt time.Time `col:"created_at"`
	UpdatedAt time.Time `col:"updated_at"`
	Deleted   bool      `col:"deleted"`
}
