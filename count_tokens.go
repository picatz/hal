package main

import "strings"

// countTokens counts the number of tokens in a string based.
//
// GPT-2 and GPT-3 use byte pair encoding to turn text into a series of integers to feed
// into the model. This just gives a rough estimate of the number of tokens in a string.
//
// https://github.com/openai/gpt-2
func countTokens(text string) int {
	// For now this works though.
	return len(strings.Fields(text))
}
