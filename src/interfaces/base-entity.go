package interfaces

import (
	"time"

	"github.com/google/uuid"
)

type BaseEntity struct {
	ID         string    `bson:"id" json:"id"`
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
	ModifiedAt time.Time `bson:"modifiedAt" json:"modifiedAt"`
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
