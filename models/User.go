package models

import (
	"time"

	"github.com/google/uuid"
)

type LevelOfCooking string

const (
	LevelNovice       LevelOfCooking = "Novice"
	LevelIntermediate LevelOfCooking = "Intermediate"
	LevelProficient   LevelOfCooking = "Proficient"
	LevelExpert       LevelOfCooking = "Expert"
)

func (l LevelOfCooking) IsValid() bool {
	switch l {
	case LevelNovice,
		LevelIntermediate,
		LevelProficient,
		LevelExpert:
		return true
	}
	return false
}

type Gender string

const (
	GenderMale   Gender = "M"
	GenderFemale Gender = "F"
	GenderOther  Gender = "O"
)

func (g Gender) IsValid() bool {
	switch g {
	case GenderMale,
		GenderFemale,
		GenderOther:
		return true
	}
	return false
}

type User struct {
	ID             uuid.UUID      `json:"ID"`
	Name           string         `json:"name"`
	Email          string         `json:"email"`
	Password       string         `json:"-"`
	Gender         Gender         `json:"gender"`
	DOB            time.Time      `json:"dob"`
	LevelOfCooking LevelOfCooking `json:"levelOfCooking"`
	CreatedAtUTC   time.Time      `json:"createdAtUTC"`
	UpdatedAtUTC   *time.Time     `json:"updatedAtUTC,omitempty"`
}

type SignUp struct {
	Name           string         `json:"name" binding:"required,min=3,max=50"`
	Email          string         `json:"email" binding:"required,email"`
	Password       string         `json:"password" binding:"required,strongPassword"`
	Gender         Gender         `json:"gender" binding:"required,gender"`
	DOB            string         `json:"dob" binding:"required,date"`
	LevelOfCooking LevelOfCooking `json:"levelOfCooking" binding:"required,levelOfCooking"`
}

type Login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}
