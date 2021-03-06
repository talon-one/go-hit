package hit

import (
	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/internal"
	"golang.org/x/xerrors"
)

// IClearExpectBody provides a clear functionality to remove previous steps from running in the Expect().Body().JSON() scope
type IClearExpectBodyJSON interface {
	IStep
	// Equal removes all previous Expect().Body().JSON().Equal() steps.
	//
	// If you specify an argument it will only remove the Expect().Body().Equal() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Body().JSON().Equal()              // will remove all Expect().Body().JSON().Equal() steps
	//     Clear().Expect().Body().JSON().Equal("Name")        // will remove all Expect().Body().JSON().Equal("Name", ...) steps
	//     Clear().Expect().Body().JSON().Equal("Name", "Joe") // will remove all Expect().Body().JSON().Equal("Name", "Joe") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Body().JSON().Equal("Name", "Joe"),
	//         Expect().Body().JSON().Equal("Id", 10),
	//         Clear().Expect().Body().JSON().Equal("Name"),
	//         Clear().Expect().Body().JSON().Equal("Id", 10),
	//         Expect().Body().JSON().Equal("Name", "Alice"),
	//     )
	Equal(value ...interface{}) IStep

	// NotEqual removes all previous Expect().Body().JSON().NotEqual() steps.
	//
	// If you specify an argument it will only remove the Expect().Body().NotEqual() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Body().JSON().NotEqual()              // will remove all Expect().Body().JSON().NotEqual() steps
	//     Clear().Expect().Body().JSON().NotEqual("Name")        // will remove all Expect().Body().JSON().NotEqual("Name") steps
	//     Clear().Expect().Body().JSON().NotEqual("Name", "Joe") // will remove all Expect().Body().JSON().NotEqual("Name", "Joe") steps
	//
	// Example:
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Body().JSON().NotEqual("Name", "Joe"),
	//         Expect().Body().JSON().NotEqual("Id", 10),
	//         Clear().Expect().Body().JSON().NotEqual("Name"),
	//         Clear().Expect().Body().JSON().NotEqual("Id", 10),
	//         Expect().Body().JSON().NotEqual("Name", "Alice"),
	//     )
	NotEqual(value ...interface{}) IStep

	// Contains removes all previous Expect().Body().JSON().Contains() steps.
	//
	// If you specify an argument it will only remove the Expect().Body().Contains() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Body().JSON().Contains()              // will remove all Expect().Body().JSON().Contains() steps
	//     Clear().Expect().Body().JSON().Contains("Name")        // will remove all Expect().Body().JSON().Contains("Name") steps
	//     Clear().Expect().Body().JSON().Contains("Name", "Joe") // will remove all Expect().Body().JSON().Contains("Name", "Joe") steps
	//
	// Example:
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Body().JSON().Contains("Name", "Joe"),
	//         Expect().Body().JSON().Contains("Id", 10),
	//         Clear().Expect().Body().JSON().Contains("Name"),
	//         Clear().Expect().Body().JSON().Contains("Id", 10),
	//         Expect().Body().JSON().Contains("Name", "Alice"),
	//     )
	Contains(value ...interface{}) IStep

	// NotContains removes all previous Expect().Body().JSON().NotContains() steps.
	//
	// If you specify an argument it will only remove the Expect().Body().NotContains() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect().Body().JSON().NotContains()              // will remove all Expect().Body().JSON().NotContains() steps
	//     Clear().Expect().Body().JSON().NotContains("Name")        // will remove all Expect().Body().JSON().NotContains("Name") steps
	//     Clear().Expect().Body().JSON().NotContains("Name", "Joe") // will remove all Expect().Body().JSON().NotContains("Name", "Joe") steps
	//
	// Example:
	//     Do(
	//         Post("https://example.com"),
	//         Expect().Body().JSON().NotContains("Name", "Joe"),
	//         Expect().Body().JSON().NotContains("Id", 10),
	//         Clear().Expect().Body().JSON().NotContains("Name"),
	//         Clear().Expect().Body().JSON().NotContains("Id", 10),
	//         Expect().Body().JSON().NotContains("Name", "Alice"),
	//     )
	NotContains(value ...interface{}) IStep
}

type clearExpectBodyJSON struct {
	clearExpectBody IClearExpectBody
	cleanPath       clearPath
	trace           *errortrace.ErrorTrace
}

func newClearExpectBodyJSON(body IClearExpectBody, cleanPath clearPath, params []interface{}) IClearExpectBodyJSON {
	if _, ok := internal.GetLastArgument(params); ok {
		// this runs if we called Clear().Expect().Body().JSON(something)
		return &finalClearExpectBodyJSON{
			removeStep(cleanPath),
			"only usable with Clear().Expect().Body().JSON() not with Clear().Expect().Body().JSON(value)",
		}
	}
	return &clearExpectBodyJSON{
		clearExpectBody: body,
		cleanPath:       cleanPath,
		trace:           ett.Prepare(),
	}
}

func (jsn *clearExpectBodyJSON) when() StepTime {
	return CleanStep
}

func (jsn *clearExpectBodyJSON) exec(hit Hit) error {
	// this runs if we called Clear().Expect().Body().JSON()
	if err := removeSteps(hit, jsn.clearPath()); err != nil {
		return jsn.trace.Format(hit.Description(), err.Error())
	}
	return nil
}

func (jsn *clearExpectBodyJSON) clearPath() clearPath {
	return jsn.cleanPath
}

func (jsn *clearExpectBodyJSON) Equal(value ...interface{}) IStep {
	return removeStep(jsn.clearPath().Push("Equal", value))
}

func (jsn *clearExpectBodyJSON) NotEqual(value ...interface{}) IStep {
	return removeStep(jsn.clearPath().Push("NotEqual", value))
}

func (jsn *clearExpectBodyJSON) Contains(value ...interface{}) IStep {
	return removeStep(jsn.clearPath().Push("Contains", value))
}

func (jsn *clearExpectBodyJSON) NotContains(value ...interface{}) IStep {
	return removeStep(jsn.clearPath().Push("NotContains", value))
}

type finalClearExpectBodyJSON struct {
	IStep
	message string
}

func (jsn *finalClearExpectBodyJSON) fail() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New(jsn.message)
		},
	}
}

func (jsn *finalClearExpectBodyJSON) Equal(...interface{}) IStep {
	return jsn.fail()
}

func (jsn *finalClearExpectBodyJSON) NotEqual(...interface{}) IStep {
	return jsn.fail()
}

func (jsn *finalClearExpectBodyJSON) Contains(...interface{}) IStep {
	return jsn.fail()
}

func (jsn *finalClearExpectBodyJSON) NotContains(...interface{}) IStep {
	return jsn.fail()
}
