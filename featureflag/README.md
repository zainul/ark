# Feature Flag

Initialization
```golang
	// Create flag module that refresh every 10 minutes and exclusive for flight 
	config := featureflag.Config{
		ConsulURL:    "http://localhost:8500",
		HTTPClient:   httpClient,
		Instances:    []featureflag.Instance{featureflag.InstanceFlight},
		CronInterval: 600,
		Prefix:       "service/featureflag",
	}
	flag, err := featureflag.New(config)
```

**Check for access**
```golang
	isGranted := flag.HasAccess("myfeaturename", featureflag.InstanceFlight, userID, isQA)
```

**Check for access by percentage user id**
```golang
	isGranted := flag.HasAccessInPercentageUser("myfeaturename", featureflag.InstanceFlight, userID, isQA)
```


**GetActiveValue** is the method use for get the value that we store in consul, and applicable for feature flag concept
for example we want store multiple value of some URL let say `abc.com` and `abc_backup.com`
so we can switch which one is currently active `abc.com` or `abc_backup.com` and we can also change the value without restart the app
what happen if the all multi value is not active the value will be handle in application code, with the example
the `url` will be use the config url in application code usually set in consul also , but changes the value will be restart the app
the `config url` usually one of the multi value , for in the example is either `abc.com` or `abc_backup.com`.



```golang
   isGranted, activeValue := flag.GetActiveValue("featurename", featureflag.InstanceTrain, userID, isQA)
```

## How to add feature
- Add consul key to `{prefix}/{instance_name}/{feature_name}`
- Value should be only status of the flag
- For specific users, the value should be `2#{list_of_userids}`. Users are seperated by semicolon. Example : `2#5000;6000;7000`

## How to add multi value
- Same the way to add feature 
	- Add consul key to `{prefix}/{instance_name}/{feature_name}`
	- For specific users, the value should be `2#{list_of_userids}`. Users are seperated by semicolon. Example : `2#5000;6000;7000`
- But for multiple value we need added
	
	example **one active** with **enable QA only**:
	```
	1;
	0|http://10.255.13.105:18082
	1|http://10.255.13.105:18081
	0|http://10.255.13.105:18081
	```

	example **one active** with **specific user**:
	```
	2#5000;6000;7000;
	0|http://10.255.13.105:18082
	1|http://10.255.13.105:18081
	0|http://10.255.13.105:18081
	```

	example **multiple active value** with **specific user**:
	```
	2#5000;6000;7000;
	1|http://10.255.13.105:18082
	1|http://10.255.13.105:18081
	0|http://10.255.13.105:18081
	```

	example **multiple active value** with **enable QA only**:
	```
	1;
	1|http://10.255.13.105:18082
	1|http://10.255.13.105:18081
	0|http://10.255.13.105:18081
	```


## Statuses of flag
ID | Description
--- | ---
0 | Disabled for all
1 | Enabled for QA only
2 | Enabled for specific users
3 | Enabled for all
4 | Enabled for specific users by percentage (saved as active value)

## Available instances:
- general
