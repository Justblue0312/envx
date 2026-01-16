package envx

import (
	"encoding"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidSpecification = errors.New("specification must be a struct pointer")

var gatherRegexp = regexp.MustCompile("([^A-Z]+|[A-Z]+[^A-Z]+|[A-Z]+)")
var acronymRegexp = regexp.MustCompile("([A-Z]+)([A-Z][^A-Z]+)")

type ParseError struct {
	KeyName   string
	FieldName string
	TypeName  string
	Value     string
	Err       error
}

type Decoder interface {
	Decode(value string) error
}

type Setter interface {
	Set(value string) error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("envx.Process: assigning %[1]s to %[2]s: converting '%[3]s' to type %[4]s. details: %[5]s", e.KeyName, e.FieldName, e.Value, e.TypeName, e.Err)
}

type varInfo struct {
	Name  string
	Alt   string
	Key   string
	Field reflect.Value
	Tags  reflect.StructTag
}

func gatherInfo(prefix string, spec interface{}) ([]varInfo, error) {
	s := reflect.ValueOf(spec)

	if s.Kind() != reflect.Ptr || s.Elem().Kind() != reflect.Struct {
		return nil, ErrInvalidSpecification
	}
	s = s.Elem()
	typeOfSpec := s.Type()

	infos := make([]varInfo, 0, s.NumField())
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		fieldType := typeOfSpec.Field(i)

		if !field.CanSet() || isTrue(fieldType.Tag.Get("ignored")) {
			continue
		}

		for field.Kind() == reflect.Ptr {
			if field.IsNil() {
				if field.Type().Elem().Kind() != reflect.Struct {
					break
				}
				field.Set(reflect.New(field.Type().Elem()))
			}
			field = field.Elem()
		}

		info := varInfo{
			Name:  fieldType.Name,
			Field: field,
			Tags:  fieldType.Tag,
			Alt:   strings.ToUpper(fieldType.Tag.Get("envconfig")),
		}

		info.Key = info.Name

		if isTrue(fieldType.Tag.Get("split_words")) {
			info.Key = toSnakeCase(info.Name)
		}

		if info.Alt != "" {
			info.Key = info.Alt
		}

		separator := "_"
		if isTrue(fieldType.Tag.Get("nested")) {
			separator = "__"
		}

		if prefix != "" {
			info.Key = fmt.Sprintf("%s%s%s", prefix, separator, info.Key)
		}
		info.Key = strings.ToUpper(info.Key)
		infos = append(infos, info)

		if field.Kind() == reflect.Struct {
			if decoderFrom(field) == nil && setterFrom(field) == nil &&
				textUnmarshaler(field) == nil && binaryUnmarshaler(field) == nil {
				innerPrefix := prefix
				if !fieldType.Anonymous {
					innerPrefix = info.Key
				}

				embeddedPtr := field.Addr().Interface()
				embeddedInfos, err := gatherInfo(innerPrefix, embeddedPtr)
				if err != nil {
					return nil, err
				}
				infos = append(infos[:len(infos)-1], embeddedInfos...)
				continue
			}
		}
	}
	return infos, nil
}

func toSnakeCase(s string) string {
	words := gatherRegexp.FindAllStringSubmatch(s, -1)
	if len(words) == 0 {
		return s
	}

	var result []string
	for _, word := range words {
		if m := acronymRegexp.FindStringSubmatch(word[0]); len(m) == 3 {
			result = append(result, m[1], m[2])
		} else {
			result = append(result, word[0])
		}
	}

	return strings.Join(result, "_")
}

func CheckDisallowed(prefix string, spec interface{}) error {
	infos, err := gatherInfo(prefix, spec)
	if err != nil {
		return err
	}

	vars := make(map[string]struct{})
	for _, info := range infos {
		vars[info.Key] = struct{}{}
	}

	if prefix != "" {
		prefix = strings.ToUpper(prefix) + "_"
	}

	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, prefix) {
			continue
		}
		v := strings.SplitN(env, "=", 2)[0]
		if _, found := vars[v]; !found {
			return fmt.Errorf("unknown environment variable %s", v)
		}
	}

	return nil
}

func Process(prefix string, spec interface{}) error {
	infos, err := gatherInfo(prefix, spec)
	if err != nil {
		return err
	}

	for _, info := range infos {
		value, ok := lookupEnv(info.Key)
		if !ok && info.Alt != "" {
			value, ok = lookupEnv(info.Alt)
		}

		if !ok {
			value = tryNestedKeys(info.Key)
			if value != "" {
				ok = true
			}
		}

		def := info.Tags.Get("default")
		if def != "" && !ok {
			value = def
		}

		req := info.Tags.Get("required")
		if !ok && def == "" {
			if isTrue(req) {
				key := info.Key
				if info.Alt != "" {
					key = info.Alt
				}
				return fmt.Errorf("required key %s missing value", key)
			}
			continue
		}

		err = processField(value, info.Field)
		if err != nil {
			return &ParseError{
				KeyName:   info.Key,
				FieldName: info.Name,
				TypeName:  info.Field.Type().String(),
				Value:     value,
				Err:       err,
			}
		}
	}

	return nil
}

func MustProcess(prefix string, spec interface{}) {
	if err := Process(prefix, spec); err != nil {
		panic(err)
	}
}

func processField(value string, field reflect.Value) error {
	typ := field.Type()

	decoder := decoderFrom(field)
	if decoder != nil {
		return decoder.Decode(value)
	}

	setter := setterFrom(field)
	if setter != nil {
		return setter.Set(value)
	}

	if t := textUnmarshaler(field); t != nil {
		return t.UnmarshalText([]byte(value))
	}

	if b := binaryUnmarshaler(field); b != nil {
		return b.UnmarshalBinary([]byte(value))
	}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		if field.IsNil() {
			field.Set(reflect.New(typ))
		}
		field = field.Elem()
	}

	switch typ.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var (
			val int64
			err error
		)
		if field.Kind() == reflect.Int64 && typ.PkgPath() == "time" && typ.Name() == "Duration" {
			var d time.Duration
			d, err = time.ParseDuration(value)
			val = int64(d)
		} else {
			val, err = strconv.ParseInt(value, 0, typ.Bits())
		}
		if err != nil {
			return err
		}
		field.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(value, 0, typ.Bits())
		if err != nil {
			return err
		}
		field.SetUint(val)
	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(value, typ.Bits())
		if err != nil {
			return err
		}
		field.SetFloat(val)
	case reflect.Slice:
		sl := reflect.MakeSlice(typ, 0, 0)
		if typ.Elem().Kind() == reflect.Uint8 {
			sl = reflect.ValueOf([]byte(value))
		} else if strings.TrimSpace(value) != "" {
			vals := strings.Split(value, ",")
			sl = reflect.MakeSlice(typ, len(vals), len(vals))
			for i, val := range vals {
				err := processField(val, sl.Index(i))
				if err != nil {
					return err
				}
			}
		}
		field.Set(sl)
	case reflect.Map:
		mp := reflect.MakeMap(typ)
		if strings.TrimSpace(value) != "" {
			pairs := strings.Split(value, ",")
			for _, pair := range pairs {
				kvpair := strings.Split(pair, ":")
				if len(kvpair) != 2 {
					return fmt.Errorf("invalid map item: %q", pair)
				}
				k := reflect.New(typ.Key()).Elem()
				err := processField(kvpair[0], k)
				if err != nil {
					return err
				}
				v := reflect.New(typ.Elem()).Elem()
				err = processField(kvpair[1], v)
				if err != nil {
					return err
				}
				mp.SetMapIndex(k, v)
			}
		}
		field.Set(mp)
	}

	return nil
}

func interfaceFrom(field reflect.Value, fn func(interface{}, *bool)) {
	if !field.CanInterface() {
		return
	}
	var ok bool
	fn(field.Interface(), &ok)
	if !ok && field.CanAddr() {
		fn(field.Addr().Interface(), &ok)
	}
}

func decoderFrom(field reflect.Value) (d Decoder) {
	interfaceFrom(field, func(v interface{}, ok *bool) { d, *ok = v.(Decoder) })
	return d
}

func setterFrom(field reflect.Value) (s Setter) {
	interfaceFrom(field, func(v interface{}, ok *bool) { s, *ok = v.(Setter) })
	return s
}

func textUnmarshaler(field reflect.Value) (t encoding.TextUnmarshaler) {
	interfaceFrom(field, func(v interface{}, ok *bool) { t, *ok = v.(encoding.TextUnmarshaler) })
	return t
}

func binaryUnmarshaler(field reflect.Value) (b encoding.BinaryUnmarshaler) {
	interfaceFrom(field, func(v interface{}, ok *bool) { b, *ok = v.(encoding.BinaryUnmarshaler) })
	return b
}

func isTrue(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}
