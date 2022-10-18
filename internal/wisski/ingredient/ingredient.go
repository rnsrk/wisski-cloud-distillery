package ingredient

import (
	"reflect"
	"strings"

	"github.com/FAU-CDI/wisski-distillery/internal/wisski/liquid"
)

// Ingredients represent a part of a WissKI instance.
// An Ingredient should be implemented as a pointer to a struct.
// Every ingredient must embed [Base] and should be initialized using [Init] inside a [lazy.Pool].
//
// By convention these are defined within their corresponding subpackage.
// This subpackage also contains all required resources.
type Ingredient interface {
	// Name returns the name of this ingredient
	// Name should be implemented by the [Base] struct.
	Name() string

	// getBase returns the underlying Base object of this Ingredient.
	// It is used internally during initialization
	getBase() *Base
}

// Base is embedded into every Ingredient
type Base struct {
	name           string // name is the name of this ingredient
	*liquid.Liquid        // the underlying liquid
}

//lint:ignore U1000 used to implement the private methods of [Component]
func (cb *Base) getBase() *Base {
	return cb
}

// Init initializes a new Ingredient.
// Init is only intended to be used within a lazy.Pool[Ingredient,*Liquid].
func Init(ingredient Ingredient, liquid *liquid.Liquid) Ingredient {
	base := ingredient.getBase() // pointer to a struct
	base.Liquid = liquid
	base.name = strings.ToLower(reflect.TypeOf(ingredient).Elem().Name())
	return ingredient
}

func (cb Base) Name() string {
	return cb.name
}
