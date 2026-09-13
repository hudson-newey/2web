package javascript

import (
	"fmt"
	"hudson-newey/2web/src/constants"
)

const JsFunctionNamespace string = constants.CompilerNamespace + "func_"
const JsVarNamespace string = constants.CompilerNamespace + "var_"
const JsElementNamespace string = "data-" + constants.CompilerNamespace + "element_"

const ValueVar string = constants.CompilerNamespace + "value"

// IdAllocator allocates the unique runtime identifiers that wire a single
// page's compiled output together (DOM selectors, JavaScript variable names,
// and handler function names).
//
// Identifiers only need to be unique within a page: the emitted HTML and the
// emitted JavaScript for a page reference each other and are served together.
//
// Each page owns an allocator (see page.Page) instead of the compiler sharing
// one build wide counter. A build wide counter made the identifiers depend on
// which pages happened to be compiled before them (and in what order), so the
// compiled output of a page would change between builds even when none of its
// inputs changed. That breaks reproducible builds and forces browsers to
// re-download unchanged assets on every rebuild.
type IdAllocator struct {
	nextId int
}

func NewIdAllocator() IdAllocator {
	return IdAllocator{}
}

func (model *IdAllocator) CreateFunctionName() string {
	return fmt.Sprint(JsFunctionNamespace, model.reserve())
}

func (model *IdAllocator) CreateVariableName() string {
	return fmt.Sprint(JsVarNamespace, model.reserve())
}

func (model *IdAllocator) CreateElementName() string {
	return fmt.Sprint(JsElementNamespace, model.reserve())
}

func (model *IdAllocator) reserve() int {
	id := model.nextId
	model.nextId++
	return id
}
