package weather

import (
	"context"
	"fmt"
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
func (r *repo) GetGeodata() {
	panic("unimplemented")
}

// GetUserCity implements repository.WeatherRepository.
func (r *repo) GetUserCity(ctx context.Context, chatID int) (string, error) {
	query := `SELECT main_city, other_cities FROM cities WHERE chat_id = $1`

	var mainCity string
	err := r.db.DB().QueryRow(ctx, query, chatID).
		Scan(&mainCity)

	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		if err != pgx.ErrNoRows {
			return "", err
		}
		return "", errors.Wrap(domain.MainCityNotSetErr, err.Error())
	}

	return mainCity, nil
}
