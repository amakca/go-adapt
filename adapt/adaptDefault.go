package adapter

import (
	"reflect"
	"strconv"
)

// Добавить проверку на соответсвие дефолта и рестрикта: дефолтное значение должно входить в интервал

func adaptDefault(defaultValue tagValue, value reflect.Value) (err error) {
	// Проверяем на nil указатели
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			// Создаем новый указатель с дефолтным значением
			newValue := reflect.New(value.Type().Elem())
			if err := adaptDefault(defaultValue, newValue.Elem()); err != nil {
				return err
			}
			value.Set(newValue)
			return nil
		}
		// Если указатель не nil, работаем с его значением
		value = value.Elem()
	}

	if !value.IsZero() {
		return nil
	}

	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = adaptDefaultInt(defaultValue, value)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		err = adaptDefaultUint(defaultValue, value)

	case reflect.Float64, reflect.Float32:
		err = adaptDefaultFloat(defaultValue, value)

	case reflect.String:
		err = adaptDefaultString(defaultValue, value)

	default:
		err = ErrInvalidTags

	}
	return
}

func adaptDefaultInt(defaultValue tagValue, value reflect.Value) error {
	defVal, err := strconv.Atoi(string(defaultValue))
	if err != nil {
		return err
	}

	value.SetInt(int64(defVal))
	return nil
}

func adaptDefaultFloat(defaultValue tagValue, value reflect.Value) error {
	defVal, err := strconv.ParseFloat(string(defaultValue), 64)
	if err != nil {
		return err
	}

	value.SetFloat(defVal)
	return nil
}

func adaptDefaultUint(defaultValue tagValue, value reflect.Value) error {
	defVal, err := strconv.ParseUint(string(defaultValue), 10, 64)
	if err != nil {
		return err
	}

	value.SetUint(defVal)
	return nil
}

func adaptDefaultString(defaultValue tagValue, value reflect.Value) error {
	value.SetString(string(defaultValue))
	return nil
}
