package adapter

import (
	"reflect"
	"regexp"
)

func adaptRegex(regexValue tagValue, value reflect.Value) error {
	// Проверяем на nil указатели
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil // Для nil указателей regex не применяется
		}
		// Если указатель не nil, работаем с его значением
		value = value.Elem()
	}

	if value.Kind() != reflect.String {
		return ErrInvalidTags
	}
	regex, err := regexp.Compile(string(regexValue))
	if err != nil {
		return err
	}
	value.SetString(regex.ReplaceAllString(value.String(), ""))
	return nil
}
