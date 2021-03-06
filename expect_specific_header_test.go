package hit_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/Eun/go-hit"
)

func TestExpectHeadersSpecificHeader_Len(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Header", "Hello")
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Post(s.URL),
		Expect().Header("X-Header").Len(5),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Expect().Header("X-Header").Len(0),
		),
		PtrStr(`"Hello" should have 0 item(s), but has 5`),
	)
}

func TestExpectHeadersSpecificHeader_Empty(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Header1", "Hello")
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Post(s.URL),
		Expect().Header("X-Header2").Empty(),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Expect().Header("X-Header1").Empty(),
		),
		PtrStr(`"Hello" should be empty, but has 5 item(s)`),
	)
}

func TestExpectSpecificHeader_Equal(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Header", "Hello")
		writer.Header().Set("X-Int", "3")
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Post(s.URL),
		Expect().Header("X-Header").Equal("Hello"),
		Expect().Header("X-Int").Equal(3),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Expect().Header("X-Header").Equal("Bye"),
		),
		PtrStr("Not equal"), PtrStr(`expected: "Bye"`), nil, nil, nil, nil, nil,
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Expect().Header("X-Int").Equal(1),
		),
		PtrStr("Not equal"), PtrStr("expected: 1"), nil, nil, nil, nil, nil,
	)
}

func TestExpectSpecificHeader_NotEqual(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Header", "Hello")
		writer.Header().Set("X-Int", "3")
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Post(s.URL),
		Expect().Header("X-Header").NotEqual("Bye"),
		Expect().Header("X-Int").NotEqual(1),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Expect().Header("X-Header").NotEqual("Hello"),
		),
		PtrStr(`should not be "Hello"`),
	)
	ExpectError(t,
		Do(
			Post(s.URL),
			Expect().Header("X-Int").NotEqual(3),
		),
		PtrStr(`should not be 3`),
	)
}

func TestExpectSpecificHeader_Contains(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Header", "Hello")
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Post(s.URL),
		Expect().Header("X-Header").Contains("Hello"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Expect().Header("X-Header").Contains("Bye"),
		),
		PtrStr(`"Hello" does not contain "Bye"`),
	)
}

func TestExpectSpecificHeader_NotContains(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Header", "Hello")
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Post(s.URL),
		Expect().Header("X-Header").NotContains("Bye"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Expect().Header("X-Header").NotContains("H"),
		),
		PtrStr(`"Hello" should not contain "H"`),
	)
}

func TestExpectSpecificHeader_OneOf(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Header", "Hello")
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Post(s.URL),
		Expect().Header("X-Header").OneOf("Hello", "World"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Expect().Header("X-Header").OneOf("Universe"),
		),
		PtrStr("[]interface {}{"), PtrStr(`"Universe",`), PtrStr(`} does not contain "Hello"`),
	)
}

func TestExpectSpecificHeader_NotOneOf(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Header", "Hello")
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	Test(t,
		Post(s.URL),
		Expect().Header("X-Header").NotOneOf("Universe"),
	)

	ExpectError(t,
		Do(
			Post(s.URL),
			Expect().Header("X-Header").NotOneOf("Hello", "World"),
		),
		nil, nil, nil, PtrStr(`} should not contain "Hello"`),
	)
}
