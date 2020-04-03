package holiday

import (
	"time"
)

type defaultProcessor struct {
	data *Config
}

type Holiday struct {
	Date  time.Time         `json:"date"`
	Label map[string]string `json:"label"` // Label for localization (EN / ID for example)
}

// Config for Holiday Processor
// Format Value for HolidayList: CountryID(Key): []Holiday (Value) List of Holiday for the selected Country
// For example map["ID"] will return slice of pointer of holiday for Indonesia Country
type HolidayList map[string][]Holiday

type Config struct {
	Holidays HolidayList
}

type Processor interface {
	Get(countryID string) []Holiday
}
