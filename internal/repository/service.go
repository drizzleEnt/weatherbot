package repository

import (
	"context"
	"weatherbot/internal/domain"
)

type WeatherRepository interface {
	GetGeodata(ctx context.Context, city string) (domain.CityInfo, error)
	GetUserInfo(ctx context.Context, chatID int) (domain.UserInfo, error)
	SetGeodata(ctx context.Context, cityInfo domain.CityInfo) error
	SetUserInfo(ctx context.Context, userInfo domain.UserInfo, cityInfo domain.CityInfo) error
}
