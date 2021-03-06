package hit_test

import (
	"testing"

	"errors"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestExpect_Custom(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Expect().Custom(func(hit Hit) {
			require.Equal(t, "Hello World", hit.Response().Body().String())
		}),
	)
}

func TestExpect_Double(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	Test(t,
		Post(s.URL),
		Send().Body("Hello World"),
		Expect().Body().Equal(`Hello World`),
		Expect().Body().Equal(`Hello World`),
	)
}

func TestExpect(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("func", func(t *testing.T) {
		t.Run("with correct parameter (using Response)", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body("Hello World"),
					Expect(func(hit Hit) {
						if hit.Response().Body().String() != "Hello Universe" {
							panic("Not equal")
						}
					}),
				),
				PtrStr("Not equal"),
			)
		})
		t.Run("with correct parameter (using Hit)", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body("Hello World"),
					Expect(func(hit Hit) {
						hit.MustDo(Expect("Hello Universe"))
					})),
				PtrStr("Not equal"), nil, nil, nil, nil, nil, nil,
			)
		})
		t.Run("with correct parameter (using Hit) and error", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body("Hello World"),
					Expect(func(hit Hit) error {
						return hit.Do(Expect("Hello Universe"))
					})),
				PtrStr("Not equal"), PtrStr(`expected: "Hello Universe"`), PtrStr(`actual: "Hello World"`), nil, nil, nil, nil,
			)
		})
		t.Run("with correct parameter (using Hit) and error (return an error)", func(t *testing.T) {
			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body("Hello World"),
					Expect(func(hit Hit) error {
						return errors.New("whoops")
					})),
				PtrStr("whoops"),
			)
		})
		t.Run("with invalid parameter", func(t *testing.T) {
			calledFunc := false
			Test(t,
				Post(s.URL),
				Send().Body("Hello World"),
				Expect(func() {
					calledFunc = true
				}),
			)
			require.True(t, calledFunc)
		})
	})

	t.Run("body", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect("Hello World"),
		)
	})
}

func TestExpect_DeepFunc(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	calledFunc := false
	ExpectError(t,
		Do(
			Post(s.URL),
			Send().Body("Hello World"),
			Expect(func(h1 Hit) {
				h1.MustDo(Expect(func(h2 Hit) {
					h2.MustDo(Expect(func(h3 Hit) {
						calledFunc = true
						h3.MustDo(Expect().Body().Equal("Hello Universe"))
					}))
				}))
			}),
		),
		PtrStr("Not equal"), nil, nil, nil, nil, nil, nil,
	)
	require.True(t, calledFunc)
}

func TestExpect_Final(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("Expect(value).Body()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect("Data").Body()),
			PtrStr("only usable with Expect() not with Expect(value)"),
		)
	})

	t.Run("Expect(value).Body().JSON()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect("Data").Body().JSON()),
			PtrStr("only usable with Expect() not with Expect(value)"),
		)
	})

	t.Run("Expect(value).Interface()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect("Data").Interface(nil)),
			PtrStr("only usable with Expect() not with Expect(value)"),
		)
	})

	t.Run("Expect(value).Header().Contains()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect("Data").Header().Contains(nil)),
			PtrStr("only usable with Expect() not with Expect(value)"),
		)
	})

	t.Run("Expect(value).Header().NotContains()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect("Data").Header().NotContains(nil)),
			PtrStr("only usable with Expect() not with Expect(value)"),
		)
	})

	t.Run("Expect(value).Header().OneOf()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect("Data").Header().OneOf()),
			PtrStr("only usable with Expect() not with Expect(value)"),
		)
	})

	t.Run("Expect(value).Header().NotOneOf()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect("Data").Header().NotOneOf()),
			PtrStr("only usable with Expect() not with Expect(value)"),
		)
	})

	t.Run("Expect(value).Header().Empty()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect("Data").Header().Empty()),
			PtrStr("only usable with Expect() not with Expect(value)"),
		)
	})

	t.Run("Expect(value).Header().Len()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect("Data").Header().Len(0)),
			PtrStr("only usable with Expect() not with Expect(value)"),
		)
	})

	t.Run("Expect(value).Header().Equal()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect("Data").Header().Equal(nil)),
			PtrStr("only usable with Expect() not with Expect(value)"),
		)
	})

	t.Run("Expect(value).Header().NotEqual()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect("Data").Header().NotEqual(nil)),
			PtrStr("only usable with Expect() not with Expect(value)"),
		)
	})

	t.Run("Expect(value).Status()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect("Data").Status()),
			PtrStr("only usable with Expect() not with Expect(value)"),
		)
	})

	t.Run("Expect(value).Status().Equal()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect("Data").Status().Equal(0)),
			PtrStr("only usable with Expect() not with Expect(value)"),
		)
	})

	t.Run("Expect(value).Custom()", func(t *testing.T) {
		ExpectError(t,
			Do(Expect("Data").Custom(nil)),
			PtrStr("only usable with Expect() not with Expect(value)"),
		)
	})
}

func TestExpect_WithoutArgument(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	ExpectError(t,
		Do(
			Post(s.URL),
			Expect(),
		),
		PtrStr("unable to run Expect() without an argument or without a chain. Please use Expect(something) or Expect().Something"),
	)
}
