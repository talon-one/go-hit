package hit

import (
	"net/http"

	"github.com/Eun/go-convert"
	"github.com/Eun/go-hit/errortrace"
	"github.com/pkg/errors"
)

type ExpectHeadersCallback func(headers *http.Header)

type expectHeaders struct {
	Hit
	expect         *defaultExpect
	specificHeader string
}

func newExpectHeaders(expect *defaultExpect, name string) *expectHeaders {
	return &expectHeaders{
		Hit:            expect.Hit,
		expect:         expect,
		specificHeader: name,
	}
}

func (hdr *expectHeaders) value(headers http.Header) interface{} {
	if hdr.specificHeader == "" {
		return headers
	}
	return headers.Get(hdr.specificHeader)
}

// Contains checks if the specified header is present or
// (if header was already specified) the header value contains the specified value
// Examples:
//           Expect().Headers().Contains("Content-Type")
//           Expect().Headers("Content-Type").Contains("application")
func (hdr *expectHeaders) Contains(v string) Hit {
	et := errortrace.Prepare()
	return hdr.expect.Custom(func(hit Hit) {
		et.Panic.Contains(hit.T(), hdr.value(hit.Response().Header), v)
	})
}

// OneOf checks if the value is one of the specified values
// Example:
//           Expect().Headers("Content-Type").OneOf("application/json", "text/x-json")
func (hdr *expectHeaders) OneOf(values ...interface{}) Hit {
	if hdr.specificHeader == "" {
		errortrace.Panic.FailNow(hdr.T(), errors.New("OneOf can only be used if a header was already specified"))
	}
	et := errortrace.Prepare()
	return hdr.expect.Custom(func(hit Hit) {
		et.Panic.Contains(hit.T(), values, hdr.value(hit.Response().Header))
	})
}

// Empty checks if the headers are empty or
// (if header was already specified) the header value is empty
// Examples:
//           Expect().Headers().Empty()
//           Expect().Headers("Content-Type").Empty()
func (hdr *expectHeaders) Empty() Hit {
	et := errortrace.Prepare()
	return hdr.expect.Custom(func(hit Hit) {
		et.Panic.Empty(hit.T(), hdr.value(hit.Response().Header))
	})
}

// Len checks if the amount of headers are equal to the specified size or
// (if header was already specified) the length of the header value is equal to the specified size
// Examples:
//           Expect().Headers().Len(0)
//           Expect().Headers("Content-Type").Len(16)
func (hdr *expectHeaders) Len(size int) Hit {
	et := errortrace.Prepare()
	return hdr.expect.Custom(func(hit Hit) {
		et.Panic.Len(hit.T(), hdr.value(hit.Response().Header), size)
	})
}

// Equal checks if the headers are equal to the specified one or
// (if header was already specified) the header value is equal to the specified value
// Examples:
//           Expect().Headers().Equal(map[string]string{"Content-Type": "application/json"})
//           Expect().Headers("Content-Type").Equal("application/json")
func (hdr *expectHeaders) Equal(v interface{}) Hit {
	et := errortrace.Prepare()
	return hdr.expect.Custom(func(hit Hit) {
		compareData, err := converter.Convert(hdr.value(hit.Response().Header), v, convert.Options.ConvertEmbeddedStructToParentType())
		et.Panic.NoError(hit.T(), err)
		et.Panic.Equal(hit.T(), v, compareData)
	})
}

// Get a specific header
// Examples:
//           Expect().Headers().Get("Content-Type").Equal("application/json")
//           Expect().Headers().Get("Content-Type").Contains("json")
func (hdr *expectHeaders) Get(name string) *expectHeaders {
	if hdr.specificHeader != "" {
		errortrace.Panic.FailNow(hdr.T(), errors.New("Get can only be used if no header was already specified"))
	}
	return &expectHeaders{
		Hit:            hdr.Hit,
		expect:         hdr.expect,
		specificHeader: name,
	}
}