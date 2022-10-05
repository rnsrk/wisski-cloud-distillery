package component

import (
	"reflect"
	"sync"

	"github.com/FAU-CDI/wisski-distillery/pkg/slicesx"
	"github.com/tkw1536/goprogram/lib/reflectx"
)

var metaCache sync.Map

// getMeta gets the component belonging to a component type
func getMeta[C Component]() meta {
	tp := reflectx.TypeOf[C]()

	// we already have a m => return it
	if m, ok := metaCache.Load(tp); ok {
		return m.(meta)
	}

	// create a new m
	var m meta
	m.init(tp)

	// store it in the cache
	metaCache.Store(tp, m)
	return m
}

// meta stores meta-information about a specific component
type meta struct {
	Name string       // the type name of this component
	Elem reflect.Type // the element type of the component

	CFields map[string]reflect.Type // fields with type C for which C implements component
	IFields map[string]reflect.Type // fields []I where I is an interface that implements component
}

// componentType is the type of components
var componentType = reflectx.TypeOf[Component]()

// init initializes this refklecttype
func (m *meta) init(tp reflect.Type) {
	if tp.Kind() != reflect.Pointer && tp.Elem().Kind() != reflect.Struct {
		panic("GetMeta: Type (" + tp.String() + ") must be backed by a pointer to slice")
	}

	m.Name = tp.Elem().PkgPath() + "." + tp.Elem().Name()
	m.Elem = tp.Elem()

	m.CFields = make(map[string]reflect.Type)
	m.IFields = make(map[string]reflect.Type)

	// fill the above variables, with a mapping of field name to struct
	count := m.Elem.NumField()
	for i := 0; i < count; i++ {
		field := m.Elem.Field(i)

		name := field.Name
		tp := field.Type

		switch {
		// field is a pointer to struct that implements a component
		case tp.Implements(componentType) && tp.Kind() == reflect.Pointer && tp.Elem().Kind() == reflect.Struct:
			m.CFields[name] = tp

		// field is []I, where I is an interface that implements component
		case tp.Kind() == reflect.Slice && tp.Elem().Kind() == reflect.Interface && tp.Elem().Implements(componentType):
			m.IFields[name] = tp.Elem()
		}
	}
}

// New creates a new ComponentDescription
func (m meta) New() Component {
	return reflect.New(m.Elem).Interface().(Component)
}

// NeedsInitComponent
func (m meta) NeedsInitComponent() bool {
	return len(m.CFields) > 0 || len(m.IFields) > 0
}

// InitComponent sets up the fields of the given instance of a component.
func (m meta) InitComponent(instance reflect.Value, all []Component) {
	elem := instance.Elem()

	// assign the component fields
	for field, eType := range m.CFields {
		c := slicesx.First(all, func(c Component) bool {
			return reflect.TypeOf(c).AssignableTo(eType)
		})

		field := elem.FieldByName(field)
		field.Set(reflect.ValueOf(c))
	}

	// assign the multi subtypes
	registryR := reflect.ValueOf(all)
	for field, eType := range m.IFields {
		cs := filterSubtype(registryR, eType)
		field := elem.FieldByName(field)
		field.Set(cs)
	}
}

// filterSubtype filters the slice of type []S into a slice of type []iface.
// S and I must be interface types.
func filterSubtype(sliceS reflect.Value, iface reflect.Type) reflect.Value {
	len := sliceS.Len()

	// convert each element
	result := reflect.MakeSlice(reflect.SliceOf(iface), 0, len)
	for i := 0; i < len; i++ {
		element := sliceS.Index(i)
		if element.Elem().Type().Implements(iface) {
			result = reflect.Append(result, element.Elem().Convert(iface))
		}
	}
	return result
}
