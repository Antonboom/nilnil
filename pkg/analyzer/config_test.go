package analyzer

import "testing"

func TestCheckedTypes(t *testing.T) {
	c := newDefaultCheckedTypes()
	assertStringEqual(t, c.String(), "chan,func,iface,map,ptr,uintptr,unsafeptr")

	err := c.Set("chan,iface,ptr")
	assertNoError(t, err)
	assertTrue(t, c.Contains(chanType))
	assertTrue(t, c.Contains(ifaceType))
	assertTrue(t, c.Contains(ptrType))
	assertFalse(t, c.Contains(funcType))
	assertFalse(t, c.Contains(mapType))
	assertFalse(t, c.Contains(uintptrType))
	assertFalse(t, c.Contains(unsafeptrType))
	assertStringEqual(t, "chan,iface,ptr", c.String())
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
	assertStringEqual(t, c.String(), "chan,func,iface,map,ptr,uintptr,unsafeptr")
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("must be not nil")
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("must be nil, got %q", err)
	}
}

func assertStringEqual(t *testing.T, a, b string) {
	t.Helper()
	if a != b {
		t.Fatalf("%q != %q", a, b)
	}
}

func assertTrue(t *testing.T, v bool) {
	t.Helper()
	if !v {
		t.Fatal("must be true")
	}
}

func assertFalse(t *testing.T, v bool) {
	t.Helper()
	if v {
		t.Fatal("must be false")
	}
}
