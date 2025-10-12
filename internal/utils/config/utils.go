package config

import "time"

// ApplyJSONStrIfEmpty хелпер функция для JSON конфига
func ApplyJSONStrIfEmpty(target *string, envValue, jsonValue string) {
	if envValue == "" && jsonValue != "" {
		*target = jsonValue
	}
}

// ApplyEnvStrIfEmpty хелпер функция для Env конфига
func ApplyEnvStrIfEmpty(target *string, envValue string) {
	if envValue != "" {
		*target = envValue
	}
}

// ApplyJSONByteIfEmpty хелпер функция для JSON конфига
func ApplyJSONByteIfEmpty(target *[]byte, envValue, jsonValue string) {
	if envValue == "" && jsonValue != "" {
		*target = []byte(jsonValue)
	}
}

// ApplyEnvByteIfEmpty хелпер функция для Env конфига
func ApplyEnvByteIfEmpty(target *[]byte, envValue string) {
	if envValue != "" {
		*target = []byte(envValue)
	}
}

// ApplyJSONDurationIfEmpty хелпер функция для JSON конфига
func ApplyJSONDurationIfEmpty(target *time.Duration, envValue, jsonValue int) {
	if envValue == 0 && jsonValue != 0 {
		*target = time.Duration(jsonValue) * time.Second
	}
}

// ApplyEnvDurationIfEmpty хелпер функция для Env конфига
func ApplyEnvDurationIfEmpty(target *time.Duration, envValue int) {
	if envValue != 0 {
		*target = time.Duration(envValue) * time.Second
	}
}

// ApplyJSONIntIfEmpty хелпер функция для JSON конфига
func ApplyJSONIntIfEmpty(target *int, envValue, jsonValue int) {
	if envValue == 0 && jsonValue != 0 {
		*target = jsonValue
	}
}

// ApplyEnvIntIfEmpty хелпер функция для Env конфига
func ApplyEnvIntIfEmpty(target *int, envValue int) {
	if envValue != 0 {
		*target = envValue
	}
}

// ApplyJSONBollIfEmpty хелпер функция для JSON конфига
func ApplyJSONBollIfEmpty(target *bool, envValue, jsonValue bool) {
	if !envValue && jsonValue {
		*target = jsonValue
	}
}

// ApplyEnvBollIfEmpty хелпер функция для Env конфига
func ApplyEnvBollIfEmpty(target *bool, envValue bool) {
	if envValue {
		*target = true
	}
}
