package adapt

import (
	"reflect"
	"strconv"
)

func adaptMin(minValue tagValue, value reflect.Value) (err error) {
	// Проверяем на nil указатели
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil // Для nil указателей min не применяется
		}
		// Если указатель не nil, работаем с его значением
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = adaptMinInt(minValue, value)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		err = adaptMinUint(minValue, value)

	case reflect.Float64:
		err = adaptMinFloat64(minValue, value)
	case reflect.Float32:
		err = adaptMinFloat32(minValue, value)

	default:
		err = ErrInvalidTags

	}
	return
}

func adaptMinInt(minValue tagValue, value reflect.Value) error {
	min, err := strconv.Atoi(string(minValue))
	if err != nil {
		return err
	}

	if value.Int() < int64(min) {
		value.SetInt(int64(min))
	}
	return nil
}

func adaptMinFloat64(minValue tagValue, value reflect.Value) error {
	min, err := strconv.ParseFloat(string(minValue), 64)
	if err != nil {
		return err
	}

	if value.Float() < min {
		value.SetFloat(min)
	}
	return nil
}

func adaptMinFloat32(minValue tagValue, value reflect.Value) error {
	min64, err := strconv.ParseFloat(string(minValue), 32)
	if err != nil {
		return err
	}
	min := float32(min64)
	if float32(value.Float()) < min {
		value.SetFloat(float64(min))
	}
	return nil
}

func adaptMinUint(minValue tagValue, value reflect.Value) error {
	min, err := strconv.ParseUint(string(minValue), 10, 64)
	if err != nil {
		return err
	}

	if value.Uint() < min {
		value.SetUint(min)
	}
	return nil
}
