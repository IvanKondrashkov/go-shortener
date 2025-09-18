package config

import "time"

func applyStrIfEmpty(target *string, envValue, jsonValue string) {
	if envValue == "" && jsonValue != "" {
		*target = jsonValue
	}
}

func applyByteIfEmpty(target *[]byte, envValue, jsonValue string) {
	if envValue == "" && jsonValue != "" {
		*target = []byte(jsonValue)
	}
}

func applyDurationIfEmpty(target *time.Duration, envValue, jsonValue int) {
	if envValue == 0 && jsonValue != 0 {
		*target = time.Duration(jsonValue) * time.Second
	}
}

func applyIntIfEmpty(target *int, envValue, jsonValue int) {
	if envValue == 0 && jsonValue != 0 {
		*target = jsonValue
	}
}

func applyBollIfEmpty(target *bool, envValue, jsonValue bool) {
	if !envValue && jsonValue {
		*target = jsonValue
	}
}
