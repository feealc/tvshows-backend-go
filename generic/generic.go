package generic

import (
	"errors"
	"strconv"
	"time"
)

func CheckParamInt(param, message string) (valueConverted int, err error) {
	if param != "" {
		valueConverted, err = strconv.Atoi(param)
		if err != nil {
			return valueConverted, errors.New(message)
		}
	}

	return valueConverted, nil
}

func GetCurrentDate() int {
	dateInt, _ := strconv.Atoi(time.Now().Format("20060102"))
	return dateInt
}
