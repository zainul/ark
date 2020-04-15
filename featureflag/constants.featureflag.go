package featureflag

// List of instances
const (
	InstanceGeneral Instance = iota
)

// List of flag statuses
const (
	StatusDisabled Status = iota
	StatusQA
	StatusSpecific
	StatusAll
	StatusInPercentageUser
)

/*
const (
	flagEnabled        = 1
	versioningRedisKey = "version-featureflag"
)*/

var (
	instanceTypeMap = map[Instance]string{
		InstanceGeneral:    "general",
	}

	// statusvalue is status in multivalue row, 1 is active and 0 is inactive
	statusValue = map[string]bool{
		"1": true,
		"0": false,
	}
)
