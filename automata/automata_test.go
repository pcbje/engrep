package automata

import (
	"reflect"
	"sort"
	"testing"
)

func Test1(t *testing.T) {
	auto := CreateAutomata([]string{"bjelland", "petter"})
	expected := []string{"petter"}
	actual := []string{}

	auto.FindAll("peter", 1, func(hit string, distance int) {
		actual = append(actual, hit)
	})

	sort.Strings(expected)
	sort.Strings(actual)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}
