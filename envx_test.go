package envx

import (
	"errors"
	"os"
	"testing"
	"time"
)

type TestConfig struct {
	StringField   string `envconfig:"TEST_STRING"`
	IntField      int    `envconfig:"TEST_INT"`
	BoolField     bool   `envconfig:"TEST_BOOL"`
	RequiredField string `envconfig:"TEST_REQUIRED"`
	DefaultField  string `envconfig:"TEST_DEFAULT" default:"default_value"`
	IgnoredField  string `envconfig:"TEST_IGNORED" ignored:"true"`
}

type NestedConfig struct {
	Database struct {
		Host     string `envconfig:"DB_HOST" nested:"true"`
		Port     int    `envconfig:"DB_PORT" nested:"true"`
		Password string `envconfig:"DB_PASSWORD" nested:"true"`
	} `nested:"true"`

	Server struct {
		Host string `envconfig:"SERVER_HOST" nested:"true"`
		Port int    `envconfig:"SERVER_PORT" nested:"true"`
	} `nested:"true"`
}

type CustomType struct {
	Value string
}

func (c *CustomType) Decode(value string) error {
	c.Value = "decoded:" + value
	return nil
}

func (c *CustomType) Set(value string) error {
	c.Value = "set:" + value
	return nil
}

type CustomUnmarshaler struct {
	Value string
}

func (c *CustomUnmarshaler) UnmarshalText(text []byte) error {
	c.Value = "unmarshaled:" + string(text)
	return nil
}

type CustomBinaryUnmarshaler struct {
	Value string
}

func (c *CustomBinaryUnmarshaler) UnmarshalBinary(data []byte) error {
	c.Value = "binary:" + string(data)
	return nil
}

type ComplexConfig struct {
	CustomDecoder   CustomType
	CustomSetter    CustomType
	CustomUnmarshal CustomUnmarshaler
	CustomBinary    CustomBinaryUnmarshaler

	SliceField []string
	IntSlice   []int
	MapField   map[string]string
	Duration   time.Duration
	PtrField   *string
	NestedPtr  *TestConfig
}

func TestProcess(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		spec     interface{}
		prefix   string
		expected interface{}
		wantErr  bool
		errType  error
	}{
		{
			name: "basic string field",
			env:  map[string]string{"TEST_STRING": "hello"},
			spec: &TestConfig{},
			expected: &TestConfig{
				StringField:  "hello",
				DefaultField: "default_value",
			},
		},
		{
			name: "int field",
			env:  map[string]string{"TEST_INT": "42"},
			spec: &TestConfig{},
			expected: &TestConfig{
				IntField:     42,
				DefaultField: "default_value",
			},
		},
		{
			name: "bool field",
			env:  map[string]string{"TEST_BOOL": "true"},
			spec: &TestConfig{},
			expected: &TestConfig{
				BoolField:    true,
				DefaultField: "default_value",
			},
		},
		{
			name: "default value",
			env:  map[string]string{},
			spec: &TestConfig{},
			expected: &TestConfig{
				DefaultField: "default_value",
			},
		},

		{
			name:   "with prefix",
			env:    map[string]string{"APP_TEST_STRING": "prefixed"},
			spec:   &TestConfig{},
			prefix: "APP",
			expected: &TestConfig{
				StringField:  "prefixed",
				DefaultField: "default_value",
			},
		},
		{
			name: "custom decoder",
			env:  map[string]string{"CUSTOMDECODER": "test"},
			spec: &ComplexConfig{},
			expected: &ComplexConfig{
				CustomDecoder: CustomType{Value: "decoded:test"},
			},
		},
		{
			name: "custom setter",
			env:  map[string]string{"CUSTOMSETTER": "test"},
			spec: &ComplexConfig{},
			expected: &ComplexConfig{
				CustomSetter: CustomType{Value: "set:test"},
			},
		},
		{
			name: "slice field",
			env:  map[string]string{"SLICEFIELD": "a,b,c"},
			spec: &ComplexConfig{},
			expected: &ComplexConfig{
				SliceField: []string{"a", "b", "c"},
			},
		},
		{
			name: "int slice",
			env:  map[string]string{"INTSLICE": "1,2,3"},
			spec: &ComplexConfig{},
			expected: &ComplexConfig{
				IntSlice: []int{1, 2, 3},
			},
		},
		{
			name: "map field",
			env:  map[string]string{"MAPFIELD": "key1:val1,key2:val2"},
			spec: &ComplexConfig{},
			expected: &ComplexConfig{
				MapField: map[string]string{"key1": "val1", "key2": "val2"},
			},
		},
		{
			name: "duration field",
			env:  map[string]string{"DURATION": "5s"},
			spec: &ComplexConfig{},
			expected: &ComplexConfig{
				Duration: 5 * time.Second,
			},
		},
		{
			name: "pointer field",
			env:  map[string]string{"PTRFIELD": "pointer"},
			spec: &ComplexConfig{},
			expected: &ComplexConfig{
				PtrField: func() *string { s := "pointer"; return &s }(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				if err := os.Setenv(k, v); err != nil {
					t.Fatalf("Failed to set env var %s: %v", k, err)
				}
				defer func(k string) {
					if err := os.Unsetenv(k); err != nil {
						t.Logf("Failed to unset env var %s: %v", k, err)
					}
				}(k)
			}

			err := Process(tt.prefix, tt.spec)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Process() expected error, got nil")
					return
				}
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("Process() expected error %v, got %v", tt.errType, err)
					return
				}
				return
			}

			if err != nil {
				t.Errorf("Process() unexpected error: %v", err)
				return
			}

			if !compareConfigs(tt.spec, tt.expected) {
				t.Errorf("Process() = %+v, want %+v", tt.spec, tt.expected)
			}
		})
	}
}

func TestNestedStructs(t *testing.T) {
	env := map[string]string{
		"DATABASE__DB_HOST":     "localhost",
		"DATABASE__DB_PORT":     "5432",
		"DATABASE__DB_PASSWORD": "secret",
		"SERVER__SERVER_HOST":   "0.0.0.0",
		"SERVER__SERVER_PORT":   "8080",
	}

	for k, v := range env {
		if err := os.Setenv(k, v); err != nil {
			t.Fatalf("Failed to set env var %s: %v", k, err)
		}
		defer func(k string) {
			if err := os.Unsetenv(k); err != nil {
				t.Logf("Failed to unset env var %s: %v", k, err)
			}
		}(k)
	}

	config := &NestedConfig{}
	err := Process("", config)
	if err != nil {
		t.Fatalf("Process() unexpected error: %v", err)
	}

	expected := &NestedConfig{
		Database: struct {
			Host     string `envconfig:"DB_HOST" nested:"true"`
			Port     int    `envconfig:"DB_PORT" nested:"true"`
			Password string `envconfig:"DB_PASSWORD" nested:"true"`
		}{
			Host:     "localhost",
			Port:     5432,
			Password: "secret",
		},
		Server: struct {
			Host string `envconfig:"SERVER_HOST" nested:"true"`
			Port int    `envconfig:"SERVER_PORT" nested:"true"`
		}{
			Host: "0.0.0.0",
			Port: 8080,
		},
	}

	if !compareNestedConfigs(config, expected) {
		t.Errorf("Process() = %+v, want %+v", config, expected)
	}
}

func TestMustProcess(t *testing.T) {
	if err := os.Setenv("TEST_STRING", "hello"); err != nil {
		t.Fatalf("Failed to set env var: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("TEST_STRING"); err != nil {
			t.Logf("Failed to unset env var: %v", err)
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustProcess() panicked unexpectedly: %v", r)
		}
	}()

	config := &TestConfig{}
	MustProcess("", config)

	if config.StringField != "hello" {
		t.Errorf("MustProcess() = %+v, want StringField=hello", config)
	}
}

func TestMustProcessPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustProcess() expected panic, but didn't panic")
		}
	}()

	config := (*TestConfig)(nil) // This will cause ErrInvalidSpecification
	MustProcess("", config)
}

func TestCheckDisallowed(t *testing.T) {
	if err := os.Setenv("TEST_STRING", "hello"); err != nil {
		t.Fatalf("Failed to set env var: %v", err)
	}
	if err := os.Setenv("UNKNOWN_VAR", "value"); err != nil {
		t.Fatalf("Failed to set env var: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("TEST_STRING"); err != nil {
			t.Logf("Failed to unset env var: %v", err)
		}
	}()
	defer func() {
		if err := os.Unsetenv("UNKNOWN_VAR"); err != nil {
			t.Logf("Failed to unset env var: %v", err)
		}
	}()

	config := &TestConfig{}
	err := CheckDisallowed("TEST", config)
	if err == nil {
		t.Errorf("CheckDisallowed() expected error for unknown var")
	}
}

func TestInvalidSpecification(t *testing.T) {
	tests := []struct {
		name string
		spec interface{}
	}{
		{"not a pointer", TestConfig{}},
		{"nil pointer", (*TestConfig)(nil)},
		{"pointer to non-struct", new(string)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Process("", tt.spec)
			if err == nil {
				t.Errorf("Process() expected error for invalid spec")
			}
			if !errors.Is(err, ErrInvalidSpecification) {
				t.Errorf("Process() expected ErrInvalidSpecification, got %v", err)
			}
		})
	}
}

func TestParseError(t *testing.T) {
	err := &ParseError{
		KeyName:   "TEST_KEY",
		FieldName: "TestField",
		TypeName:  "string",
		Value:     "invalid",
		Err:       errors.New("conversion error"),
	}

	expected := "envx.Process: assigning TEST_KEY to TestField: converting 'invalid' to type string. details: conversion error"
	if err.Error() != expected {
		t.Errorf("ParseError.Error() = %q, want %q", err.Error(), expected)
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"TestField", "Test_Field"},
		{"XMLHttpRequest", "XML_Http_Request"},
		{"Simple", "Simple"},
		{"UserID", "User_ID"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := toSnakeCase(tt.input); got != tt.want {
				t.Errorf("toSnakeCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsTrue(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"True", true},
		{"TRUE", true},
		{"1", true},
		{"t", true},
		{"T", true},
		{"false", false},
		{"False", false},
		{"FALSE", false},
		{"0", false},
		{"f", false},
		{"F", false},
		{"", false},
		{"invalid", false},
		{"yes", false}, // strconv.ParseBool doesn't accept "yes"
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := isTrue(tt.input); got != tt.want {
				t.Errorf("isTrue(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func compareConfigs(a, b interface{}) bool {
	switch x := a.(type) {
	case *TestConfig:
		y, ok := b.(*TestConfig)
		if !ok {
			return false
		}
		return x.StringField == y.StringField &&
			x.IntField == y.IntField &&
			x.BoolField == y.BoolField &&
			x.RequiredField == y.RequiredField &&
			x.DefaultField == y.DefaultField
	case *ComplexConfig:
		y, ok := b.(*ComplexConfig)
		if !ok {
			return false
		}
		if x.SliceField != nil && y.SliceField != nil {
			if len(x.SliceField) != len(y.SliceField) {
				return false
			}
			for i := range x.SliceField {
				if x.SliceField[i] != y.SliceField[i] {
					return false
				}
			}
		}
		if x.IntSlice != nil && y.IntSlice != nil {
			if len(x.IntSlice) != len(y.IntSlice) {
				return false
			}
			for i := range x.IntSlice {
				if x.IntSlice[i] != y.IntSlice[i] {
					return false
				}
			}
		}
		if x.MapField != nil && y.MapField != nil {
			if len(x.MapField) != len(y.MapField) {
				return false
			}
			for k, v := range x.MapField {
				if y.MapField[k] != v {
					return false
				}
			}
		}
		return x.Duration == y.Duration
	}
	return false
}

func compareNestedConfigs(a, b *NestedConfig) bool {
	return a.Database.Host == b.Database.Host &&
		a.Database.Port == b.Database.Port &&
		a.Database.Password == b.Database.Password &&
		a.Server.Host == b.Server.Host &&
		a.Server.Port == b.Server.Port
}

func BenchmarkProcess(b *testing.B) {
	if err := os.Setenv("TEST_STRING", "benchmark"); err != nil {
		b.Fatalf("Failed to set env var: %v", err)
	}
	if err := os.Setenv("TEST_INT", "42"); err != nil {
		b.Fatalf("Failed to set env var: %v", err)
	}
	if err := os.Setenv("TEST_BOOL", "true"); err != nil {
		b.Fatalf("Failed to set env var: %v", err)
	}
	if err := os.Setenv("TEST_REQUIRED", "value"); err != nil {
		b.Fatalf("Failed to set env var: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("TEST_STRING"); err != nil {
			b.Logf("Failed to unset env var: %v", err)
		}
	}()
	defer func() {
		if err := os.Unsetenv("TEST_INT"); err != nil {
			b.Logf("Failed to unset env var: %v", err)
		}
	}()
	defer func() {
		if err := os.Unsetenv("TEST_BOOL"); err != nil {
			b.Logf("Failed to unset env var: %v", err)
		}
	}()
	defer func() {
		if err := os.Unsetenv("TEST_REQUIRED"); err != nil {
			b.Logf("Failed to unset env var: %v", err)
		}
	}()

	config := &TestConfig{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := Process("", config); err != nil {
			b.Fatalf("Process failed: %v", err)
		}
	}
}
