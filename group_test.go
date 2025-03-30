/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * Tests fop groups of attributes
 */

package goipp

import "testing"

// TestGroupEqualSimilar tests Group.Equal and Group.Similar
func TestGroupEqualSimilar(t *testing.T) {
	type testData struct {
		g1, g2  Group // A pair of Attributes slice
		equal   bool  // Expected g1.Equal(g2) output
		similar bool  // Expected g2.Similar(g2) output
	}

	attrs1 := Attributes{
		MakeAttr("attr1", TagInteger, Integer(1)),
		MakeAttr("attr2", TagInteger, Integer(2)),
		MakeAttr("attr3", TagInteger, Integer(3)),
	}

	attrs2 := Attributes{
		MakeAttr("attr3", TagInteger, Integer(3)),
		MakeAttr("attr2", TagInteger, Integer(2)),
		MakeAttr("attr1", TagInteger, Integer(1)),
	}

	tests := []testData{
		{
			g1:      Group{TagJobGroup, nil},
			g2:      Group{TagJobGroup, nil},
			equal:   true,
			similar: true,
		},

		{
			g1:      Group{TagJobGroup, Attributes{}},
			g2:      Group{TagJobGroup, Attributes{}},
			equal:   true,
			similar: true,
		},

		{
			g1:      Group{TagJobGroup, Attributes{}},
			g2:      Group{TagJobGroup, nil},
			equal:   false,
			similar: true,
		},

		{
			g1:      Group{TagJobGroup, attrs1},
			g2:      Group{TagJobGroup, attrs1},
			equal:   true,
			similar: true,
		},

		{
			g1:      Group{TagJobGroup, attrs1},
			g2:      Group{TagJobGroup, attrs2},
			equal:   false,
			similar: true,
		},
	}

	for _, test := range tests {
		equal := test.g1.Equal(test.g2)
		similar := test.g1.Similar(test.g2)

		if equal != test.equal {
			t.Errorf("testing Group.Equal:\n"+
				"attrs 1:   %s\n"+
				"attrs 2:   %s\n"+
				"expected:  %v\n"+
				"present:   %v\n",
				test.g1, test.g2,
				test.equal, equal,
			)
		}

		if similar != test.similar {
			t.Errorf("testing Group.Similar:\n"+
				"attrs 1:  %s\n"+
				"attrs 2:  %s\n"+
				"expected: %v\n"+
				"present:  %v\n",
				test.g1, test.g2,
				test.similar, similar,
			)
		}
	}
}

// TestGroupAdd tests Group.Add
func TestGroupAdd(t *testing.T) {
	g1 := Group{
		TagJobGroup,
		Attributes{
			MakeAttr("attr1", TagInteger, Integer(1)),
			MakeAttr("attr2", TagInteger, Integer(2)),
			MakeAttr("attr3", TagInteger, Integer(3)),
		},
	}

	g2 := Group{Tag: TagJobGroup}
	for _, attr := range g1.Attrs {
		g2.Add(attr)
	}

	if !g1.Equal(g2) {
		t.Errorf("Group.Add test failed:\n"+
			"expected: %#v\n"+
			"present:  %#v\n",
			g1, g2,
		)
	}
}

// TestGroupCopy tests Group.Clone and Group.DeepCopy
func TestGroupCopy(t *testing.T) {
	type testData struct {
		g Group
	}

	attrs := Attributes{
		MakeAttr("attr1", TagInteger, Integer(1)),
		MakeAttr("attr2", TagInteger, Integer(2)),
		MakeAttr("attr3", TagInteger, Integer(3)),
	}

	tests := []testData{
		{Group{TagJobGroup, nil}},
		{Group{TagJobGroup, Attributes{}}},
		{Group{TagJobGroup, attrs}},
	}

	for _, test := range tests {
		clone := test.g.Clone()

		if !test.g.Equal(clone) {
			t.Errorf("testing Group.Clone\n"+
				"expected: %#v\n"+
				"present:  %#v\n",
				test.g,
				clone,
			)
		}

		copy := test.g.DeepCopy()
		if !test.g.Equal(copy) {
			t.Errorf("testing Group.DeepCopy\n"+
				"expected: %#v\n"+
				"present:  %#v\n",
				test.g,
				copy,
			)
		}
	}
}

// TestGroupEqualSimilar tests Group.Equal and Group.Similar
func TestGroupsEqualSimilar(t *testing.T) {
	type testData struct {
		groups1, groups2 Groups // A pair of Attributes slice
		equal            bool   // Expected g1.Equal(g2) output
		similar          bool   // Expected g2.Similar(g2) output
	}

	g1 := Group{
		TagJobGroup,
		Attributes{MakeAttr("attr1", TagInteger, Integer(1))},
	}

	g2 := Group{
		TagJobGroup,
		Attributes{MakeAttr("attr2", TagInteger, Integer(2))},
	}

	g3 := Group{
		TagPrinterGroup,
		Attributes{MakeAttr("attr2", TagInteger, Integer(2))},
	}

	tests := []testData{
		{
			// nil equal and similar to nil
			groups1: nil,
			groups2: nil,
			equal:   true,
			similar: true,
		},

		{
			// Empty groups equal and similar to empty groups
			groups1: Groups{},
			groups2: Groups{},
			equal:   true,
			similar: true,
		},

		{
			// nil similar but not equal to empty groups
			groups1: nil,
			groups2: Groups{},
			equal:   false,
			similar: true,
		},

		{
			// groups of different size neither equal nor similar
			groups1: Groups{g1, g2, g3},
			groups2: Groups{g1, g2},
			equal:   false,
			similar: false,
		},

		{
			// Same list of groups: equal and similar
			groups1: Groups{g1, g2, g3},
			groups2: Groups{g1, g2, g3},
			equal:   true,
			similar: true,
		},

		{
			// Groups with different group tags reordered.
			// Similar but not equal.
			groups1: Groups{g1, g2, g3},
			groups2: Groups{g3, g1, g2},
			equal:   false,
			similar: true,
		},

		{
			// Groups with the same group tags reordered.
			// Neither equal nor similar
			groups1: Groups{g1, g2, g3},
			groups2: Groups{g2, g1, g3},
			equal:   false,
			similar: false,
		},
	}

	for _, test := range tests {
		equal := test.groups1.Equal(test.groups2)
		similar := test.groups1.Similar(test.groups2)

		if equal != test.equal {
			t.Errorf("testing Groups.Equal:\n"+
				"attrs 1:   %s\n"+
				"attrs 2:   %s\n"+
				"expected:  %v\n"+
				"present:   %v\n",
				test.groups1, test.groups2,
				test.equal, equal,
			)
		}

		if similar != test.similar {
			t.Errorf("testing Groups.Similar:\n"+
				"attrs 1:  %s\n"+
				"attrs 2:  %s\n"+
				"expected: %v\n"+
				"present:  %v\n",
				test.groups1, test.groups2,
				test.similar, similar,
			)
		}
	}
}

// TestGroupsAdd tests Groups.Add
func TestGroupsAdd(t *testing.T) {
	g1 := Group{
		TagJobGroup,
		Attributes{MakeAttr("attr1", TagInteger, Integer(1))},
	}

	g2 := Group{
		TagJobGroup,
		Attributes{MakeAttr("attr2", TagInteger, Integer(2))},
	}

	g3 := Group{
		TagPrinterGroup,
		Attributes{MakeAttr("attr2", TagInteger, Integer(2))},
	}

	groups1 := Groups{g1, g2, g3}

	groups2 := Groups{}
	groups2.Add(g1)
	groups2.Add(g2)
	groups2.Add(g3)

	if !groups1.Equal(groups2) {
		t.Errorf("Groups.Add test failed:\n"+
			"expected: %#v\n"+
			"present:  %#v\n",
			groups1, groups2,
		)
	}
}

// TestGroupsCopy tests Groups.Clone and Groups.DeepCopy
func TestGroupsCopy(t *testing.T) {
	g1 := Group{
		TagJobGroup,
		Attributes{MakeAttr("attr1", TagInteger, Integer(1))},
	}

	g2 := Group{
		TagJobGroup,
		Attributes{MakeAttr("attr2", TagInteger, Integer(2))},
	}

	g3 := Group{
		TagPrinterGroup,
		Attributes{MakeAttr("attr2", TagInteger, Integer(2))},
	}

	type testData struct {
		groups Groups
	}

	tests := []testData{
		{nil},
		{Groups{}},
		{Groups{g1, g2, g3}},
	}

	for _, test := range tests {
		clone := test.groups.Clone()

		if !test.groups.Equal(clone) {
			t.Errorf("testing Groups.Clone\n"+
				"expected: %#v\n"+
				"present:  %#v\n",
				test.groups,
				clone,
			)
		}

		copy := test.groups.DeepCopy()
		if !test.groups.Equal(copy) {
			t.Errorf("testing Groups.DeepCopy\n"+
				"expected: %#v\n"+
				"present:  %#v\n",
				test.groups,
				copy,
			)
		}
	}
}
