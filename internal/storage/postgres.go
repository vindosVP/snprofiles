package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"github.com/vindosVP/snprofiles/internal/models"
)

type PostgresStorage struct {
	db *pgxpool.Pool
}

func NewPostgresStorage(db *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (p *PostgresStorage) CreateProfile(ctx context.Context, profile *models.Profile) (*models.Profile, error) {
	_, err := p.GetProfile(ctx, profile.UserId)
	if err != nil && !errors.Is(err, ErrProfileDoesNotExist) {
		return nil, errors.Wrap(err, "failed to update profile in database")
	}
	if err == nil {
		return nil, ErrProfileAlreadyExist
	}

	query := `INSERT INTO profiles (user_id, first_name, last_name, description, phone_number, city)
				VALUES ($1, $2, $3, $4, $5, $6)
				RETURNING user_id, first_name, last_name, description, phone_number, city, photo_uuid`
	rows, err := p.db.Query(ctx, query,
		profile.UserId,
		profile.FirstName,
		profile.LastName,
		profile.Description,
		profile.PhoneNumber,
		profile.City)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create profile in database")
	}
	cp, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.Profile])
	if err != nil {
		return nil, errors.Wrap(err, "failed to create profile in database")
	}
	return cp, nil
}

func (p *PostgresStorage) GetProfile(ctx context.Context, userId int64) (*models.Profile, error) {
	query := `SELECT user_id, first_name, last_name, description, phone_number, city, photo_uuid FROM profiles 
            	WHERE user_id = $1`
	rows, err := p.db.Query(ctx, query, userId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get profile from database")
	}
	cp, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.Profile])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProfileDoesNotExist
		}
		return nil, errors.Wrap(err, "failed to get profile from database")
	}
	return cp, nil
}

func (p *PostgresStorage) GetProfiles(ctx context.Context) ([]*models.Profile, error) {
	query := `SELECT user_id, first_name, last_name, description, phone_number, city, photo_uuid FROM profiles`
	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get profiles from database")
	}
	profiles, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.Profile])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return make([]*models.Profile, 0), nil
		}
		return nil, errors.Wrap(err, "failed to get profiles from database")
	}
	return profiles, nil
}

func (p *PostgresStorage) UpdateProfile(ctx context.Context, userId int64, profile *models.UpdateProfile) (*models.Profile, error) {
	_, err := p.GetProfile(ctx, userId)
	if err != nil {
		if errors.Is(err, ErrProfileDoesNotExist) {
			return nil, ErrProfileDoesNotExist
		}
		return nil, errors.Wrap(err, "failed to update profile in database")
	}

	query := `UPDATE profiles SET 
					first_name = $1,last_name= $2,description = $3,phone_number = $4,city = $5
				WHERE user_id = $6
				RETURNING user_id, first_name, last_name, description, phone_number, city, photo_uuid`
	rows, err := p.db.Query(ctx, query,
		profile.FirstName,
		profile.LastName,
		profile.Description,
		profile.PhoneNumber,
		profile.City,
		userId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update profile in database")
	}
	cp, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.Profile])
	if err != nil {
		return nil, errors.Wrap(err, "failed to update profile in database")
	}
	return cp, nil
}

func (p *PostgresStorage) SetProfilePhoto(ctx context.Context, userID int64, photoUUID *string) (*string, error) {
	var cPhoto *string
	query := `UPDATE profiles SET photo_uuid = $1 WHERE user_id = $2 RETURNING photo_uuid`
	err := p.db.QueryRow(ctx, query, photoUUID, userID).Scan(&cPhoto)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set profile photo in database")
	}
	return cPhoto, nil
}
