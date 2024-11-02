package models

import profilesv1 "github.com/vindosVP/snprofiles/gen/go"

type Profile struct {
	UserId      int64   `json:"userId" db:"user_id"`
	FirstName   *string `json:"firstName" db:"first_name"`
	LastName    *string `json:"lastName" db:"last_name"`
	Description *string `json:"description" db:"description"`
	PhoneNumber *string `json:"phoneNumber" db:"phone_number"`
	City        *string `json:"city" db:"city"`
	PhotoUUID   *string `json:"photoUUID" db:"photo_uuid"`
}

type UpdateProfile struct {
	FirstName   *string `json:"firstName" db:"first_name"`
	LastName    *string `json:"lastName" db:"last_name"`
	Description *string `json:"description" db:"description"`
	PhoneNumber *string `json:"phoneNumber" db:"phone_number"`
	City        *string `json:"city" db:"city"`
}

func (p *Profile) ToGRPC() *profilesv1.Profile {
	return &profilesv1.Profile{
		UserId:      p.UserId,
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		Description: p.Description,
		PhoneNumber: p.PhoneNumber,
		City:        p.City,
		PhotoUUID:   p.PhotoUUID,
	}
}

func UpdateProfileFromGRPC(in *profilesv1.PutProfile) *UpdateProfile {
	return &UpdateProfile{
		FirstName:   in.FirstName,
		LastName:    in.LastName,
		Description: in.Description,
		PhoneNumber: in.PhoneNumber,
		City:        in.City,
	}
}

func ProfileFromGRPC(in *profilesv1.Profile) *Profile {
	return &Profile{
		UserId:      in.GetUserId(),
		FirstName:   in.FirstName,
		LastName:    in.LastName,
		Description: in.Description,
		PhoneNumber: in.PhoneNumber,
		City:        in.City,
		PhotoUUID:   in.PhotoUUID,
	}
}
