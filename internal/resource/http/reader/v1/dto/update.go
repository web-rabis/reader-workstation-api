package dto

import (
	"time"

	reader "github.com/web-rabis/reader-client/model"
)

type UpdateReaderRequest struct {
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
	ReadersHall    string `json:"readersHall,omitempty"`
	IsEmployee     bool   `json:"isEmployee,omitempty"`
}

func (r *UpdateReaderRequest) Reader(reader_ *reader.Reader) (*reader.Reader, error) {
	birthDate, err := time.Parse("2006-01-02", r.BirthDate)
	if err != nil {
		return nil, err
	}
	finalReader := &reader.Reader{
		TicketNumber:       reader_.TicketNumber,
		RFID:               reader_.RFID,
		RegistrationDate:   reader_.RegistrationDate,
		ReregistrationDate: reader_.ReregistrationDate,
		LastVisitDate:      reader_.LastVisitDate,
		DateCorrection:     reader_.DateCorrection,
		ExpirationDate:     reader_.ExpirationDate,
		Barcode:            r.BarCode,
		ReadersHall:        r.ReadersHall,
		Status:             reader_.Status,
		BranchId:           reader_.BranchId,
		Iin:                r.Iin,
		Firstname:          r.Firstname,
		Middlename:         r.Middlename,
		Lastname:           r.Lastname,
		YearBirth:          int64(birthDate.Year()),
		Address:            r.Address,
		Email:              r.Email,
		Phone:              r.Phone,
		WorkPhone:          r.WorkPhone,
		PlaceEmployment:    r.WorkPlace,
		Position:           r.Position,
		Course:             r.Course,
		Note:               reader_.Note,
		PriznakUdal:        reader_.PriznakUdal,
		Department:         reader_.Department,
		BirthDate:          &birthDate,

		TypeCard:       typeCardPtr(r.TypeCard),
		Faculty:        dictPtr(r.Faculty),
		AcademicDegree: dictPtr(r.AcademicDegree),
		Institution:    dictPtr(r.Institution),
		Education:      dictPtr(r.Education),
		Speciality:     dictPtr(r.Speciality),
		Nationality:    dictPtr(r.Nationality),
		SocialStatus:   socialStatusPtr(r.SocialStatus),
		Sex:            dictPtr(r.Sex),
		IsEmployee:     r.IsEmployee,
	}
	return finalReader, nil
}
func dictPtr(id *int64) *reader.Dictionary {
	if id == nil || *id == 0 {
		return nil
	}
	return &reader.Dictionary{Id: *id}
}

func typeCardPtr(id *int64) *reader.TypeCard {
	if id == nil || *id == 0 {
		return nil
	}
	return &reader.TypeCard{Id: *id}
}

func socialStatusPtr(id *int64) *reader.SocialStatus {
	if id == nil || *id == 0 {
		return nil
	}
	return &reader.SocialStatus{Id: *id}
}
