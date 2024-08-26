package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoteService_ValidateText(t *testing.T) {
	note := newNoteService(nil)

	testCases := []struct {
		testName   string
		input      string
		expectWord string
	}{
		{
			testName:   "Correct test",
			input:      "Превет, как дела?",
			expectWord: "Превет",
		},
	}

	for _, tc := range testCases {
		textErrs, err := note.ValidateText(tc.input)

		assert.Equal(t, nil, err)
		assert.Len(t, textErrs, 1)
		assert.Equal(t, tc.expectWord, textErrs[0].Word)
	}
}
