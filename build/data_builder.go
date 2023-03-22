package build

import "reflect"

const (
	builder_cap = 7
)

func Builder() DataBuilder {
	return &builder{data: make(map[string]any, builder_cap)}
}

type DataBuilder interface {
	Field(key string, value any) DataBuilder
	FileValue(key string, get func(val any) any) DataBuilder
	FieldValid(key string, value any) DataBuilder
	Reset() DataBuilder
	Data() map[string]any
}

type builder struct {
	data map[string]any
}

func (builder *builder) Field(key string, value any) DataBuilder {
	builder.data[key] = value
	return builder
}

func (builder *builder) FileValue(key string, get func(val any) any) DataBuilder {
	v := builder.data[key]
	val := get(v)
	if val != nil {
		builder.data[key] = val
	}
	return builder
}

// not nil, not empty string ,not empty slice|arr|chan|map
func (builder *builder) FieldValid(key string, value any) DataBuilder {
	if value != nil {
		rval := reflect.ValueOf(value)
		if rval.Kind() == reflect.Ptr {
			rval = rval.Elem()
		}
		kind := rval.Kind()
		switch kind {
		case reflect.String:
		case reflect.Array:
		case reflect.Chan:
		case reflect.Map:
			if rval.Len() == 0 {
				return builder
			}
		default:
		}
		builder.data[key] = value
	}
	return builder
}

func (builder *builder) Reset() DataBuilder {
	builder.data = make(map[string]any, builder_cap)
	return builder
}

func (builder *builder) Data() map[string]any {
	return builder.data
}
