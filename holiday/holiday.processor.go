package holiday

import (
	"sync"
)

var instance *defaultProcessor
var once sync.Once

func init() {
	once = sync.Once{}
}

// GetProcessor Holiday processor
func GetProcessor(data *Config) Processor {

	if data.Holidays == nil {
		return nil
	}

	instance = &defaultProcessor{
		data: data,
	}

	return instance

}

// Get list of holiday based on country id
func (c *defaultProcessor) Get(countryID string) []Holiday {
	return c.data.Holidays[countryID]
}
