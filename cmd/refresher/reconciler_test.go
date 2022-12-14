package main

import (
	"sort"
	"testing"
)

func TestComparer(t *testing.T) {
	la := []string{"a", "b"}
	lb := []string{"b", "a"}
	lc := []string{"a", "d"}
	ld := []string{"a", "b", "c"}
	if !lstcmp(la, la) || !lstcmp(la, lb) {
		t.Fail()
	}
	if lstcmp(la, lc) || lstcmp(la, ld) {
		t.Fail()
	}
}

func TestReconciler(t *testing.T) {
	plan := []string{"a", "b", "c", "d"}
	fact := []string{"c", "f", "a"}
	c := []string{"b", "d"}
	u := []string{"a", "c"}
	d := []string{"f"}
	c1, u1, d1 := reconcileLists(plan, fact)
	if !lstcmp(c1, c) {
		t.Errorf("Expected %s, got %s", c, c1)
	}
	if !lstcmp(u1, u) {
		t.Errorf("Expected %s, got %s", u, u1)
	}
	if !lstcmp(d1, d) {
		t.Errorf("Expected %s, got %s", d, d1)
	}

}

func lstcmp(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)
	for i, ia := range a {
		if b[i] != ia {
			return false
		}
	}
	return true
}
