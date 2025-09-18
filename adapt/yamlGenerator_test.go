package adapter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestStruct для тестирования генерации YAML
type TestStruct struct {
	IntForMax       int            `json:"int-for-max" rst-min:"5" info:"Счетчик"`
	FloatForDefault float64        `json:"float-for-default" rst-default:"3.14" info:"Число Пи"`
	StringForChoice string         `json:"string-for-choice" rst-choice:"apple||banana||orange" info:"Фрукт"`
	UintForbidden   uint           `json:"uint-forbidden" rst-forbidden:"1||2||3**10" info:"Идентификатор"`
	StringRegex     string         `json:"string-regex" rst-regex:"[^a-zA-Z]+" info:"Текст без букв"`
	CombinedField   int            `json:"combined-field" rst-min:"5" rst-max:"100" rst-default:"50" info:"Комбинированное поле"`
	SliceField      []int          `json:"slice-field" rst-max:"5" info:"Список чисел"`
	MapField        map[string]int `json:"map-field" rst-min:"3" info:"Карта значений"`
	PtrField        *string        `json:"ptr-field" rst-choice:"yes||no" info:"Указатель на строку"`
}

// NestedStruct для тестирования вложенных структур
type NestedStruct struct {
	Name   string `json:"name" rst-regex:"[a-zA-Z]+" info:"Имя пользователя"`
	Age    int    `rst-min:"18" rst-max:"120" info:"Возраст"`
	Email  string `rst-regex:"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$" info:"Email адрес"`
	Active bool   `rst-default:"true" info:"Активен ли пользователь"`
}

// ComplexStruct для тестирования сложных структур
type ComplexStruct struct {
	User      NestedStruct           `info:"Информация о пользователе"`
	Users     []NestedStruct         `rst-max:"10" info:"Список пользователей"`
	Settings  map[string]interface{} `info:"Настройки"`
	Count     int                    `rst-min:"0" rst-max:"1000" rst-default:"0" info:"Счетчик"`
	Status    string                 `rst-choice:"active||inactive||pending" info:"Статус"`
	Forbidden float64                `rst-forbidden:"0.0||-1.0**-10.0" info:"Запрещенные значения"`
}

func Test_GenerateStructYAML(t *testing.T) {
	ptrString := "yes"

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "Simple struct with RST tags",
			input: TestStruct{
				IntForMax:       10,
				FloatForDefault: 3.14,
				StringForChoice: "apple",
				UintForbidden:   5,
				StringRegex:     "123",
				CombinedField:   75,
				SliceField:      []int{1, 2, 3},
				MapField:        map[string]int{"a": 1, "b": 2},
				PtrField:        &ptrString,
			},
			expected: `# Generated YAML structure with RST tags comments

# Счетчик; минимальное значение - 5
  int-for-max: 10
# Число Пи; значение по умолчанию - 3.14
  float-for-default: 3.14
# Фрукт; допустимые значения: apple, banana, orange
  string-for-choice: "apple"
# Идентификатор; запрещенные значения: 1, 2, 3, подменное значение: 10
  uint-forbidden: 5
# Текст без букв; регулярное выражение: [^a-zA-Z]+
  string-regex: "123"
# Комбинированное поле; минимальное значение - 5; максимальное значение - 100; значение по умолчанию - 50
  combined-field: 75
# Список чисел; максимальное значение - 5
  slice-field:
    - 1
    - 2
    - 3
# Карта значений; минимальное значение - 3
  map-field:
    a: 1
    b: 2
# Указатель на строку; допустимые значения: yes, no
  ptr-field: "yes"
`,
		},
		{
			name: "Complex nested struct",
			input: ComplexStruct{
				User: NestedStruct{
					Name:   "John",
					Age:    25,
					Email:  "john@example.com",
					Active: true,
				},
				Users: []NestedStruct{
					{Name: "Alice", Age: 30, Email: "alice@example.com", Active: true},
					{Name: "Bob", Age: 35, Email: "bob@example.com", Active: false},
				},
				Settings: map[string]interface{}{
					"theme": "dark",
					"lang":  "en",
				},
				Count:     100,
				Status:    "active",
				Forbidden: 5.5,
			},
			expected: `# Generated YAML structure with RST tags comments

# Информация о пользователе
  user:
# Имя пользователя; регулярное выражение: [a-zA-Z]+
    name: "John"
# Возраст; минимальное значение - 18; максимальное значение - 120
    age: 25
# Email адрес; регулярное выражение: ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
    email: "john@example.com"
# Активен ли пользователь; значение по умолчанию - true
    active: true
# Список пользователей; максимальное значение - 10
  users:
    - 
# Имя пользователя; регулярное выражение: [a-zA-Z]+
        name: "Alice"
# Возраст; минимальное значение - 18; максимальное значение - 120
        age: 30
# Email адрес; регулярное выражение: ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
        email: "alice@example.com"
# Активен ли пользователь; значение по умолчанию - true
        active: true
    - 
# Имя пользователя; регулярное выражение: [a-zA-Z]+
        name: "Bob"
# Возраст; минимальное значение - 18; максимальное значение - 120
        age: 35
# Email адрес; регулярное выражение: ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
        email: "bob@example.com"
# Активен ли пользователь; значение по умолчанию - true
        active: false
# Настройки
  settings:
    lang: "en"
    theme: "dark"
# Счетчик; минимальное значение - 0; максимальное значение - 1000; значение по умолчанию - 0
  count: 100
# Статус; допустимые значения: active, inactive, pending
  status: "active"
# Запрещенные значения; запрещенные значения: 0.0, -1.0, подменное значение: -10.0
  forbidden: 5.5
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateStructYAML(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_GenerateStructYAML_ErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
	}{
		{
			name:        "Not a struct",
			input:       "not a struct",
			expectError: true,
		},
		{
			name:        "Nil input",
			input:       nil,
			expectError: true,
		},
		{
			name:        "Valid struct",
			input:       TestStruct{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateStructYAML(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
			}
		})
	}
}
