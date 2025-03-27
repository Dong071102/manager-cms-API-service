package models

import (
	"github.com/google/uuid"
	pgvector "github.com/pgvector/pgvector-go"

)

type Admin struct {
	AdminID       uuid.UUID `json:"admin_id" gorm:"type:uuid;primaryKey"`
	AdminCode     string    `json:"admin_code" gorm:"type:varchar(100);unique;not null"`
	FaceEmbedding pgvector.Vector `json:"face_embedding" gorm:"type:vector(512)"`
	User          User      `gorm:"foreignKey:AdminID;references:UserID"`
}
