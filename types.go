package phpserialize

import (
	"fmt"
	"github.com/vmihailenco/tagparser"
	"reflect"
	"sync"
)

var structs = newStructCache()

type structCache struct {
	m sync.Map
}

type structCacheKey struct {
	tag string
	typ reflect.Type
}

func newStructCache() *structCache {
	return new(structCache)
}

func (m *structCache) Fields(typ reflect.Type, tag string) *fields {
	key := structCacheKey{tag: tag, typ: typ}

	if v, ok := m.m.Load(key); ok {
		return v.(*fields)
	}

	fs := getFields(typ, tag)
	m.m.Store(key, fs)

	return fs
}

type field struct {
	name  string
	index []int
	// omitEmpty bool
	// encoder   encoderFunc
	decoder decoderFunc
}

func newFields(typ reflect.Type) *fields {
	return &fields{
		Type: typ,
		Map:  make(map[string]*field, typ.NumField()),
		List: make([]*field, 0, typ.NumField()),
	}
}

var (
	defaultStructTag = `php`
)

func getFields(typ reflect.Type, fallbackTag string) *fields {
	fs := newFields(typ)

	// var omitEmpty bool
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)

		tagStr := f.Tag.Get(defaultStructTag)
		if tagStr == "" && fallbackTag != "" {
			tagStr = f.Tag.Get(fallbackTag)
		}

		tag := tagparser.Parse(tagStr)
		if tag.Name == "-" {
			continue
		}

		field := &field{
			name:  tag.Name,
			index: f.Index,
			// omitEmpty: omitEmpty || tag.HasOption("omitempty"),
		}
		field.decoder = getDecoder(f.Type)

		if field.name == "" {
			field.name = f.Name
		}

		fs.Add(field)
	}

	return fs
}

type fields struct {
	Type reflect.Type
	Map  map[string]*field
	List []*field
	// AsArray bool

	// hasOmitEmpty bool
}

func (fs *fields) Add(field *field) {
	// fs.warnIfFieldExists(field.name)
	fs.Map[field.name] = field
	fs.List = append(fs.List, field)
	//if field.omitEmpty {
	//	fs.hasOmitEmpty = true
	//}
}

func (f *field) DecodeValue(d *Decoder, strct reflect.Value) error {
	v := fieldByIndexAlloc(strct, f.index)
	if f.decoder == nil {
		return fmt.Errorf(`phpserialize: could not find decoder for field %s`, f.name)
	}
	return f.decoder(d, v)
}

func fieldByIndexAlloc(v reflect.Value, index []int) reflect.Value {
	if len(index) == 1 {
		return v.Field(index[0])
	}

	/*
		for i, idx := range index {
			if i > 0 {
				var ok bool
				v, ok = indirectNil(v)
				if !ok {
					return v
				}
			}
			v = v.Field(idx)
		}

		return v
	*/

	panic(`unsupported`)
}
