package event_emitter

import "testing"

func TestDifferentEventHandlers(t *testing.T) {
	one := New()
	two := New()

	one.On("one", func() {})

	if !one.Has("one") {
		t.Errorf("one emitter excpect to have \"one\" event")
	}

	if two.Has("one") {
		t.Errorf("two emitter should not have \"one\" event")
	}
}

func TestReturn1Value(t *testing.T) {
	e := New()

	e.On("event", func(a int, b int) int {
		return a + b
	})

	values, err := e.Fire("event", 3, 4)

	if err != nil {
		t.Errorf("Got error %s", err)
	}

	if len(values) != 1 {
		t.Errorf("Number of returned value is wrong")
	}
}

func TestReturn2Values(t *testing.T) {
	e := New()

	e.On("event", func(a int, b int) int {
		return a + b
	})

	e.On("event", func(a int, b int) int {
		return a * b
	})

	values, err := e.Fire("event", 3, 4)

	if err != nil {
		t.Errorf("Got error %s", err)
	}

	if len(values) != 2 {
		t.Errorf("Number of returned value is wrong")
	}
}
