package main

import (
	"fmt"
	"log"
	"time"

	"github.com/envx/envx"
)

type CustomType struct {
	Value string
}

func (c *CustomType) Decode(value string) error {
	c.Value = "decoded:" + value
	return nil
}

type AdvancedConfig struct {
	StringField   string        `envconfig:"STRING_FIELD"`
	IntField      int           `envconfig:"INT_FIELD"`
	BoolField     bool          `envconfig:"BOOL_FIELD"`
	FloatField    float64       `envconfig:"FLOAT_FIELD"`
	DurationField time.Duration `envconfig:"DURATION_FIELD"`

	SliceOfStrings []string          `envconfig:"SLICE_STRINGS"`
	SliceOfInts    []int             `envconfig:"SLICE_INTS"`
	MapField       map[string]string `envconfig:"MAP_FIELD"`

	CustomField CustomType `envconfig:"CUSTOM_FIELD"`
	PtrField    *string    `envconfig:"PTR_FIELD"`

	DefaultField string `envconfig:"DEFAULT_FIELD" default:"default_value"`
	IgnoredField string `envconfig:"IGNORED_FIELD" ignored:"true"`
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
