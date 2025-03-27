package models

import (
	"github.com/google/uuid"
	pgvector "github.com/pgvector/pgvector-go"
)

type Lecturer struct {
	LecturerID    uuid.UUID       `json:"lecturer_id" gorm:"type:uuid;primaryKey"`
	LectainerCode string          `json:"lecturer_code" gorm:"type:varchar(100);unique"`
	FaceEmbedding pgvector.Vector `json:"face_embedding" gorm:"type:vector(512)"`
	User          User            `gorm:"foreignKey:LecturerID;references:UserID"`
}
