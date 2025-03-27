package models

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	pgvector "github.com/pgvector/pgvector-go"
)
func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateStudentCode() string {
	return fmt.Sprintf("%08d", rand.Intn(100000000))
}

type Student struct {
	StudentID     uuid.UUID `json:"student_id" gorm:"type:uuid;primaryKey"`
	StudentCode   string    `json:"student_code" gorm:"type:varchar(100);unique;not null"`
	FaceEmbedding pgvector.Vector `json:"face_embedding" gorm:"type:vector(512)"`
	User          User      `gorm:"foreignKey:StudentID;references:UserID"`
}