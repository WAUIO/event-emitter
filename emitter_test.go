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

func TestReturn1ArrayOf1Value(t *testing.T) {
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

func TestReturn2ArrayOf1Value(t *testing.T) {
	e := New()

	// (listener #1)
	e.On("event", func(a int, b int) int {
		return a + b
	})

	// (listener #2)
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

	// 3 + 4 = 7 (listener #1)
	if value1, isInt := values[0][0].(int); isInt {
		if value1 != 7 {
			t.Errorf("Got wrong value in values[0][0], got %d, expected 7", values[0][0])
		}
	} else {
		t.Errorf("got wrong value type, got %T, expected int", values[0][0])
	}

	// 3 x 4 = 12 (listener #2)
	if value2, isInt := values[1][0].(int); isInt {
		if value2 != 12 {
			t.Errorf("Got wrong value in values[0][0], got %d, expected 12", values[1][0])
		}
	} else {
		t.Errorf("got wrong value type, got %T, expected int", values[1][0])
	}
}
