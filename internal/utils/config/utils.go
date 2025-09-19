package config

import "time"

func ApplyJSONStrIfEmpty(target *string, envValue, jsonValue string) {
	if envValue == "" && jsonValue != "" {
		*target = jsonValue
	}
}

func ApplyEnvStrIfEmpty(target *string, envValue string) {
	if envValue != "" {
		*target = envValue
	}
}

func ApplyJSONByteIfEmpty(target *[]byte, envValue, jsonValue string) {
	if envValue == "" && jsonValue != "" {
		*target = []byte(jsonValue)
	}
}

func ApplyEnvByteIfEmpty(target *[]byte, envValue string) {
	if envValue != "" {
		*target = []byte(envValue)
	}
}

func ApplyJSONDurationIfEmpty(target *time.Duration, envValue, jsonValue int) {
	if envValue == 0 && jsonValue != 0 {
		*target = time.Duration(jsonValue) * time.Second
	}
}

func ApplyEnvDurationIfEmpty(target *time.Duration, envValue int) {
	if envValue != 0 {
		*target = time.Duration(envValue) * time.Second
	}
}

func ApplyJSONIntIfEmpty(target *int, envValue, jsonValue int) {
	if envValue == 0 && jsonValue != 0 {
		*target = jsonValue
	}
}

func ApplyEnvIntIfEmpty(target *int, envValue int) {
	if envValue != 0 {
		*target = envValue
	}
}

func ApplyJSONBollIfEmpty(target *bool, envValue, jsonValue bool) {
	if !envValue && jsonValue {
		*target = jsonValue
	}
}

func ApplyEnvBollIfEmpty(target *bool, envValue bool) {
	if envValue {
		*target = true
	}
}
