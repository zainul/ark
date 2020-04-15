package featureflag_test

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	. "github.com/zainul/ark/featureflag"
)

func TestNewError(t *testing.T) {
	config := Config{
		ConsulURL:  "error",
		HTTPClient: client,
	}
	_, err := New(config)
	assert.Error(t, err)

}
func TestNewSuccess(t *testing.T) {
	config := Config{
		ConsulURL:  "http://localhost",
		HTTPClient: client,
	}
	_, err := New(config)
	assert.Nil(t, err)

}

func TestHasAccess(t *testing.T) {

	defer gock.Off()

	gock.New("http://127.0.0.1/v1/").Get("").Reply(200).BodyString(testingFlag)

	var testingUserID int64 = 999

	config := Config{
		ConsulURL:    "http://127.0.0.1",
		HTTPClient:   client,
		Instances:    []Instance{InstanceFlight},
		CronInterval: 300,
		Prefix:       "service",
	}
	module, _ := New(config)
	assert.NotNil(t, module)

	// Denied userID
	assert.False(t, module.HasAccess("specific", InstanceFlight, 1000, false))

	// Granted userID
	assert.True(t, module.HasAccess("specific", InstanceFlight, testingUserID, false))

	// Denied QA
	assert.False(t, module.HasAccess("qa", InstanceFlight, testingUserID, false))

	// Granted QA
	assert.True(t, module.HasAccess("qa", InstanceFlight, testingUserID, true))

	// Denied all
	assert.False(t, module.HasAccess("disabled", InstanceFlight, testingUserID, true))
	assert.False(t, module.HasAccess("disabled", InstanceFlight, testingUserID, false))
	assert.False(t, module.HasAccess("disabled", InstanceFlight, 1000, true))
	assert.False(t, module.HasAccess("disabled", InstanceFlight, 1000, false))

	// Granted all
	assert.True(t, module.HasAccess("all", InstanceFlight, testingUserID, true))
	assert.True(t, module.HasAccess("all", InstanceFlight, testingUserID, false))
	assert.True(t, module.HasAccess("all", InstanceFlight, 1000, true))
	assert.True(t, module.HasAccess("all", InstanceFlight, 1000, false))

	// Deny unrecognized key
	assert.False(t, module.HasAccess("qqq", InstanceFlight, testingUserID, true))

}

func TestHasAccessInPercentage(t *testing.T) {

	defer gock.Off()

	gock.New("http://127.0.0.1/v1/").Get("").Reply(200).BodyString(testingFlag)

	config := Config{
		ConsulURL:    "http://127.0.0.1",
		HTTPClient:   client,
		Instances:    []Instance{InstanceHotel},
		CronInterval: 300,
		Prefix:       "service",
	}
	module, _ := New(config)
	assert.NotNil(t, module)

	// Granted userID
	assert.True(t, module.HasAccessInPercentageUser("general", InstanceHotel, 1000, false))
	assert.False(t, module.HasAccessInPercentageUser("general", InstanceHotel, 84848484849494, false))
	assert.True(t, module.HasAccessInPercentageUser("isqa", InstanceHotel, 999, true))
	assert.False(t, module.HasAccessInPercentageUser("disabled_all", InstanceHotel, 1000, true))
	assert.False(t, module.HasAccessInPercentageUser("disabled_all", InstanceHotel, 84848484849494, false))

	assert.True(t, module.HasAccessInPercentageUser("enable_all", InstanceHotel, 1000, false))
	assert.True(t, module.HasAccessInPercentageUser("enable_all", InstanceHotel, 84848484849494, false))
}

func TestActiveValue(t *testing.T) {
	defer gock.Off()

	gock.New("http://127.0.0.1/v1/").Get("").Reply(200).BodyString(testingFlag)

	var testingUserID int64 = 999

	config := Config{
		ConsulURL:    "http://127.0.0.1",
		HTTPClient:   client,
		Instances:    []Instance{InstanceTrain},
		CronInterval: 300,
		Prefix:       "service",
	}
	module, _ := New(config)
	assert.NotNil(t, module)

	// Denied userID
	denied, _ := module.GetActiveValue("urls", InstanceTrain, 1000, false)
	assert.False(t, false, denied)

	// Granted User
	granted, grantedValue := module.GetActiveValue("urls", InstanceTrain, testingUserID, true)
	assert.Equal(t, true, granted)
	assert.Equal(t, 1, len(grantedValue))

	// 2 active value with 1 in active value
	grantedTwoActiveValue, twoActive := module.GetActiveValue("two_value_active_one_inactive", InstanceTrain, testingUserID, true)
	assert.Equal(t, true, grantedTwoActiveValue)
	assert.Equal(t, 2, len(twoActive))

	// Disable all user with 2 active value with 1 in active value
	disableStatus, disableValue := module.GetActiveValue("disable_all", InstanceTrain, 8900, false)
	assert.Equal(t, false, disableStatus)
	assert.Equal(t, 0, len(disableValue))
}

func TestActiveValueMultiUser(t *testing.T) {
	defer gock.Off()

	gock.New("http://127.0.0.1/v1/").Get("").Reply(200).BodyString(testingFlag)

	config := Config{
		ConsulURL:    "http://127.0.0.1",
		HTTPClient:   client,
		Instances:    []Instance{InstanceTrain},
		CronInterval: 300,
		Prefix:       "service",
	}
	module, _ := New(config)
	assert.NotNil(t, module)

	// Denied userID
	denied, _ := module.GetActiveValue("specific", InstanceTrain, 1000, false)
	assert.False(t, false, denied)

	// test for 2#5000;6000;7000;
	grantedSpecificUser, grantedValueSpecificUser := module.GetActiveValue("specific", InstanceTrain, 5000, true)

	assert.Equal(t, true, grantedSpecificUser)
	assert.Equal(t, 1, len(grantedValueSpecificUser))

	// multiple user and 2 active value with 1 in active value
	grantedTwoActiveValue, twoActive := module.GetActiveValue("specific_user_with_two_value_active_one_inactive", InstanceTrain, 5000, true)
	assert.Equal(t, true, grantedTwoActiveValue)
	assert.Equal(t, 2, len(twoActive))

	// multiple user and Denied user id with 2 active value with 1 in active value
	grantedTwoActiveValue, twoActive = module.GetActiveValue("specific_user_with_two_value_active_one_inactive", InstanceTrain, 8900, false)
	assert.Equal(t, false, grantedTwoActiveValue)
	assert.Equal(t, 0, len(twoActive))

	// multiple value , one inactive and allowed user
	statusTwoWithValue, statusTwo := module.GetActiveValue("specific_user_with_two_value_active_one_inactive", InstanceTrain, 5000, false)
	assert.Equal(t, true, statusTwoWithValue)
	assert.Equal(t, 2, len(statusTwo))

	// disable user and for multi value 2 active 1 inactive
	grantedTwoActiveValue, twoActive = module.GetActiveValue("disable_all_multi_user", InstanceTrain, 8900, false)
	assert.Equal(t, false, grantedTwoActiveValue)
	assert.Equal(t, 0, len(twoActive))
}
