package weather

import (
	"context"
	"weatherbot/internal/clients"
	"weatherbot/internal/domain"
	"weatherbot/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type repo struct {
	db clients.DBClient
}

var _ repository.WeatherRepository = (*repo)(nil)

func NewRepository(client clients.DBClient) *repo {
	return &repo{
		db: client,
	}
}

// GetGeodata implements repository.WeatherRepository.
func (r *repo) GetGeodata(ctx context.Context, city string) (domain.CityInfo, error) {
	query := `SELECT latitude, longitude FROM geodata WHERE main_city = $1`

	var cityInfo domain.CityInfo
	err := r.db.DB().QueryRow(ctx, query, city).Scan(&cityInfo.Latitude, &cityInfo.Longitude)
	if err != nil {

		if err != pgx.ErrNoRows {
			return domain.CityInfo{}, err
		}
		return domain.CityInfo{}, errors.Wrap(domain.CityNotFoundErr, err.Error())
	}

	return cityInfo, nil
}

// GetUserCity implements repository.WeatherRepository.
func (r *repo) GetUserInfo(ctx context.Context, chatID int) (domain.UserInfo, error) {
	query := `SELECT 
	c.main_city, u.language
	FROM cities c 
	JOIN users u ON c.chat_id = u.chat_id
	WHERE c.chat_id = $1`

	var userInfo domain.UserInfo
	err := r.db.DB().QueryRow(ctx, query, chatID).
		Scan(&userInfo.City, &userInfo.Language)

	if err != nil {
		if err != pgx.ErrNoRows {
			return domain.UserInfo{}, err
		}
		return domain.UserInfo{}, errors.Wrap(domain.MainCityNotSetErr, err.Error())
	}

	return userInfo, nil
}

// SetGeodata implements repository.WeatherRepository.
func (r *repo) SetGeodata(ctx context.Context, cityInfo domain.CityInfo) error {
	query := `INSERT INTO geodata ("main_city","longitude","latitude") VALUES($1, $2, $3)`

	_, err := r.db.DB().Exec(ctx, query, cityInfo.CityName, cityInfo.Longitude, cityInfo.Latitude)
	if err != nil {
		return errors.Wrap(err, "repo.SetGeodata")
	}
	return nil
}

// SetUserInfo implements repository.WeatherRepository.
func (r *repo) SetUserInfo(ctx context.Context, userInfo domain.UserInfo, cityInfo domain.CityInfo) error {
	tx, err := r.db.DB().Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "repo.SetUserInfo")
	}
	defer tx.Rollback(ctx)

	query := `
	INSERT INTO users 
	(username, chat_id, language) 
	VALUES($1, $2, $3)
	ON CONFLICT (chat_id) 
	DO NOTHING`
	_, err = tx.Exec(ctx, query, userInfo.Username, userInfo.ChatID, userInfo.Language)
	if err != nil {
		return errors.Wrap(err, "repo.SetUserInfo")
	}

	query = `
	INSERT INTO cities 
	(chat_id, main_city) 
	VALUES ($1, $2)
	ON CONFLICT (chat_id) 
	DO UPDATE SET 
	other_cities = array_append(cities.other_cities, EXCLUDED.main_city)`

	_, err = tx.Exec(ctx, query, userInfo.ChatID, cityInfo.CityName)
	if err != nil {
		return errors.Wrap(err, "repo.SetUserInfo")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "repo.SetUserInfo")
	}

	return nil
}
