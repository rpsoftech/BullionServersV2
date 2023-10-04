package interfaces

import (
	"time"

	"github.com/google/uuid"
)

type BaseEntity struct {
	ID         string    `bson:"id" json:"id" validate:"required,uuid"`
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt" validate:"required"`
	ModifiedAt time.Time `bson:"modifiedAt" json:"modifiedAt" validate:"required"`
}

func (b *BaseEntity) CreateNewId() (r *BaseEntity) {
	id := uuid.New().String()
	b.ID = id
	b.CreatedAt = time.Now()
	b.ModifiedAt = time.Now()
	return b
}

func (b *BaseEntity) Updated() (r *BaseEntity) {
	b.ModifiedAt = time.Now()
	return b
}
