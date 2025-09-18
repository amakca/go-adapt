package adapt

import (
	"reflect"
	"strconv"
)

func adaptMax(maxValue tagValue, value reflect.Value) (err error) {
	// Проверяем на nil указатели
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil // Для nil указателей max не применяется
		}
		// Если указатель не nil, работаем с его значением
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = adaptMaxInt(maxValue, value)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		err = adaptMaxUint(maxValue, value)

	case reflect.Float64:
		err = adaptMaxFloat64(maxValue, value)
	case reflect.Float32:
		err = adaptMaxFloat32(maxValue, value)

	default:
		err = ErrInvalidTags

	}
	return
}

func adaptMaxInt(maxValue tagValue, value reflect.Value) error {
	max, err := strconv.Atoi(string(maxValue))
	if err != nil {
		return err
	}

	if value.Int() > int64(max) {
		value.SetInt(int64(max))
	}
	return nil
}

func adaptMaxFloat64(maxValue tagValue, value reflect.Value) error {
	max, err := strconv.ParseFloat(string(maxValue), 64)
	if err != nil {
		return err
	}

	if value.Float() > max {
		value.SetFloat(max)
	}
	return nil
}

func adaptMaxFloat32(maxValue tagValue, value reflect.Value) error {
	max64, err := strconv.ParseFloat(string(maxValue), 32)
	if err != nil {
		return err
	}
	max := float32(max64)
	if float32(value.Float()) > max {
		value.SetFloat(float64(max))
	}
	return nil
}

func adaptMaxUint(maxValue tagValue, value reflect.Value) error {
	max, err := strconv.ParseUint(string(maxValue), 10, 64)
	if err != nil {
		return err
	}

	if value.Uint() > max {
		value.SetUint(max)
	}
	return nil
}
