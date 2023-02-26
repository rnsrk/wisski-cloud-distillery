package wisski

import (
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/liquid"
	"github.com/tkw1536/pkglib/collection"
	"github.com/tkw1536/pkglib/lifetime"
)

//
//  ==== init ====
//

func (wisski *WissKI) init() {
	wisski.lifetimeInit.Do(func() {
		wisski.lifetime.Init = ingredient.Init
		lifetime.RegisterGroup[ingredient.WissKIFetcher](&wisski.lifetime)
	})
}

//
//  ==== registration ====
//

// manual initializes a component from the provided distillery.
func manual[I ingredient.Ingredient](init func(ingredient I)) initFunc {
	return func(context ctx) ingredient.Ingredient {
		return lifetime.Make(context, init)
	}
}

// use is like r, but does not provided additional initialization
func auto[I ingredient.Ingredient](context ctx) ingredient.Ingredient {
	return lifetime.Make[ingredient.Ingredient, I](context, nil)
}

// register returns all components of the distillery
func (wisski *WissKI) register(context ctx) []ingredient.Ingredient {
	return collection.MapSlice(
		wisski.allIngredients(),
		func(f initFunc) ingredient.Ingredient {
			return f(context)
		},
	)
}

// ctx is a context for component initialization
type ctx = *lifetime.InjectorContext[ingredient.Ingredient]

//
//  ==== export ====
//

// export is a convenience function to export a single component
func export[I ingredient.Ingredient](wisski *WissKI) I {
	wisski.init()
	return lifetime.ExportComponent[ingredient.Ingredient, *liquid.Liquid, I](&wisski.lifetime, &wisski.Liquid, wisski.register)
}

//lint:ignore U1000 for future use
func exportAll[I ingredient.Ingredient](wisski *WissKI) []I {
	wisski.init()
	return lifetime.ExportComponents[ingredient.Ingredient, *liquid.Liquid, I](&wisski.lifetime, &wisski.Liquid, wisski.register)
}

type initFunc = func(context ctx) ingredient.Ingredient
