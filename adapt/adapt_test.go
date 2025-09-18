package adapter

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type customType float32

type IncorrectTags struct {
	PtrUint   *uint `json:"pointerUint-for-min" rst-min:"5" rst-max:"500"`
	IntForMax int   `json:"int-for-max" rst-min:"fail" rst-max:"fail"`
}

type IncorrectField struct {
	IntForMax string `json:"int-for-max" rst-min:"5"`
}

type CorrectNested1 struct {
	StringForRegex string `json:"string-for-regex" rst-usage:"test string" rst-regex:"[^a-zA-Z0-9 ]+"`
	nonExportable  bool
}

type CorrectNested2 struct {
	IntForMin     int     `json:"int-for-min" rst-default:"100" rst-min:"5" rst-max:"500"`
	Float32ForMax float32 `json:"int-for-max" rst-max:"50.54"`
}

type CorrectNested3 struct {
	FloatForMax   customType `json:"float32-for-max" rst-default:"100" rst-min:"5" rst-max:"500"`
	nonExportable bool
}

type CorrectRoot struct {
	IfaceWithValue   any
	IfaceWithPointer any
	NestedStrValue   CorrectNested1
	NestedStrPointer *CorrectNested3 `json:"NestedStrValue" rst-info:"change name"`
	UintForMin       *uint16         `json:"uint16-for-min" rst-default:"100" rst-min:"5" rst-max:"500"`
	nonExportable    float32
	//New
	StringForDefault  string  `rst-default:"default"`
	IntForChoice      int     `rst-choice:"2000||2001||2003"`
	FloatForForbidden float32 `rst-forbidden:"23.34||30.56**40.65"`
	//Slice&Map
	SliceForMax    []int             `rst-max:"5"`
	Slice2ForMax   [][]int           `rst-max:"5"`
	MapForMin      map[string]int    `rst-min:"5"`
	SliceStruct    []CorrectNested1  `rst-max:"5"`
	SlicePtrStruct []*CorrectNested1 `rst-max:"5"`
	SliceAny       []any             `rst-max:"5"`
}

// Новые структуры для расширенного тестирования
type ExtendedTestStruct struct {
	// Choice тесты
	IntChoice    int     `rst-choice:"10||20||30"`
	UintChoice   uint    `rst-choice:"100||200||300"`
	FloatChoice  float64 `rst-choice:"1.5||2.5||3.5"`
	StringChoice string  `rst-choice:"apple||banana||orange"`

	// Default тесты
	IntDefault    int     `rst-default:"42"`
	UintDefault   uint    `rst-default:"123"`
	FloatDefault  float64 `rst-default:"3.14"`
	StringDefault string  `rst-default:"hello"`

	// Forbidden тесты
	IntForbidden    int     `rst-forbidden:"1||2||3**10"`
	UintForbidden   uint    `rst-forbidden:"10||20||30**100"`
	FloatForbidden  float64 `rst-forbidden:"1.1||2.2||3.3**10.0"`
	StringForbidden string  `rst-forbidden:"bad||evil||wrong**good"`

	// Min/Max тесты
	IntMin   int     `rst-min:"10"`
	IntMax   int     `rst-max:"100"`
	UintMin  uint    `rst-min:"5"`
	UintMax  uint    `rst-max:"50"`
	FloatMin float64 `rst-min:"1.5"`
	FloatMax float64 `rst-max:"10.5"`

	// Regex тесты
	StringRegex string `rst-regex:"[^a-zA-Z]+"`

	// Комбинированные тесты
	CombinedField int `rst-min:"5" rst-max:"100" rst-default:"50"`
}

type SliceTestStruct struct {
	IntSlice    []int            `rst-min:"5"`
	FloatSlice  []float64        `rst-max:"10.0"`
	StringSlice []string         `rst-choice:"a||b||c"`
	StructSlice []CorrectNested1 `rst-max:"3"`
}

type MapTestStruct struct {
	IntMap    map[string]int     `rst-min:"10"`
	FloatMap  map[string]float64 `rst-max:"20.0"`
	StringMap map[string]string  `rst-regex:"[^0-9]+"`
}

type PointerTestStruct struct {
	IntPtr    *int     `rst-min:"5" rst-max:"100"`
	FloatPtr  *float64 `rst-default:"3.14"`
	StringPtr *string  `rst-choice:"yes||no"`
}

type ArrayTestStruct struct {
	IntArray    [3]int     `rst-min:"1"`
	FloatArray  [2]float64 `rst-max:"5.0"`
	StringArray [2]string  `rst-regex:"[^x]+"`
}

var a = adapter{logger: log.New(os.Stdout, "adapter ", log.LstdFlags)}

func Test_AdaptIncorrectInput(t *testing.T) {
	uintValue := uint(2)

	incorrectTags := IncorrectTags{
		PtrUint:   &uintValue,
		IntForMax: 501,
	}

	incorrectField := IncorrectField{
		IntForMax: "fail",
	}

	t.Run("Incorrect input type", func(t *testing.T) {
		strct, err := a.AdaptStruct(1)
		assert.ErrorIs(t, err, ErrNotStruct)
		assert.Nil(t, strct)
	})

	t.Run("Error adapt struct", func(t *testing.T) {
		strct, err := a.AdaptStruct(incorrectField)
		assert.ErrorIs(t, err, ErrInvalidTags)
		assert.Nil(t, strct)

		strct, err = a.AdaptStruct(incorrectTags)
		assert.Error(t, err)
		assert.Nil(t, strct)
		assert.Equal(t, uint(2), *incorrectTags.PtrUint)
	})
}

func Test_AdaptCorrectInput(t *testing.T) {
	uintValue := uint16(2)
	ptrValue := int(6)
	expectedUintValue := uint16(5)
	expectedPtrValue := int(5)

	correctNested1 := CorrectNested1{
		StringForRegex: "-678))++2",
		nonExportable:  false,
	}

	expectedNested1 := CorrectNested1{
		StringForRegex: "6782",
		nonExportable:  false,
	}

	correctPtrNested1 := CorrectNested1{
		StringForRegex: "-678))++2",
		nonExportable:  false,
	}

	expectedPtrNested1 := CorrectNested1{
		StringForRegex: "6782",
		nonExportable:  false,
	}

	correctNested2 := CorrectNested2{
		IntForMin:     2,
		Float32ForMax: 50.6,
	}

	expectedNested2 := CorrectNested2{
		IntForMin:     5,
		Float32ForMax: 50.54,
	}

	correctNested3 := CorrectNested3{
		FloatForMax:   506,
		nonExportable: false,
	}

	expectedNested3 := CorrectNested3{
		FloatForMax:   500,
		nonExportable: false,
	}

	correctRoot := CorrectRoot{
		IfaceWithValue:    correctNested1,
		IfaceWithPointer:  &correctNested2,
		NestedStrValue:    correctNested1,
		NestedStrPointer:  &correctNested3,
		UintForMin:        &uintValue,
		nonExportable:     10,
		StringForDefault:  "",
		IntForChoice:      1999,
		FloatForForbidden: 30.56,
		SliceForMax:       []int{6, 4, 7},
		MapForMin:         map[string]int{"1": 1, "2": 3},
		SliceStruct:       []CorrectNested1{correctNested1},
		SlicePtrStruct:    []*CorrectNested1{&correctPtrNested1},
		SliceAny:          []any{&ptrValue, correctNested1, 7},
	}

	sl := make([]int, 1)
	sl[0] = 6

	correctRoot.Slice2ForMax = make([][]int, 1)
	correctRoot.Slice2ForMax[0] = sl

	expectedRoot := CorrectRoot{
		IfaceWithValue:    expectedNested1,
		IfaceWithPointer:  &expectedNested2,
		NestedStrValue:    expectedNested1,
		NestedStrPointer:  &expectedNested3,
		UintForMin:        &expectedUintValue,
		nonExportable:     10,
		StringForDefault:  "default",
		IntForChoice:      2000,
		FloatForForbidden: 40.65,
		SliceForMax:       []int{5, 4, 5},
		Slice2ForMax:      [][]int{},
		MapForMin:         map[string]int{"1": 5, "2": 5},
		SliceStruct:       []CorrectNested1{expectedNested1},
		SlicePtrStruct:    []*CorrectNested1{&expectedPtrNested1},
		SliceAny:          []any{&expectedPtrValue, expectedNested1, 5},
	}

	exsl := make([]int, 1)
	exsl[0] = 5

	expectedRoot.Slice2ForMax = make([][]int, 1)
	expectedRoot.Slice2ForMax[0] = exsl

	t.Run("Adapt Struct", func(t *testing.T) {
		actualRoot, err := a.AdaptStruct(correctRoot)
		assert.NoError(t, err)
		assert.Equal(t, expectedRoot, actualRoot)

		uintValue = 3
		correctNested2.IntForMin = 1
		correctNested3.FloatForMax = 521

		actualRoot, err = a.AdaptStruct(&correctRoot)
		assert.Equal(t, expectedRoot, actualRoot)
		assert.Equal(t, expectedRoot, correctRoot)
		assert.NoError(t, err)
	})
}

func Test_ExtendedChoice(t *testing.T) {
	t.Run("Int Choice", func(t *testing.T) {
		test := ExtendedTestStruct{IntChoice: 15}
		expected := ExtendedTestStruct{IntChoice: 10}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected.IntChoice, result.(ExtendedTestStruct).IntChoice)

		// Тест с валидным значением
		test.IntChoice = 20
		expected.IntChoice = 20
		result, err = a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected.IntChoice, result.(ExtendedTestStruct).IntChoice)
	})

	t.Run("Uint Choice", func(t *testing.T) {
		test := ExtendedTestStruct{UintChoice: 150}
		expected := ExtendedTestStruct{UintChoice: 100}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected.UintChoice, result.(ExtendedTestStruct).UintChoice)
	})

	t.Run("Float Choice", func(t *testing.T) {
		test := ExtendedTestStruct{FloatChoice: 2.0}
		expected := ExtendedTestStruct{FloatChoice: 1.5}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected.FloatChoice, result.(ExtendedTestStruct).FloatChoice)
	})

	t.Run("String Choice", func(t *testing.T) {
		test := ExtendedTestStruct{StringChoice: "grape"}
		expected := ExtendedTestStruct{StringChoice: "apple"}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected.StringChoice, result.(ExtendedTestStruct).StringChoice)
	})

	t.Run("Zero Value Choice", func(t *testing.T) {
		test := ExtendedTestStruct{} // все поля zero value
		expected := ExtendedTestStruct{}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected.IntChoice, result.(ExtendedTestStruct).IntChoice)
		assert.Equal(t, expected.UintChoice, result.(ExtendedTestStruct).UintChoice)
		assert.Equal(t, expected.FloatChoice, result.(ExtendedTestStruct).FloatChoice)
		assert.Equal(t, expected.StringChoice, result.(ExtendedTestStruct).StringChoice)
	})
}

func Test_ExtendedDefault(t *testing.T) {
	t.Run("Int Default", func(t *testing.T) {
		test := ExtendedTestStruct{IntDefault: 0} // zero value
		expected := ExtendedTestStruct{IntDefault: 42}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected.IntDefault, result.(ExtendedTestStruct).IntDefault)

		// Тест с не-zero значением
		test.IntDefault = 100
		expected.IntDefault = 100
		result, err = a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected.IntDefault, result.(ExtendedTestStruct).IntDefault)
	})

	t.Run("Uint Default", func(t *testing.T) {
		test := ExtendedTestStruct{UintDefault: 0}
		expected := ExtendedTestStruct{UintDefault: 123}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected.UintDefault, result.(ExtendedTestStruct).UintDefault)
	})

	t.Run("Float Default", func(t *testing.T) {
		test := ExtendedTestStruct{FloatDefault: 0.0}
		expected := ExtendedTestStruct{FloatDefault: 3.14}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected.FloatDefault, result.(ExtendedTestStruct).FloatDefault)
	})

	t.Run("String Default", func(t *testing.T) {
		test := ExtendedTestStruct{StringDefault: ""}
		expected := ExtendedTestStruct{StringDefault: "hello"}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected.StringDefault, result.(ExtendedTestStruct).StringDefault)
	})
}

func Test_ExtendedForbidden(t *testing.T) {
	t.Run("Int Forbidden", func(t *testing.T) {
		type TestStruct struct {
			IntForbidden int `rst-forbidden:"1||2||3**10"`
		}
		test := TestStruct{IntForbidden: 2}
		expected := TestStruct{IntForbidden: 10}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)

		// Тест с не-запрещенным значением
		test.IntForbidden = 5
		expected.IntForbidden = 5
		result, err = a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Uint Forbidden", func(t *testing.T) {
		type TestStruct struct {
			UintForbidden uint `rst-forbidden:"10||20||30**100"`
		}
		test := TestStruct{UintForbidden: 20}
		expected := TestStruct{UintForbidden: 100}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Float Forbidden", func(t *testing.T) {
		type TestStruct struct {
			FloatForbidden float64 `rst-forbidden:"1.1||2.2||3.3**10.0"`
		}
		test := TestStruct{FloatForbidden: 2.2}
		expected := TestStruct{FloatForbidden: 10.0}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("String Forbidden", func(t *testing.T) {
		type TestStruct struct {
			StringForbidden string `rst-forbidden:"bad||evil||wrong**good"`
		}
		test := TestStruct{StringForbidden: "evil"}
		expected := TestStruct{StringForbidden: "good"}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func Test_ExtendedMinMax(t *testing.T) {
	t.Run("Int Min", func(t *testing.T) {
		type TestStruct struct {
			IntMin int `rst-min:"10"`
		}
		test := TestStruct{IntMin: 5}
		expected := TestStruct{IntMin: 10}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Int Max", func(t *testing.T) {
		type TestStruct struct {
			IntMax int `rst-max:"100"`
		}
		test := TestStruct{IntMax: 150}
		expected := TestStruct{IntMax: 100}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Uint Min", func(t *testing.T) {
		type TestStruct struct {
			UintMin uint `rst-min:"5"`
		}
		test := TestStruct{UintMin: 3}
		expected := TestStruct{UintMin: 5}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Uint Max", func(t *testing.T) {
		type TestStruct struct {
			UintMax uint `rst-max:"50"`
		}
		test := TestStruct{UintMax: 60}
		expected := TestStruct{UintMax: 50}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Float Min", func(t *testing.T) {
		type TestStruct struct {
			FloatMin float64 `rst-min:"1.5"`
		}
		test := TestStruct{FloatMin: 1.0}
		expected := TestStruct{FloatMin: 1.5}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Float Max", func(t *testing.T) {
		type TestStruct struct {
			FloatMax float64 `rst-max:"10.5"`
		}
		test := TestStruct{FloatMax: 15.0}
		expected := TestStruct{FloatMax: 10.5}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func Test_ExtendedRegex(t *testing.T) {
	t.Run("String Regex", func(t *testing.T) {
		type TestStruct struct {
			StringRegex string `rst-regex:"[^a-zA-Z]+"`
		}
		test := TestStruct{StringRegex: "hello123world"}
		expected := TestStruct{StringRegex: "helloworld"}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("String Regex with special chars", func(t *testing.T) {
		type TestStruct struct {
			StringRegex string `rst-regex:"[^a-zA-Z]+"`
		}
		test := TestStruct{StringRegex: "test@#$%^&*()test"}
		expected := TestStruct{StringRegex: "testtest"}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func Test_ExtendedCombined(t *testing.T) {
	t.Run("Combined Field", func(t *testing.T) {
		type TestStruct struct {
			CombinedField int `rst-default:"50" rst-min:"5" rst-max:"100"`
		}
		test := TestStruct{CombinedField: 0}      // zero value
		expected := TestStruct{CombinedField: 50} // default

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)

		// Тест с значением меньше минимума
		test.CombinedField = 3
		expected.CombinedField = 5 // min
		result, err = a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)

		// Тест с значением больше максимума
		test.CombinedField = 150
		expected.CombinedField = 100 // max
		result, err = a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func Test_SliceAdaptation(t *testing.T) {
	t.Run("Int Slice Min", func(t *testing.T) {
		type TestStruct struct {
			IntSlice []int `rst-min:"5"`
		}
		test := TestStruct{
			IntSlice: []int{3, 7, 2, 8},
		}
		expected := TestStruct{
			IntSlice: []int{5, 7, 5, 8},
		}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Float Slice Max", func(t *testing.T) {
		type TestStruct struct {
			FloatSlice []float64 `rst-max:"10.0"`
		}
		test := TestStruct{
			FloatSlice: []float64{15.0, 8.0, 12.0},
		}
		expected := TestStruct{
			FloatSlice: []float64{10.0, 8.0, 10.0},
		}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("String Slice Choice", func(t *testing.T) {
		type TestStruct struct {
			StringSlice []string `rst-choice:"a||b||c"`
		}
		test := TestStruct{
			StringSlice: []string{"d", "a", "e"},
		}
		expected := TestStruct{
			StringSlice: []string{"a", "a", "a"},
		}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func Test_MapAdaptation(t *testing.T) {
	t.Run("Int Map Min", func(t *testing.T) {
		type TestStruct struct {
			IntMap map[string]int `rst-min:"10"`
		}
		test := TestStruct{
			IntMap: map[string]int{"a": 5, "b": 15, "c": 8},
		}
		expected := TestStruct{
			IntMap: map[string]int{"a": 10, "b": 15, "c": 10},
		}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Float Map Max", func(t *testing.T) {
		type TestStruct struct {
			FloatMap map[string]float64 `rst-max:"20.0"`
		}
		test := TestStruct{
			FloatMap: map[string]float64{"a": 25.0, "b": 15.0, "c": 30.0},
		}
		expected := TestStruct{
			FloatMap: map[string]float64{"a": 20.0, "b": 15.0, "c": 20.0},
		}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("String Map Regex", func(t *testing.T) {
		type TestStruct struct {
			StringMap map[string]string `rst-regex:"[0-9]+"`
		}
		test := TestStruct{
			StringMap: map[string]string{"a": "hello123", "b": "world456", "c": "test"},
		}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		resultStruct := result.(TestStruct)
		assert.Equal(t, "hello", resultStruct.StringMap["a"])
		assert.Equal(t, "world", resultStruct.StringMap["b"])
		assert.Equal(t, "test", resultStruct.StringMap["c"])
	})
}

func Test_PointerAdaptation(t *testing.T) {
	t.Run("Int Pointer", func(t *testing.T) {
		type TestStruct struct {
			IntPtr *int `rst-min:"5" rst-max:"100"`
		}
		val := 3
		test := TestStruct{IntPtr: &val}
		expectedVal := 5
		expected := TestStruct{IntPtr: &expectedVal}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, *expected.IntPtr, *result.(TestStruct).IntPtr)
	})

	t.Run("Float Pointer Default", func(t *testing.T) {
		type TestStruct struct {
			FloatPtr *float64 `rst-default:"3.14"`
		}
		test := TestStruct{FloatPtr: nil}
		expectedVal := 3.14
		expected := TestStruct{FloatPtr: &expectedVal}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		resultStruct := result.(TestStruct)
		assert.NotNil(t, resultStruct.FloatPtr)
		assert.Equal(t, *expected.FloatPtr, *resultStruct.FloatPtr)
	})

	t.Run("String Pointer Choice", func(t *testing.T) {
		type TestStruct struct {
			StringPtr *string `rst-choice:"yes||no"`
		}
		val := "maybe"
		test := TestStruct{StringPtr: &val}
		expectedVal := "yes"
		expected := TestStruct{StringPtr: &expectedVal}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, *expected.StringPtr, *result.(TestStruct).StringPtr)
	})
}

func Test_ArrayAdaptation(t *testing.T) {
	t.Run("Int Array Min", func(t *testing.T) {
		type TestStruct struct {
			IntArray [3]int `rst-min:"1"`
		}
		test := TestStruct{
			IntArray: [3]int{0, 2, 5},
		}
		expected := TestStruct{
			IntArray: [3]int{1, 2, 5},
		}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Float Array Max", func(t *testing.T) {
		type TestStruct struct {
			FloatArray [2]float64 `rst-max:"5.0"`
		}
		test := TestStruct{
			FloatArray: [2]float64{6.0, 3.0},
		}
		expected := TestStruct{
			FloatArray: [2]float64{5.0, 3.0},
		}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("String Array Regex", func(t *testing.T) {
		type TestStruct struct {
			StringArray [2]string `rst-regex:"[0-9]+"`
		}
		test := TestStruct{
			StringArray: [2]string{"hello123", "world456"},
		}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		resultStruct := result.(TestStruct)
		assert.Equal(t, "hello", resultStruct.StringArray[0])
		assert.Equal(t, "world", resultStruct.StringArray[1])
	})
}

func Test_ErrorCases(t *testing.T) {
	t.Run("Invalid Forbidden Format", func(t *testing.T) {
		type InvalidForbidden struct {
			Field int `rst-forbidden:"invalid"`
		}

		test := InvalidForbidden{Field: 1}
		_, err := a.AdaptStruct(test)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidTags)
	})

	t.Run("Invalid Choice Format", func(t *testing.T) {
		type InvalidChoice struct {
			Field int `rst-choice:"invalid||format"`
		}

		test := InvalidChoice{Field: 1}
		_, err := a.AdaptStruct(test)
		assert.Error(t, err)
	})

	t.Run("Invalid Min/Max Values", func(t *testing.T) {
		type InvalidMinMax struct {
			Field int `rst-min:"invalid" rst-max:"also_invalid"`
		}

		test := InvalidMinMax{Field: 1}
		_, err := a.AdaptStruct(test)
		assert.Error(t, err)
	})

	t.Run("Invalid Default Value", func(t *testing.T) {
		type InvalidDefault struct {
			Field int `rst-default:"invalid"`
		}

		test := InvalidDefault{Field: 0}
		_, err := a.AdaptStruct(test)
		assert.Error(t, err)
	})
}

func Test_EdgeCases(t *testing.T) {
	t.Run("Empty String Regex", func(t *testing.T) {
		type TestStruct struct {
			StringRegex string `rst-regex:"[^a-zA-Z]+"`
		}
		test := TestStruct{StringRegex: ""}
		expected := TestStruct{StringRegex: ""}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Zero Values", func(t *testing.T) {
		type TestStruct struct {
			IntDefault    int     `rst-default:"42"`
			UintDefault   uint    `rst-default:"123"`
			FloatDefault  float64 `rst-default:"3.14"`
			StringDefault string  `rst-default:"hello"`
		}
		test := TestStruct{}
		expected := TestStruct{
			IntDefault:    42,
			UintDefault:   123,
			FloatDefault:  3.14,
			StringDefault: "hello",
		}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Nil Pointer", func(t *testing.T) {
		type TestStruct struct {
			FloatPtr *float64 `rst-default:"3.14"`
		}
		test := TestStruct{}
		expectedVal := 3.14
		expected := TestStruct{FloatPtr: &expectedVal}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, *expected.FloatPtr, *result.(TestStruct).FloatPtr)
	})

	t.Run("Empty Slice", func(t *testing.T) {
		type TestStruct struct {
			IntSlice []int `rst-min:"5"`
		}
		test := TestStruct{IntSlice: []int{}}
		expected := TestStruct{IntSlice: []int{}}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Empty Map", func(t *testing.T) {
		type TestStruct struct {
			IntMap map[string]int `rst-min:"10"`
		}
		test := TestStruct{IntMap: map[string]int{}}
		expected := TestStruct{IntMap: map[string]int{}}

		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func Test_ParseStructTag(t *testing.T) {
	t.Run("Valid Tags", func(t *testing.T) {
		type TestStruct struct {
			Field string `rst-default:"hello"`
		}

		test := TestStruct{Field: ""}
		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, "hello", result.(TestStruct).Field)
	})

	t.Run("Invalid Tag Format", func(t *testing.T) {
		type TestStruct struct {
			Field string `invalid:"format"`
		}

		test := TestStruct{Field: ""}
		result, err := a.AdaptStruct(test)
		assert.NoError(t, err) // неизвестные теги игнорируются безопасным парсером
		assert.Equal(t, test, result)
	})

	t.Run("Empty Tag", func(t *testing.T) {
		type TestStruct struct {
			Field string ``
		}

		test := TestStruct{Field: ""}
		result, err := a.AdaptStruct(test)
		assert.NoError(t, err)
		assert.Equal(t, test, result)
	})
}
