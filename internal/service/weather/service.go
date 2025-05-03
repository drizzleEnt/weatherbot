package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"weatherbot/internal/domain"
	"weatherbot/internal/repository"
	"weatherbot/internal/service"

	"github.com/pkg/errors"
)

type weatherService struct {
	r repository.WeatherRepository
}

var _ service.WeatherService = (*weatherService)(nil)

func NewService(r repository.WeatherRepository) *weatherService {
	return &weatherService{
		r: r,
	}
}

// GetWeather implements service.WeatherService.
func (w *weatherService) GetWeather(ctx context.Context, userInfoIncome domain.UserInfo) (*domain.WeatherDataResponse, error) {
	cityInfo := domain.CityInfo{
		CityName: userInfoIncome.City,
	}

	if cityInfo.CityName == "" {
		userInfo, err := w.r.GetUserInfo(ctx, userInfoIncome.ChatID)
		if err != nil {
			return nil, err
		}
		cityInfo.CityName = userInfo.City
	}

	cityInfoDB, err := w.r.GetGeodata(ctx, cityInfo.CityName)
	if err == nil {
		cityInfo.Latitude = cityInfoDB.Latitude
		cityInfo.Longitude = cityInfoDB.Longitude

		weatherData, err := w.getWeatherFromAPI(strconv.FormatFloat(cityInfo.Latitude, 'f', 6, 64),
			strconv.FormatFloat(cityInfo.Longitude, 'f', 6, 64))
		if err != nil {
			return nil, err
		}

		return weatherData, nil
	}

	if errors.Is(err, domain.CityNotFoundErr) {
		cityInfoAPI, err := w.GetGeodataFromApi(cityInfo.CityName, userInfoIncome.Language)
		if err != nil {
			return nil, err
		}

		cityInfo.Latitude = cityInfoAPI.Latitude
		cityInfo.Longitude = cityInfoAPI.Longitude

		err = w.r.SetGeodata(ctx, cityInfo)
		if err != nil {
			return nil, err
		}

		err = w.r.SetUserInfo(ctx, userInfoIncome, cityInfo)
		if err != nil {
			return nil, err
		}

		weatherData, err := w.getWeatherFromAPI(strconv.FormatFloat(cityInfo.Latitude, 'f', 6, 64),
			strconv.FormatFloat(cityInfo.Longitude, 'f', 6, 64))
		if err != nil {
			return nil, err
		}

		fmt.Printf("weatherData: %v\n", weatherData)
		return weatherData, nil
	}

	return nil, nil
}

func (w *weatherService) GetGeodataFromApi(city string, language string) (domain.CityInfo, error) {
	q := url.Values{}
	q.Add("name", city)
	q.Add("count", "1")
	q.Add("language", "ru")
	q.Add("format", "json")

	body, err := doApiRequest(q, "geocoding-api.open-meteo.com", "/v1/search")
	if err != nil {
		return domain.CityInfo{}, errors.Wrap(err, "GetGeodataFromApi")
	}

	var result domain.GeodataResponse

	if err := json.Unmarshal(body, &result); err != nil {
		return domain.CityInfo{}, errors.Wrap(err, "GetGeodata from API")
	}

	if len(result.Results) == 0 {
		return domain.CityInfo{}, errors.Wrap(domain.CityNotGetErr, "api not found city")
	}

	cityInfo := domain.CityInfo{
		CityName:  result.Results[0].Name,
		Longitude: result.Results[0].Longitude,
		Latitude:  result.Results[0].Latitude,
	}

	return cityInfo, nil
}

// https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&hourly=temperature_2m&forecast_days=1
func (w *weatherService) getWeatherFromAPI(latitude string, longitude string) (*domain.WeatherDataResponse, error) {
	q := url.Values{}
	q.Add("hourly", "temperature_2m")
	q.Add("forecast_days", "1")
	q.Add("longitude", longitude)
	q.Add("latitude", latitude)

	body, err := doApiRequest(q, "api.open-meteo.com", "/v1/forecast")
	if err != nil {
		return nil, errors.Wrap(err, "getWeatherFromAPI")
	}
	var weatherData domain.WeatherDataResponse
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return nil, errors.Wrap(err, "getWeatherFromAPI.unmarshal")
	}

	return &weatherData, nil
}

func doApiRequest(q url.Values, host string, path string) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed New request")
	}

	req.URL.RawQuery = q.Encode()

	cl := http.Client{}

	resp, err := cl.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed Do request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed read response body")
	}

	return body, nil
}
