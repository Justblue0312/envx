package main

import (
	"fmt"
	"log"
	"time"

	"github.com/justblue0312/envx"
)

type CustomType struct {
	Value string
}

func (c *CustomType) Decode(value string) error {
	c.Value = "decoded:" + value
	return nil
}

type AdvancedConfig struct {
	StringField   string        `envx:"STRING_FIELD"`
	IntField      int           `envx:"INT_FIELD"`
	BoolField     bool          `envx:"BOOL_FIELD"`
	FloatField    float64       `envx:"FLOAT_FIELD"`
	DurationField time.Duration `envx:"DURATION_FIELD"`

	SliceOfStrings []string          `envx:"SLICE_STRINGS"`
	SliceOfInts    []int             `envx:"SLICE_INTS"`
	MapField       map[string]string `envx:"MAP_FIELD"`

	CustomField CustomType `envx:"CUSTOM_FIELD"`
	PtrField    *string    `envx:"PTR_FIELD"`

	DefaultField string `envx:"DEFAULT_FIELD" default:"default_value"`
	IgnoredField string `envx:"IGNORED_FIELD" ignored:"true"`
}

func main() {
	config := &AdvancedConfig{}

	if err := envx.Process("", config); err != nil {
		log.Fatalf("Failed to process config: %v", err)
	}

	fmt.Printf("Advanced Config:\n")
	fmt.Printf("  StringField: %s\n", config.StringField)
	fmt.Printf("  IntField: %d\n", config.IntField)
	fmt.Printf("  BoolField: %v\n", config.BoolField)
	fmt.Printf("  FloatField: %f\n", config.FloatField)
	fmt.Printf("  DurationField: %v\n", config.DurationField)
	fmt.Printf("  SliceOfStrings: %v\n", config.SliceOfStrings)
	fmt.Printf("  SliceOfInts: %v\n", config.SliceOfInts)
	fmt.Printf("  MapField: %v\n", config.MapField)
	fmt.Printf("  CustomField: %+v\n", config.CustomField)
	if config.PtrField != nil {
		fmt.Printf("  PtrField: %s\n", *config.PtrField)
	}
	fmt.Printf("  DefaultField: %s\n", config.DefaultField)
	fmt.Printf("  IgnoredField: %s\n", config.IgnoredField)
}
