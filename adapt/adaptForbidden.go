package adapter

import (
	"reflect"
	"strconv"
	"strings"
)

func adaptForbidden(forbiddenValue tagValue, value reflect.Value) (err error) {
	// Проверяем на nil указатели
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil // Для nil указателей forbidden не применяется
		}
		// Если указатель не nil, работаем с его значением
		value = value.Elem()
	}

	options := strings.Split(string(forbiddenValue), SET_DELIMITER)
	withDef := strings.Split(options[len(options)-1], VAL_DELIMITER)
	if len(withDef) != 2 {
		return ErrInvalidTags
	}
	options = options[:len(options)-1]
	options = append(options, withDef...)

	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = adaptForbiddenInt(options, value)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		err = adaptForbiddenUint(options, value)

	case reflect.Float32:
		err = adaptForbiddenFloat32(options, value)

	case reflect.Float64:
		err = adaptForbiddenFloat64(options, value)

	case reflect.String:
		err = adaptForbiddenString(options, value)

	default:
		err = ErrInvalidTags

	}
	return
}

func adaptForbiddenInt(forbiddenValue []string, value reflect.Value) error {
	lenght := len(forbiddenValue)
	for i := 0; i < lenght-1; i++ {

		forVal, err := strconv.Atoi(forbiddenValue[i])
		if err != nil {
			return err
		}

		if int64(forVal) == value.Int() {
			defVal, err := strconv.Atoi(forbiddenValue[lenght-1])
			if err != nil {
				return err
			}
			value.SetInt(int64(defVal))
			return nil
		}
	}

	return nil
}

func adaptForbiddenFloat64(forbiddenValue []string, value reflect.Value) error {
	lenght := len(forbiddenValue)
	for i := 0; i < lenght-1; i++ {

		forVal, err := strconv.ParseFloat(forbiddenValue[i], 64)
		if err != nil {
			return err
		}

		if forVal == value.Float() {
			defVal, err := strconv.ParseFloat(forbiddenValue[lenght-1], 64)
			if err != nil {
				return err
			}
			value.SetFloat(defVal)
			return nil
		}
	}

	return nil
}

func adaptForbiddenFloat32(forbiddenValue []string, value reflect.Value) error {
	lenght := len(forbiddenValue)
	for i := 0; i < lenght-1; i++ {

		forVal, err := strconv.ParseFloat(forbiddenValue[i], 32)
		if err != nil {
			return err
		}

		if float32(forVal) == float32(value.Float()) {
			defVal, err := strconv.ParseFloat(forbiddenValue[lenght-1], 32)
			if err != nil {
				return err
			}
			value.SetFloat(defVal)
			return nil
		}
	}

	return nil
}

func adaptForbiddenUint(forbiddenValue []string, value reflect.Value) error {
	lenght := len(forbiddenValue)
	for i := 0; i < lenght-1; i++ {

		forVal, err := strconv.ParseUint(forbiddenValue[i], 10, 64)
		if err != nil {
			return err
		}

		if forVal == value.Uint() {
			defVal, err := strconv.ParseUint(forbiddenValue[lenght-1], 10, 64)
			if err != nil {
				return err
			}
			value.SetUint(defVal)
			return nil
		}
	}

	return nil
}

func adaptForbiddenString(forbiddenValue []string, value reflect.Value) error {
	lenght := len(forbiddenValue)
	for i := 0; i < lenght-1; i++ {

		if forbiddenValue[i] == value.String() {
			value.SetString(forbiddenValue[lenght-1])
			return nil
		}
	}

	return nil
}
