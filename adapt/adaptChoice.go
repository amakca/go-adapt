package adapter

import (
	"reflect"
	"strconv"
	"strings"
)

// Вынести основные проверки структурных тегов на уровень выше, добавить проверку пустых значений

func adaptChoice(set tagValue, value reflect.Value) (err error) {
	// Проверяем на nil указатели
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil // Для nil указателей choice не применяется
		}
		// Если указатель не nil, работаем с его значением
		value = value.Elem()
	}

	if value.IsZero() {
		return nil
	}

	options := strings.Split(string(set), SET_DELIMITER)

	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = adaptChoiceInt(options, value)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		err = adaptChoiceUint(options, value)

	case reflect.Float64, reflect.Float32:
		err = adaptChoiceFloat(options, value)

	case reflect.String:
		err = adaptChoiceString(options, value)

	default:
		err = ErrInvalidTags

	}
	return
}

func adaptChoiceInt(set []string, value reflect.Value) error {
	for _, option := range set {
		val, err := strconv.Atoi(option)
		if err != nil {
			return err
		}

		if int64(val) == value.Int() {
			return nil
		}
	}

	defVal, _ := strconv.Atoi(set[0])
	value.SetInt(int64(defVal))
	return nil
}

func adaptChoiceFloat(set []string, value reflect.Value) error {
	for _, option := range set {
		val, err := strconv.ParseFloat(option, 64)
		if err != nil {
			return err
		}

		if val == value.Float() {
			return nil
		}
	}

	defVal, _ := strconv.ParseFloat(set[0], 64)
	value.SetFloat(defVal)
	return nil
}

func adaptChoiceUint(set []string, value reflect.Value) error {
	for _, option := range set {
		val, err := strconv.ParseUint(option, 10, 64)
		if err != nil {
			return err
		}

		if val == value.Uint() {
			return nil
		}
	}

	defVal, _ := strconv.ParseUint(set[0], 10, 64)
	value.SetUint(defVal)
	return nil
}

func adaptChoiceString(set []string, value reflect.Value) error {
	for _, option := range set {
		if option == value.String() {
			return nil
		}
	}

	value.SetString(set[0])
	return nil
}
