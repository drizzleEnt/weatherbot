package domain

import "errors"

var MainCityNotSetErr = errors.New("main city not set")
var CityNotFoundErr = errors.New("city not saved in db")
var CityNotGetErr = errors.New("city not ofund in API")
