package holiday

import (
	"os"
	"testing"
	"time"
)

var dummyHolidays = HolidayList{}

func TestMain(m *testing.M) {

	initDummyData()

	os.Exit(m.Run())
}

func initDummyData() {
	loc, _ := time.LoadLocation("UTC")

	holidays := []Holiday{
		Holiday{
			Date: time.Date(2018, time.December, 31, 0, 0, 0, 0, loc),
			Label: map[string]string{
				"ID": "1 Hari sebelum Tahun Baru",
				"EN": "New Year Eve",
			},
		},
		Holiday{
			Date: time.Date(2019, time.January, 1, 0, 0, 0, 0, loc),
			Label: map[string]string{
				"ID": "Libur tahun Baru",
				"EN": "New Year",
			},
		},
	}

	dummyHolidays["ID"] = holidays
}
