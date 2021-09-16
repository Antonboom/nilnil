package analyzer

import "testing"

func TestCheckedTypes(t *testing.T) {
	c := newDefaultCheckedTypes()

	for _, tt := range knownTypes {
		assertTrue(t, c.Contains(tt))
	}
	assertStringEqual(t, c.String(), "chan,func,iface,map,ptr")

	err := c.Set("chan,iface,ptr")
	assertNoError(t, err)
	assertTrue(t, c.Contains(chanType))
	assertTrue(t, c.Contains(ifaceType))
	assertTrue(t, c.Contains(ptrType))
	assertStringEqual(t, "chan,iface,ptr", c.String())

	for _, tt := range knownTypes {
		err := c.Set(tt.S())
		assertNoError(t, err)

		for _, tt2 := range knownTypes {
			if tt2 == tt {
				assertTrue(t, c.Contains(tt2))
			} else {
				assertFalse(t, c.Contains(tt2))
			}
		}

		assertStringEqual(t, tt.S(), c.String())
	}
}

func TestCheckedTypes_SetError(t *testing.T) {
	c := newDefaultCheckedTypes()

	err := c.Set("unknown")
	assertError(t, err)
	t.Log(err)
}

func TestCheckedTypes_SetWithoutArg(t *testing.T) {
	c := newDefaultCheckedTypes()

	err := c.Set("")
	assertNoError(t, err)

	for _, tt := range knownTypes {
		assertTrue(t, c.Contains(tt))
	}
	assertStringEqual(t, c.String(), "chan,func,iface,map,ptr")
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.FailNow()
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.FailNow()
	}
}

func assertStringEqual(t *testing.T, a, b string) {
	t.Helper()
	if a != b {
		t.Logf("%q != %q", a, b)
		t.FailNow()
	}
}

func assertTrue(t *testing.T, v bool) {
	t.Helper()
	if !v {
		t.FailNow()
	}
}

func assertFalse(t *testing.T, v bool) {
	t.Helper()
	if v {
		t.FailNow()
	}
}
