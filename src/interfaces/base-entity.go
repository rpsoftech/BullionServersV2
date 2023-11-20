package interfaces

import (
	"time"

	"github.com/google/uuid"
)

type BaseEntity struct {
	ID                 string    `bson:"id" json:"id" validate:"required,uuid"`
	CreatedAtExported  time.Time `bson:"-" json:"createdAt,omitempty" validate:"required"`
	ModifiedAtExported time.Time `bson:"-" json:"modifiedAt,omitempty" validate:"required"`
	CreatedAt          time.Time `bson:"createdAt" json:"-" validate:"required"`
	ModifiedAt         time.Time `bson:"modifiedAt" json:"-" validate:"required"`
}

// type baseEntityAlias BaseEntity
// type baseEntityWithTimes struct {
// 	CreatedAt  time.Time `json:"createdAt"`
// 	ModifiedAt time.Time `json:"modifiedAt"`
// 	*baseEntityAlias
// }

// func (base *BaseEntity) MarshalJSON() ([]byte, error) {
// 	if base.ExportCreatedAt {
// 		return json.Marshal(baseEntityWithTimes{
// 			CreatedAt:       base.CreatedAt,
// 			ModifiedAt:      base.ModifiedAt,
// 			baseEntityAlias: (*baseEntityAlias)(base),
// 		})
// 	}
// 	return json.Marshal(base)
// }

func (b *BaseEntity) AddTimeStamps() (r *BaseEntity) {
	b.CreatedAtExported = b.CreatedAt
	b.ModifiedAtExported = b.ModifiedAt
	return b
}
func (b *BaseEntity) RestoreTimeStamp() (r *BaseEntity) {
	b.CreatedAt = b.CreatedAtExported
	b.ModifiedAt = b.ModifiedAtExported
	return b
}
func (b *BaseEntity) createNewId() (r *BaseEntity) {
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
