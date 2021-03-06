package hit

import (
	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/internal"
	"golang.org/x/xerrors"
)

// IClearSend provides a clear functionality to remove previous steps from running in the Send() scope
type IClearSend interface {
	IStep
	// Body removes all previous Send().Body() steps and all steps chained to Send().Body() e.g. Send().Body().Interface("Hello World").
	//
	// If you specify an argument it will only remove the Send().Body() steps matching that argument.
	//
	// Usage:
	//     Clear().Send().Body()                      // will remove all Send().Body() steps and all chained steps to Send() e.g. Send().Body("Hello World")
	//     Clear().Send().Body("Hello World")         // will remove all Send().Body("Hello World") steps
	//     Clear().Send().Body().Interface()              // will remove all Send().Body().Interface() steps
	//     Clear().Send().Body().Interface("Hello World") // will remove all Send().Body().Interface("Hello World") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Send().Body("Hello Earth"),
	//         Clear().Send().Body(),
	//         Send().Body("Hello World"),
	//     )
	Body(value ...interface{}) IClearSendBody

	// Interface removes all previous Send().Interface() steps.
	//
	// If you specify an argument it will only remove the Send().Interface() steps matching that argument.
	//
	// Usage:
	//     Clear().Send().Interface()              // will remove all Send().Interface() steps
	//     Clear().Send().Interface("Hello World") // will remove all Send().Interface("Hello World") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Send().Interface("Hello Earth"),
	//         Clear().Send().Interface(),
	//         Send().Interface("Hello World"),
	//     )
	Interface(value ...interface{}) IStep

	// JSON removes all previous Send().JSON() steps.
	//
	// If you specify an argument it will only remove the Send().JSON() steps matching that argument.
	//
	// Usage:
	//     Clear().Send().JSON()                                      // will remove all Send().JSON() steps
	//     Clear().Send().JSON(map[string]interface{}{"Name": "Joe"}) // will remove all Send().JSON("Hello World") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Send().JSON(map[string]interface{}{"Name": "Joe"}),
	//         Clear().Send().JSON(),
	//         Send().JSON(map[string]interface{}{"Name": "Alice"}),
	//     )
	JSON(value ...interface{}) IStep

	// Header removes all previous Send().Header() steps.
	//
	// If you specify an argument it will only remove the Send().Header() steps matching that argument.
	//
	// Usage:
	//     Clear().Send().Header()                                   // will remove all Send().Header() steps
	//     Clear().Send().Header("Content-Type")                     // will remove all Send().Header("Content-Type", ...) step
	//     Clear().Send().Header("Content-Type", "application/json") // will remove all Send().Header("Content-Type", "application/json") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Send().Header("Content-Type", "application/xml"),
	//         Clear().Send().Header("Content-Type"),
	//         Send().Header("Content-Type", "application/json"),
	//     )
	Header(values ...interface{}) IStep

	// Custom removes all previous Send().Custom() steps.
	//
	// If you specify an argument it will only remove the Send().Custom() steps matching that argument.
	//
	// Usage:
	//     Clear().Send().Custom(fn) // will remove all Send().Custom(fn) steps
	//     Clear().Send().Custom()   // will remove all Send().Custom() steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Send().Custom(func(hit Hit) {
	//             hit.Request().Body().SetString("Hello Earth")
	//         }),
	//         Clear().Send().Custom(),
	//         Send().Custom(func(hit Hit) {
	//             hit.Request().Body().SetString("Hello World")
	//         }),
	//     )
	Custom(fn ...Callback) IStep
}

type clearSend struct {
	cleanPath clearPath
	trace     *errortrace.ErrorTrace
}

func newClearSend(clearPath clearPath, params []interface{}) IClearSend {
	if _, ok := internal.GetLastArgument(params); ok {
		// this runs if we called Clear().Send(something)
		return &finalClearSend{
			removeStep(clearPath),
			"only usable with Clear().Send() not with Clear().Send(value)",
		}
	}
	return &clearSend{
		cleanPath: clearPath,
		trace:     ett.Prepare(),
	}
}

func (*clearSend) when() StepTime {
	return CleanStep
}

func (snd *clearSend) exec(hit Hit) error {
	// this runs if we called Clear().Send()
	if err := removeSteps(hit, snd.clearPath()); err != nil {
		return snd.trace.Format(hit.Description(), err.Error())
	}
	return nil
}

func (snd *clearSend) clearPath() clearPath {
	return snd.cleanPath
}

func (snd *clearSend) Body(value ...interface{}) IClearSendBody {
	return newClearSendBody(snd.clearPath().Push("Body", value), value)
}

func (snd *clearSend) Interface(value ...interface{}) IStep {
	return removeStep(snd.clearPath().Push("Interface", value))
}

// custom can be used to send a custom behaviour
func (snd *clearSend) Custom(fn ...Callback) IStep {
	args := make([]interface{}, len(fn))
	for i := range fn {
		args[i] = fn[i]
	}
	return removeStep(snd.clearPath().Push("Custom", args))
}

// JSON sets the body to the specific data (shortcut for Body().JSON()
func (snd *clearSend) JSON(value ...interface{}) IStep {
	return removeStep(snd.clearPath().Push("JSON", value))
}

func (snd *clearSend) Header(values ...interface{}) IStep {
	return removeStep(snd.clearPath().Push("Header", values))
}

type finalClearSend struct {
	IStep
	message string
}

func (snd *finalClearSend) fail() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New(snd.message)
		},
	}
}

func (snd *finalClearSend) Body(...interface{}) IClearSendBody {
	return &finalClearSendBody{
		snd.fail(),
		snd.message,
	}
}

func (snd *finalClearSend) Custom(...Callback) IStep {
	return snd.fail()
}

func (snd *finalClearSend) JSON(...interface{}) IStep {
	return snd.fail()
}

func (snd *finalClearSend) Header(...interface{}) IStep {
	return snd.fail()
}

func (snd *finalClearSend) Interface(...interface{}) IStep {
	return snd.fail()
}
