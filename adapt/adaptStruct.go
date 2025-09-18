package adapter

import (
	"fmt"
	"log"
	"os"
	"reflect"
)

type adapter struct {
	logger *log.Logger
}

func NewAdapter() adapter {
	return adapter{logger: log.New(os.Stdout, "adapter ", log.LstdFlags)}
}

func (a *adapter) SetLogger(l *log.Logger) {
	a.logger = l
}

func (a *adapter) DisableLogger() {
	a.logger = nil
}

func (a *adapter) logf(format string, args ...any) {
	if a.logger == nil {
		return
	}
	a.logger.Printf(format, args...)
}

// AdaptStruct applies rules described in structure
// tags to fields of input structure.
// It takes as input pointer/value of structure, returns edited copy.
func (a *adapter) AdaptStruct(input any) (any, error) {
	inputValue := reflect.ValueOf(input)

	if reflect.Indirect(inputValue).Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	editedCopies := newStack()

	if err := a.processField(inputValue, "", editedCopies, ""); err != nil {
		return nil, err
	}

	editedCopies.applyChanges()

	return editedCopies.addrCopy.Interface(), nil
}

func (a *adapter) processField(input reflect.Value, tags reflect.StructTag, editedCopies *stackEditedCopies, path string) error {

	if !input.CanAddr() || input.Kind() == reflect.Ptr || input.Kind() == reflect.Interface {
		copyInput := makeCopy(input)

		// Проверяем, что копия не пустая (для nil указателей)
		if !copyInput.IsValid() {
			// Для nil указателей применяем теги напрямую
			if tags != "" {
				if err := a.adaptValue(input, parseStructTag(tags), path); err != nil {
					return err
				}
			}
			return nil
		}

		if !input.CanAddr() {
			editedCopies.updateCopy(copyInput)
		}

		switch input.Kind() {
		case reflect.Ptr:
			if !input.IsNil() {
				editedCopies.add(reflect.Indirect(input), copyInput)
			}

		case reflect.Interface:
			if input.Elem().Kind() == reflect.Pointer {
				editedCopies.add(input.Elem().Elem(), copyInput)
			} else {
				editedCopies.add(input, copyInput)
			}
		}

		input = copyInput
	}

	if input.Kind() == reflect.Struct {
		if err := a.processFields(input, editedCopies, path); err != nil {
			return err
		}
		return nil
	}

	if tags == "" {
		return nil
	}

	switch input.Kind() {
	case reflect.Array, reflect.Slice:
		// Обрабатываем элементы слайса
		for i := 0; i < input.Len(); i++ {

			val := input.Index(i)

			if isSimpleType(val) {
				if err := a.adaptValue(val, parseStructTag(tags), path); err != nil {
					return err
				}
			} else {
				if err := a.processField(val, tags, editedCopies, path); err != nil {
					return err
				}
			}
		}

	case reflect.Map:
		// Создаем копию карты для работы с адресуемыми значениями
		mapCopy := reflect.MakeMap(input.Type())

		for _, key := range input.MapKeys() {
			val := input.MapIndex(key)

			// Создаем копию значения
			valCopy := reflect.New(val.Type()).Elem()
			valCopy.Set(val)

			if isSimpleType(valCopy) {
				if err := a.adaptValue(valCopy, parseStructTag(tags), path); err != nil {
					return err
				}
			} else {
				if err := a.processField(valCopy, tags, editedCopies, path); err != nil {
					return err
				}
			}

			// Устанавливаем обработанное значение в копию карты
			mapCopy.SetMapIndex(key, valCopy)
		}

		// Заменяем оригинальную карту копией
		input.Set(mapCopy)

	default:
		if err := a.adaptValue(input, parseStructTag(tags), path); err != nil {
			return err
		}
	}

	return nil
}

// processFields iteratively processes fields of structure.
// If field has struct tag, field will be processed accordingly.
// If field is pointer or structure, processing will be
// recursively called for them.
func (a *adapter) processFields(input reflect.Value, editedCopies *stackEditedCopies, parentPath string) error {
	inputType := input.Type()

	for i := 0; i < input.NumField(); i++ {
		field := inputType.Field(i)
		value := input.Field(i)

		name := field.Name
		if jsonTag, ok := field.Tag.Lookup(TAG_JSON); ok && jsonTag != "" {
			comma := len(jsonTag)
			for i := 0; i < len(jsonTag); i++ {
				if jsonTag[i] == ',' {
					comma = i
					break
				}
			}
			if comma > 0 {
				tagName := jsonTag[:comma]
				if tagName != "-" && tagName != "" {
					name = tagName
				}
			}
		}
		var path string
		if parentPath == "" {
			path = name
		} else if name != "" {
			path = parentPath + "." + name
		} else {
			path = parentPath
		}

		if err := a.processField(value, field.Tag, editedCopies, path); err != nil {
			return err
		}
	}
	return nil
}

func makeCopy(inputValue reflect.Value) reflect.Value {
	if inputValue.Kind() == reflect.Interface {
		inputValue = inputValue.Elem()
	}

	// Проверяем на nil указатели
	if inputValue.Kind() == reflect.Ptr && inputValue.IsNil() {
		return reflect.Value{}
	}

	inputValue = reflect.Indirect(inputValue)

	copyInput := reflect.New(inputValue.Type()).Elem()
	copyInput.Set(inputValue)

	return copyInput
}

// adaptValue takes as input value of structure field and
// tag map for it. Processing method will be called for each tag.
func (a *adapter) adaptValue(value reflect.Value, tagsList tagsList, path string) error {
	// Проверяем на nil указатели
	if value.Kind() == reflect.Ptr && value.IsNil() {
		// Для nil указателей применяем только default тег
		if defaultTag, exists := tagsList[RST_DEFAULT]; exists {
			if err := tagsMap[RST_DEFAULT](defaultTag, value); err != nil {
				return err
			}
			if path != "" {
				a.logf("field=%q reason=%q new_value=%v", path, RST_DEFAULT, indirectInterface(value))
			}
		}
		return nil
	}

	// Apply tags in deterministic priority order
	ordered := []tagName{RST_DEFAULT, RST_MIN, RST_MAX, RST_CHOICE, RST_FORBIDDEN, RST_REGEX}
	for _, tn := range ordered {
		if tv, ok := tagsList[tn]; ok {
			before := indirectInterface(value)
			if err := tagsMap[tn](tv, value); err != nil {
				if path != "" {
					return fmt.Errorf("field %s, tag %s: %w", path, tn, err)
				}
				return err
			}
			after := indirectInterface(value)
			if path != "" && !reflect.DeepEqual(before, after) {
				a.logf("field=%q reason=%q new_value=%v", path, tn, after)
			}
		}
	}
	return nil
}

func isSimpleType(value reflect.Value) bool {
	switch value.Kind() {
	case
		reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Float64, reflect.Float32,
		reflect.String, reflect.Bool:

		return true

	default:

		return false

	}
}

func indirectInterface(v reflect.Value) any {
	if v.IsValid() && v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	if v.IsValid() && v.CanInterface() {
		return v.Interface()
	}
	return nil
}
