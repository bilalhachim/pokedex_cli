package main

import "testing"
func TestCleanInput(t *testing.T) {
    cases := []struct {
	input    string
	expected []string
}{
	{
		input:    "  hello  world  ",
		expected: []string{"hello", "world"},
	},
	// add more cases here
}
for _, c := range cases {
	actual := cleanInput(c.input)
	if(len(actual)!=len(c.expected)){
		t.Errorf("the slices don't have same size")
				t.Fail()
	}

	// Check the length of the actual slice against the expected slice
	// if they don't match, use t.Errorf to print an error message
	// and fail the test
	for i := range actual {
		word := actual[i]
		expectedWord := c.expected[i]
		for i := 0; i < len(actual); i++ {
			if(word!=expectedWord){
				t.Errorf("word don't match")
				t.Fail()
			}
		}
	}
}
}