package dto

type WeatherDTO struct {
	Temperature Current `json:"current"`
}

type Current struct {
	Celcius    float32 `json:"temp_c"`
	Fahrenheit float32 `json:"temp_f"`
}
