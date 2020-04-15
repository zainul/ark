package featureflag_test

import (
	"net/http"
	"os"
	"testing"

	. "github.com/zainul/ark/featureflag"
)

var client = &http.Client{}

const testingFlag = `
[
    {
        "LockIndex": 0,
        "Key": "service/flight/all",
        "Flags": 0,
        "Value": "Mw==",
        "CreateIndex": 21,
        "ModifyIndex": 21
    },
    {
        "LockIndex": 0,
        "Key": "service/flight/disabled",
        "Flags": 0,
        "Value": "MA==",
        "CreateIndex": 15,
        "ModifyIndex": 15
    },
    {
        "LockIndex": 0,
        "Key": "service/flight/qa",
        "Flags": 0,
        "Value": "MQ==",
        "CreateIndex": 17,
        "ModifyIndex": 17
    },
    {
        "LockIndex": 0,
        "Key": "service/flight/specific",
        "Flags": 0,
        "Value": "MiM5OTk7ODg4",
        "CreateIndex": 20,
        "ModifyIndex": 20
    },
    {
        "LockIndex": 0,
        "Key": "service/general/all",
        "Flags": 0,
        "Value": "Mw==",
        "CreateIndex": 55,
        "ModifyIndex": 55
    },
    {
        "Key": "service/train/specific",
        "CreateIndex": 23030645,
        "ModifyIndex": 23458931,
        "LockIndex": 0,
        "Flags": 0,
        "Value": "MiM1MDAwOzYwMDA7NzAwMDsKMHxodHRwOi8vMTAuMjU1LjEzLjEwNToxODA4MgoxfGh0dHA6Ly8xMC4yNTUuMTMuMTA1OjE4MDgxCjB8aHR0cDovLzEwLjI1NS4xMy4xMDU6MTgwODE=",
        "Session": ""
    },
    {
        "Key": "service/train/two_value_active_one_inactive",
        "CreateIndex": 23030645,
        "ModifyIndex": 23643828,
        "LockIndex": 0,
        "Flags": 0,
        "Value": "MTsKMXxodHRwOi8vMTAuMjU1LjEzLjEwNToxODA4MgoxfGh0dHA6Ly8xMC4yNTUuMTMuMTA1OjE4MDgxCjB8aHR0cDovLzEwLjI1NS4xMy4xMDU6MTgwODE=",
        "Session": ""
    },
    {
        "Key": "service/train/specific_user_with_two_value_active_one_inactive",
        "CreateIndex": 23030645,
        "ModifyIndex": 23644472,
        "LockIndex": 0,
        "Flags": 0,
        "Value": "MiM1MDAwOzYwMDA7NzAwMDsKMXxodHRwOi8vMTAuMjU1LjEzLjEwNToxODA4MgoxfGh0dHA6Ly8xMC4yNTUuMTMuMTA1OjE4MDgxCjB8aHR0cDovLzEwLjI1NS4xMy4xMDU6MTgwODE=",
        "Session": ""
    },
    {
        "Key": "service/train/urls",
        "CreateIndex": 23030645,
        "ModifyIndex": 23456291,
        "LockIndex": 0,
        "Flags": 0,
        "Value": "MTsKMHxodHRwOi8vMTAuMjU1LjEzLjEwNToxODA4MQoxfGh0dHA6Ly8xMC4yNTUuMTMuMTA1OjE4MDgxCjB8aHR0cDovLzEwLjI1NS4xMy4xMDU6MTgwODE=",
        "Session": ""
    },
    {
        "Key": "service/train/disable_all",
        "CreateIndex": 23030645,
        "ModifyIndex": 23652229,
        "LockIndex": 0,
        "Flags": 0,
        "Value": "MDsKMXxodHRwOi8vMTAuMjU1LjEzLjEwNToxODA4MgoxfGh0dHA6Ly8xMC4yNTUuMTMuMTA1OjE4MDgxCjB8aHR0cDovLzEwLjI1NS4xMy4xMDU6MTgwODE=",
        "Session": ""
    },
    {
        "Key": "service/train/disable_all_multi_user",
        "CreateIndex": 23030645,
        "ModifyIndex": 23652469,
        "LockIndex": 0,
        "Flags": 0,
        "Value": "MCM1MDAwOzYwMDA7NzAwMDsKMXxodHRwOi8vMTAuMjU1LjEzLjEwNToxODA4MgoxfGh0dHA6Ly8xMC4yNTUuMTMuMTA1OjE4MDgxCjB8aHR0cDovLzEwLjI1NS4xMy4xMDU6MTgwODE=",
        "Session": ""
    },
    {
        "Key": "service/hotel/general",
        "CreateIndex": 23030645,
        "ModifyIndex": 23652470,
        "LockIndex": 0,
        "Flags": 0,
        "Value": "NDsKMXwxMA==",
        "Session": ""
    },
    {
        "Key": "service/hotel/isqa",
        "CreateIndex": 23030645,
        "ModifyIndex": 23652471,
        "LockIndex": 0,
        "Flags": 0,
        "Value": "MTsKMXwxMA==",
        "Session": ""
    },
    {
        "Key": "service/hotel/disabled_all",
        "CreateIndex": 23030645,
        "ModifyIndex": 23652472,
        "LockIndex": 0,
        "Flags": 0,
        "Value": "MDsKMXwxMA==",
        "Session": ""
    },
    {
        "Key": "service/hotel/enable_all",
        "CreateIndex": 23030645,
        "ModifyIndex": 23652473,
        "LockIndex": 0,
        "Flags": 0,
        "Value": "MzsKMXwxMA==",
        "Session": ""
    }
]
`

func TestMain(m *testing.M) {

	TFuncPatch()

	os.Exit(m.Run())
}
