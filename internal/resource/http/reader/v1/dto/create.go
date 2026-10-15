package dto

import (
	"errors"
	"strings"
	"time"

	reader "github.com/web-rabis/reader-client/model"
	"github.com/web-rabis/reader-workstation-api/internal/utils"
)

type CreateReaderRequest struct {
	TypeCard       *int64 `json:"typeCard,omitempty"`
	BarCode        string `json:"barCode,omitempty"`
	Iin            string `json:"iin,omitempty"`
	Lastname       string `json:"lastname,omitempty"`
	Middlename     string `json:"middlename,omitempty"`
	Firstname      string `json:"firstname,omitempty"`
	BirthDate      string `json:"birthDate,omitempty"`
	Address        string `json:"address,omitempty"`
	Phone          string `json:"phone,omitempty"`
	Email          string `json:"email,omitempty"`
	Sex            *int64 `json:"sex,omitempty"`
	Nationality    *int64 `json:"nationality,omitempty"`
	Education      *int64 `json:"education,omitempty"`
	SocialStatus   *int64 `json:"socialStatus,omitempty"`
	Speciality     *int64 `json:"speciality,omitempty"`
	Institution    *int64 `json:"institution,omitempty"`
	Faculty        *int64 `json:"faculty,omitempty"`
	Course         string `json:"course,omitempty"`
	WorkPlace      string `json:"workPlace,omitempty"`
	Position       string `json:"position,omitempty"`
	WorkPhone      string `json:"workPhone,omitempty"`
	AcademicDegree *int64 `json:"academicDegree,omitempty"`
	Photo          string `json:"photo,omitempty"`
	Status         string `json:"status,omitempty"`
}

func (r *CreateReaderRequest) Reader() (*reader.Reader, error) {
	now := time.Now()
	birthDate, err := time.Parse("2006-01-02", r.BirthDate)
	if err != nil {
		return nil, err
	}
	finalReader := &reader.Reader{
		RegistrationDate: &now,
		ReadersHall:      "",
		Status:           "",
		BranchId:         "0",
		Iin:              r.Iin,
		Firstname:        r.Firstname,
		Middlename:       r.Middlename,
		Lastname:         r.Lastname,
		YearBirth:        int64(birthDate.Year()),
		Address:          r.Address,
		Email:            r.Email,
		Phone:            r.Phone,
		WorkPhone:        r.WorkPhone,
		PlaceEmployment:  r.WorkPlace,
		Position:         r.Position,
		Course:           r.Course,
		Note:             "new-registration",
		Department:       "*****",
		BirthDate:        &birthDate,
		TypeCard:         typeCardPtr(r.TypeCard),
		Faculty:          dictPtr(r.Faculty),
		AcademicDegree:   dictPtr(r.AcademicDegree),
		Institution:      dictPtr(r.Institution),
		Education:        dictPtr(r.Education),
		Speciality:       dictPtr(r.Speciality),
		Nationality:      dictPtr(r.Nationality),
		SocialStatus:     socialStatusPtr(r.SocialStatus),
		Sex:              dictPtr(r.Sex),
	}
	return finalReader, nil
}
func (r *CreateReaderRequest) Validate() error {
	if r.Photo == "" {
		return errors.New("загрузите фото")
	}
	if r.Iin == "" {
		return errors.New("введите ИИН")
	}
	if r.BirthDate == "" {
		return errors.New("введите дату рождения")
	}
	if r.Lastname == "" {
		return errors.New("введите фамилию")
	}
	if r.Firstname == "" {
		return errors.New("введите имя")
	}
	if r.Sex == nil || *r.Sex == 0 {
		return errors.New("укажите пол")
	}
	if r.Nationality == nil || *r.Nationality == 0 {
		return errors.New("укажите национальность")
	}
	if r.Phone == "" {
		return errors.New("введите номер телефона")
	}
	if r.Email == "" {
		return errors.New("введите email")
	}
	if r.Address == "" {
		return errors.New("введите адрес")
	}
	if len(strings.Trim(r.Iin, " ")) != 12 {
		return errors.New("иин неверный")
	}
	birthDate, err := time.Parse("2006-01-02", r.BirthDate)
	if err != nil {
		return errors.New("неверный формат даты рождения")
	}
	beginIin := birthDate.Format("060102")
	if len(strings.Trim(r.Iin, " ")) != 12 {
		return errors.New("иин неверный")
	}
	if beginIin != r.Iin[:6] {
		return errors.New("иин неверный.")
	}
	if !utils.ValidateIIN(r.Iin) {
		return errors.New("иин неверный..")
	}
	return nil
}
