package env

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func EnvLoadFiles(paths ...string) {
	if len(paths) > 0 {
		err := godotenv.Overload(paths...)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func EnvLoadValues(values map[string]string) {
	if len(values) > 0 {

		for key, value := range values {
			err := os.Setenv(strings.TrimSpace(key), strings.TrimSpace(value))
			if err != nil {
				log.Fatal(err)
			}
		}

	}

}

func EnvReadFile(paths ...string) map[string]string {
	values, err := godotenv.Read(paths...)

	if err != nil {
		log.Fatal(err)
	}

	return values
}

func MustString(key string) (string, error) {
	if value, ok := getValue(key); ok {
		return value, nil
	}
	return "", fmt.Errorf("MustGetDuration value not found , key [%s]", key)
}

func StringDefault(key string, valueDefault string) string {
	value, err := MustString(key)
	if err != nil {
		return valueDefault
	}
	return value
}

func String(key string) string {
	return StringDefault(key, "")
}

func MustInt64(key string) (int64, error) {
	if value, ok := getValue(key); ok {
		return toInt64(value)
	}

	return 0, fmt.Errorf("MustGetInt64 value not found , key [%s]", key)
}

func Int64Default(key string, valueDefault int64) int64 {
	value, err := MustInt64(key)
	if err != nil {
		return valueDefault
	}
	return value
}

func Int64(key string) int64 {
	return Int64Default(key, 0)
}

func MustInt(key string) (int, error) {

	value, err := MustInt64(key)

	if err != nil {
		return 0, fmt.Errorf("MustGetInt value not found , key [%s]", key)
	}

	return int(value), nil
}

func IntDefault(key string, valueDefault int) int {
	value, err := MustInt(key)
	if err != nil {
		return valueDefault
	}
	return value
}

func Int(key string) int {
	return IntDefault(key, 0)
}

func MustBool(key string) (bool, error) {
	if value, ok := getValue(key); ok {
		return toBoolean(value)
	}

	return false, fmt.Errorf("MustGetBool value not found , key [%s]", key)
}

func BoolDefault(key string, valueDefault bool) bool {
	value, err := MustBool(key)
	if err != nil {
		return valueDefault
	}
	return value
}

func Bool(key string) bool {
	return BoolDefault(key, false)
}

func MustFloat64(key string) (float64, error) {
	if value, ok := getValue(key); ok {
		return toFloat64(value)
	}

	return 0, fmt.Errorf("MustGetFloat64 value not found , key [%s]", key)
}

func Float64Default(key string, valueDefault float64) float64 {
	value, err := MustFloat64(key)
	if err != nil {
		return valueDefault
	}
	return value
}

func Float64(key string) float64 {
	return Float64Default(key, 0)
}

func MustTimeLayout(key string, layout string) (time.Time, error) {
	if value, ok := getValue(key); ok {
		return toTimeByLayout(value, layout)
	}

	return time.Time{}, fmt.Errorf("MustGetTime value not found , key [%s]", key)
}

func TimeLayoutDefault(key string, layout string, valueDefault time.Time) time.Time {
	value, err := MustTimeLayout(key, layout)
	if err != nil {
		return valueDefault
	}
	return value
}

func TimeLayout(key string, layout string) time.Time {
	return TimeLayoutDefault(key, layout, time.Now())
}

func MustTimeSeconds(key string) (time.Time, error) {
	if value, ok := getValue(key); ok {
		return toTimeBySeconds(value)
	}

	return time.Time{}, fmt.Errorf("MustGetTime value not found , key [%s]", key)
}

func TimeSecondsDefault(key string, valueDefault time.Time) time.Time {
	value, err := MustTimeSeconds(key)
	if err != nil {
		return valueDefault
	}
	return value
}

func TimeSeconds(key string) time.Time {
	return TimeSecondsDefault(key, time.Now())
}

func MustDuration(key string) (time.Duration, error) {
	if value, ok := getValue(key); ok {
		return toDuration(value)
	}

	return time.Duration(0), fmt.Errorf("MustGetDuration value not found , key [%s]", key)
}

func DurationDefault(key string, valueDefault time.Duration) time.Duration {
	value, err := MustDuration(key)
	if err != nil {
		return valueDefault
	}
	return value
}

func Duration(key string) time.Duration {
	return DurationDefault(key, time.Duration(0))
}

func MustStrings(key string) ([]string, error) {
	if value, ok := getValue(key); ok {
		return toStrings(value), nil
	}

	return []string{}, fmt.Errorf("MustGetStrings value not found , key [%s]", key)
}

func StringsDefault(key string, valueDefault []string) []string {
	value, err := MustStrings(key)
	if err != nil {
		return valueDefault
	}
	return value
}

func Strings(key string) []string {
	return StringsDefault(key, []string{})
}

func toInt64(value string) (int64, error) {
	f, err := toFloat64(value)
	if err != nil {
		return 0, err
	}
	return int64(f), nil
}

func toStrings(value string) []string {
	return strings.Split(value, ",")
}

func toFloat64(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

func toBoolean(value string) (bool, error) {
	return strconv.ParseBool(value)
}

func toDuration(value string) (time.Duration, error) {
	return time.ParseDuration(value)
}

func toTimeBySeconds(value string) (time.Time, error) {
	v, err := toInt64(value)
	return time.Unix(v, 0), err
}

func toTimeByLayout(value string, layout string) (time.Time, error) {
	return time.Parse(layout, value)
}

func getValue(key string) (string, bool) {
	if value, isValue := os.LookupEnv(key); isValue {
		return value, true
	}
	return "", false
}
