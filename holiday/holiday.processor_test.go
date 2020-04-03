package holiday

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncompleteNew(t *testing.T) {
	data := &Config{}
	assert.Nil(t, GetProcessor(data))
}

func TestGet(t *testing.T) {

	data := &Config{
		Holidays: dummyHolidays,
	}
	p := GetProcessor(data)
	indonesia := p.Get("ID")

	assert.Equal(t, "New Year Eve", indonesia[0].Label["EN"])
	assert.Equal(t, "1 Hari sebelum Tahun Baru", indonesia[0].Label["ID"])
	assert.Equal(t, "New Year", indonesia[1].Label["EN"])
	assert.Equal(t, "Libur tahun Baru", indonesia[1].Label["ID"])
}
