package main

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/justblue0312/envx"
)

// CustomType implements Decoder interface
type CustomType struct {
	Value string
}

func (c *CustomType) Decode(value string) error {
	c.Value = "decoded:" + value
	return nil
}

// CustomSetter implements Setter interface
type CustomSetter struct {
	Data string
}

func (c *CustomSetter) Set(value string) error {
	c.Data = "set:" + value
	return nil
}

type ComprehensiveConfig struct {
	// Basic Types
	StringField  string  `envx:"STRING_FIELD"`
	IntField     int     `envx:"INT_FIELD"`
	Int8Field    int8    `envx:"INT8_FIELD"`
	Int16Field   int16   `envx:"INT16_FIELD"`
	Int32Field   int32   `envx:"INT32_FIELD"`
	Int64Field   int64   `envx:"INT64_FIELD"`
	UintField    uint    `envx:"UINT_FIELD"`
	Uint8Field   uint8   `envx:"UINT8_FIELD"`
	Uint16Field  uint16  `envx:"UINT16_FIELD"`
	Uint32Field  uint32  `envx:"UINT32_FIELD"`
	Uint64Field  uint64  `envx:"UINT64_FIELD"`
	BoolField    bool    `envx:"BOOL_FIELD"`
	Float32Field float32 `envx:"FLOAT32_FIELD"`
	Float64Field float64 `envx:"FLOAT64_FIELD"`

	// Time Types
	DurationField time.Duration  `envx:"DURATION_FIELD"`
	LocationField *time.Location `envx:"LOCATION_FIELD"`

	// URL Type
	URLField *url.URL `envx:"URL_FIELD"`

	// Slices
	StringSlice []string `envx:"STRING_SLICE"`
	IntSlice    []int    `envx:"INT_SLICE"`
	ByteSlice   []byte   `envx:"BYTE_SLICE"`

	// Slices of Pointers
	StrPtrSlice []*string `envx:"STR_PTR_SLICE"`
	IntPtrSlice []*int    `envx:"INT_PTR_SLICE"`

	// Maps
	StringMap  map[string]string `envx:"STRING_MAP"`
	IntMap     map[string]int    `envx:"INT_MAP"`
	ComplexMap map[int]string    `envx:"COMPLEX_MAP"`

	// Pointers
	StringPtr *string  `envx:"STRING_PTR"`
	IntPtr    *int     `envx:"INT_PTR"`
	BoolPtr   *bool    `envx:"BOOL_PTR"`
	FloatPtr  *float64 `envx:"FLOAT_PTR"`

	// Custom Types
	DecoderField CustomType   `envx:"DECODER_FIELD"`
	SetterField  CustomSetter `envx:"SETTER_FIELD"`

	// Nested struct
	NestedConfig struct {
		NestedString string `envx:"NESTED_STRING"`
		NestedInt    int    `envx:"NESTED_INT"`
	} `nested:"true"`
}

func main() {
	config := &ComprehensiveConfig{}

	if err := envx.Process("", config); err != nil {
		log.Fatalf("Failed to process config: %v", err)
	}

	// Print all results for verification
	fmt.Printf("=== Basic Types ===\n")
	fmt.Printf("StringField: %s\n", config.StringField)
	fmt.Printf("IntField: %d\n", config.IntField)
	fmt.Printf("Int8Field: %d\n", config.Int8Field)
	fmt.Printf("Int16Field: %d\n", config.Int16Field)
	fmt.Printf("Int32Field: %d\n", config.Int32Field)
	fmt.Printf("Int64Field: %d\n", config.Int64Field)
	fmt.Printf("UintField: %d\n", config.UintField)
	fmt.Printf("Uint8Field: %d\n", config.Uint8Field)
	fmt.Printf("Uint16Field: %d\n", config.Uint16Field)
	fmt.Printf("Uint32Field: %d\n", config.Uint32Field)
	fmt.Printf("Uint64Field: %d\n", config.Uint64Field)
	fmt.Printf("BoolField: %t\n", config.BoolField)
	fmt.Printf("Float32Field: %f\n", config.Float32Field)
	fmt.Printf("Float64Field: %f\n", config.Float64Field)

	fmt.Printf("\n=== Time Types ===\n")
	fmt.Printf("DurationField: %v\n", config.DurationField)
	if config.LocationField != nil {
		fmt.Printf("LocationField: %v\n", config.LocationField)
	}

	fmt.Printf("\n=== URL Type ===\n")
	if config.URLField != nil {
		fmt.Printf("URLField: %s\n", config.URLField.String())
	}

	fmt.Printf("\n=== Slices ===\n")
	fmt.Printf("StringSlice: %v\n", config.StringSlice)
	fmt.Printf("IntSlice: %v\n", config.IntSlice)
	fmt.Printf("ByteSlice: %v\n", config.ByteSlice)
	fmt.Printf("StrPtrSlice: %v\n", formatStrPtrSlice(config.StrPtrSlice))
	fmt.Printf("IntPtrSlice: %v\n", formatIntPtrSlice(config.IntPtrSlice))

	fmt.Printf("\n=== Maps ===\n")
	fmt.Printf("StringMap: %v\n", config.StringMap)
	fmt.Printf("IntMap: %v\n", config.IntMap)
	fmt.Printf("ComplexMap: %v\n", config.ComplexMap)

	fmt.Printf("\n=== Pointers ===\n")
	if config.StringPtr != nil {
		fmt.Printf("StringPtr: %s\n", *config.StringPtr)
	}
	if config.IntPtr != nil {
		fmt.Printf("IntPtr: %d\n", *config.IntPtr)
	}
	if config.BoolPtr != nil {
		fmt.Printf("BoolPtr: %t\n", *config.BoolPtr)
	}
	if config.FloatPtr != nil {
		fmt.Printf("FloatPtr: %f\n", *config.FloatPtr)
	}

	fmt.Printf("\n=== Custom Types ===\n")
	fmt.Printf("DecoderField: %+v\n", config.DecoderField)
	fmt.Printf("SetterField: %+v\n", config.SetterField)

	fmt.Printf("\n=== Nested Struct ===\n")
	fmt.Printf("NestedString: %s\n", config.NestedConfig.NestedString)
	fmt.Printf("NestedInt: %d\n", config.NestedConfig.NestedInt)

	fmt.Printf("\n=== All tests completed successfully! ===\n")
}

func formatStrPtrSlice(slice []*string) string {
	if slice == nil {
		return "nil"
	}
	result := "["
	for i, s := range slice {
		if i > 0 {
			result += ", "
		}
		if s != nil {
			result += fmt.Sprintf("%s", *s)
		} else {
			result += "nil"
		}
	}
	result += "]"
	return result
}

func formatIntPtrSlice(slice []*int) string {
	if slice == nil {
		return "nil"
	}
	result := "["
	for i, n := range slice {
		if i > 0 {
			result += ", "
		}
		if n != nil {
			result += fmt.Sprintf("%d", *n)
		} else {
			result += "nil"
		}
	}
	result += "]"
	return result
}
