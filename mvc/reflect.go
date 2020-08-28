package mvc

import (
	"reflect"

	"github.com/kataras/iris/v12/context"
)

var baseControllerTyp = reflect.TypeOf((*BaseController)(nil)).Elem()

func isBaseController(ctrlTyp reflect.Type) bool {
	return ctrlTyp.Implements(baseControllerTyp)
}

// indirectType returns the value of a pointer-type "typ".
// If "typ" is a pointer, array, chan, map or slice it returns its Elem,
// otherwise returns the typ as it's.
func indirectType(typ reflect.Type) reflect.Type {
	switch typ.Kind() {
	case reflect.Ptr, reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return typ.Elem()
	}
	return typ
}

func getSourceFileLine(ctrlType reflect.Type, m reflect.Method) (string, int) { // used for debug logs.
	sourceFileName, sourceLineNumber := context.HandlerFileLineRel(m.Func)
	if sourceFileName == "<autogenerated>" {
		elem := indirectType(ctrlType)

		for i, n := 0, elem.NumField(); i < n; i++ {
			if f := elem.Field(i); f.Anonymous {
				typ := indirectType(f.Type)
				if typ.Kind() != reflect.Struct {
					continue // field is not a struct.
				}

				// why we do that?
				// because if the element is not Ptr
				// then it's probably used as:
				// type ctrl {
				//   BaseCtrl
				// }
				// but BaseCtrl has not the method, *BaseCtrl does:
				// (c *BaseCtrl) HandleHTTPError(...)
				// so we are creating a new temporarly value ptr of that type
				// and searching inside it for the method instead.
				typ = reflect.New(typ).Type()

				if embeddedMethod, ok := typ.MethodByName(m.Name); ok {
					sourceFileName, sourceLineNumber = context.HandlerFileLineRel(embeddedMethod.Func)
				}
			}
		}
	}

	return sourceFileName, sourceLineNumber
}
