package domain

type UserInfo struct {
	Username string
	ChatID   int
	City     string
	Language string
}

type CityInfo struct {
	CityName  string
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type GeodataResponse struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}
