package handlers

import (
	"bytes"
	"testing"
)

func TestSplitFile(t *testing.T) {
	// Prepare the test data
	data := []byte("This is a test file that we will split into parts.")
	file := bytes.NewReader(data)

	// Define the number of parts you want to split the file into
	numParts := 5

	// Call the splitFile function
	parts, err := splitFile(file, numParts)
	if err != nil {
		t.Fatalf("splitFile returned an error: %v", err)
	}

	// Check the number of parts
	if len(parts) != numParts {
		t.Errorf("Expected %d parts, got %d", numParts, len(parts))
	}

	// Check the size of each part (except for the last part)
	expectedSizes := []int{10, 10, 10, 10, 10} // Adjust this based on your expected part sizes
	for i, part := range parts {
		if len(part) != expectedSizes[i] {
			t.Errorf("Part %d has size %d, expected %d", i, len(part), expectedSizes[i])
		}
	}

	// Verify the content of the parts
	expectedParts := []string{
		"This is a ",
		"test file ",
		"that we wi",
		"ll split i",
		"nto parts.",
	}
	for i, part := range parts {
		if string(part) != expectedParts[i] {
			t.Errorf("Part %d content %q, expected %q", i, string(part), expectedParts[i])
		}
	}
}
