package featureflag

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	consulapi "github.com/hashicorp/consul/api"
	c "github.com/robfig/cron"
	"github.com/zainul/ark/convert"
)

var newConsulClient = consulapi.NewClient

// New instance of featureflag
// Please save this in global variable to prevent unnessessary call
func New(config Config) (FeatureFlag, error) {

	// Remove trailing slash from prefix
	config.Prefix = filepath.Clean(config.Prefix)

	// Get a new client
	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = config.ConsulURL
	consulConfig.HttpClient = config.HTTPClient

	client, err := newConsulClient(consulConfig)
	if err != nil {
		log.Println("[ERR] New Featureflag - Failed to initialize consul client")
		return nil, err
	}

	newModule := &flagModule{
		config: config,
		kv:     client.KV(),
		flags:  sync.Map{},
	}

	// Initialize flags
	newModule._refreshData()

	// Initialize cron for updating version if the cron interval is set
	if config.CronInterval > 0 {
		newModule._initCron()
	}

	return newModule, nil
}

// GetActiveValue is get active value base on userID and roles, for example:
// if you want to use feature flag but the feature also give you the value that func can provide the requirement
// for example:
// feature flag url it will be added from all or specific user , and also the feature URL also give you the value that want to use
// like the active url you can use this function
func (f *flagModule) GetActiveValue(flagKey string, instance Instance, userID int64, isQA bool) (bool, map[int]string) {
	// Check the key and instance
	instanceName := instanceTypeMap[instance]

	if fg, ok := f.flags.Load(fmt.Sprintf("%s/%s", instanceName, flagKey)); ok {
		selectedFlag := fg.(flag)
		status, value := getSelectedFlag(selectedFlag, userID, isQA)
		return status, value
	}

	return false, nil
}

func (f *flagModule) HasAccessInPercentageUser(flagKey string, instance Instance, userID int64, isQA bool) bool {
	// Check the key and instance
	instanceName := instanceTypeMap[instance]

	if fg, ok := f.flags.Load(fmt.Sprintf("%s/%s", instanceName, flagKey)); ok {
		selectedFlag := fg.(flag)
		status, _ := getSelectedFlag(selectedFlag, userID, isQA)
		return status
	}

	return false
}

// get selected Flag is func to get the flag and return the value , the return is status of the feature flag and the value that setting in consul
func getSelectedFlag(selectedFlag flag, userID int64, isQA bool) (bool, map[int]string) {
	switch selectedFlag.Status {
	case StatusDisabled:
		// Access is disabled
		return false, nil
	case StatusQA:
		// Access is only for QA
		return isQA, selectedFlag.Value
	case StatusSpecific:
		// Access is for specific users
		for j := range selectedFlag.Users {
			if userID == selectedFlag.Users[j] {
				return true, selectedFlag.Value
			}
		}
		return false, nil
	case StatusAll:
		// Access granted for all
		return true, selectedFlag.Value
	case StatusInPercentageUser:
		for _, val := range selectedFlag.Value {
			if intVal, err := strconv.Atoi(val); err != nil {
				continue
			} else if userID%50 < int64(intVal) {
				return true, selectedFlag.Value
			}
		}

	}

	return false, nil
}

// HasAccess checking according to instance, userID and roles
func (f *flagModule) HasAccess(flagKey string, instance Instance, userID int64, isQA bool) bool {

	// Check the key and instance
	instanceName := instanceTypeMap[instance]

	if fg, ok := f.flags.Load(fmt.Sprintf("%s/%s", instanceName, flagKey)); ok {
		selectedFlag := fg.(flag)
		status, _ := getSelectedFlag(selectedFlag, userID, isQA)

		return status
	}

	// Key isn't recognized
	return false
}

// Update flags from consul
func (f *flagModule) _refreshData() {

	/*
		Value syntax:
		status#array_of_users

		e.g. : 3#1234;5566;909090
	*/
	options := &consulapi.QueryOptions{}
	result, _, _ := f.kv.List(f.config.Prefix, options)
	for _, v := range result {

		// Validate limited instances
		var key = strings.TrimPrefix(v.Key, f.config.Prefix+"/")
		for i := range f.config.Instances {
			key = ""
			instanceName := instanceTypeMap[f.config.Instances[i]]
			if strings.HasPrefix(v.Key, f.config.Prefix+"/"+instanceName) {
				key = strings.TrimPrefix(v.Key, f.config.Prefix+"/")
				break
			}
		}

		// Prepare flag
		if key != "" {

			flagValue := strings.Split(string(v.Value), "#")

			var userStr = []string{}
			if len(flagValue) > 1 {
				userStr = strings.Split(flagValue[1], ";")
			}

			// Get users
			users := []int64{}
			for r := range userStr {
				result, err := strconv.ParseInt(userStr[r], 10, 64)
				if err == nil {
					users = append(users, result)
				}
			}

			// Store to map
			ff := flag{
				Status: Status(convert.ToInt(flagValue[0])),
				Users:  users,
			}

			if len(users) > 0 {
				setValue(&ff, flagValue[1], true)
			} else {
				setValue(&ff, flagValue[0], false)
			}

			f.flags.Store(key, ff)
		}
	}

}

// set value is func to split the value that configure in consul,
// the func can handle when value has specific user to allow,
// the func has 3 parameter ,
// pointer of Flag struct use for modify the struct value, and set the featureFlag value before store in sync.Map,
// value is the flagValue in consul, it will be maybe contain value or only status,
// hasUser is flag check is has specific user or not
func setValue(ff *flag, value string, hasUser bool) {
	// we split the value using ;\n to make easy understand the status value (1,0,2 or 3), and the value
	multiValue := strings.Split(value, ";\n")

	// if has value but doesn't have user we need re-set status from array index (0) because the data after split with # in method _refreshData() will look like:
	//
	// 1;
	// 0 | samplevalue
	//
	// because doesn't have (#) separator.
	// so in this method we need set re-set the status from splitter (;\n) and get the index(0) for status
	//
	// but for specific user the data look like :
	// 2#5000;6000;7000;
	// so we not set the status any more because already corret
	if !hasUser {
		(*ff).Status = Status(convert.ToInt(multiValue[0]))
	}

	if len(multiValue) > 1 {

		// this split use for split multi value each row split with (\n) => new line
		values := strings.Split(multiValue[1], "\n")

		mapValue := make(map[int]string)

		// after got the value of multivalue, store to map int string
		for key, val := range values {

			// split the row to build the column with separator |
			cells := strings.Split(val, "|")

			// if the column not complete at least 2 column then continue
			if len(cells) <= 1 {
				continue
			}

			// check the status is in pre defined our map variable
			if _, ok := statusValue[cells[0]]; ok {

				// check f the status is convertable to int and the value is 1 store the value with key (iterator) as a map key
				// key iterator will be increment from 0, 1, 2, 3 until len(multivalues)
				// if not active or wrong status , will be not store
				// store value from column (1)
				if convert.ToInt(cells[0]) == 1 {
					mapValue[key] = cells[1]
				}
			}

		}
		// after complete iterate multivalue (row)
		// store the mapValue to pointer stored value
		// it will be change the value of pointer (*Flag.Value)
		(*ff).Value = mapValue
	}
}

func (f *flagModule) _initCron() {

	f.cron = c.New()

	// Update feature flag every X seconds
	f.cron.AddFunc(fmt.Sprintf("@every %ds", f.config.CronInterval), f._refreshData)
	f.cron.Start()
}
