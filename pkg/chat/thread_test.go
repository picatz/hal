package chat

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/picatz/openai"
)

func TestThreadSummarize(t *testing.T) {
	thread := &Thread{
		Name: "Test Thread",
		ChatHistory: []openai.ChatMessage{
			{
				Role:    openai.ChatRoleUser,
				Content: "Who is Jon Snow's father? ",
			},
			{
				Role: openai.ChatRoleAssistant,
				Content: "It is revealed in the show that Jon Snow's father is Rhaegar Targaryen, " +
					"making him a true Targaryen heir. However, in the books, it remains a popular " +
					"theory that his father is also Rhaegar, making him the legitimate heir to the Iron Throne.",
			},
			{
				Role:    openai.ChatRoleUser,
				Content: "What is his mother?",
			},
			{
				Role: openai.ChatRoleAssistant,
				Content: "In the TV show, Jon Snow's mother is revealed to be Lyanna Stark. " +
					"She is the younger sister of Ned Stark, who is Jon Snow's adoptive father. " +
					"In the books, it is strongly suggested that the same is true, but it has not yet been explicitly confirmed.",
			},
		},
	}

	summary, err := thread.Summarize(context.Background(), openai.NewClient(os.Getenv("OPENAI_API_KEY")))
	if err != nil {
		t.Fatal(err)
	}

	// Must contain the following words
	words := []string{
		"Jon Snow",
		"Rhaegar Targaryen",
		"Lyanna Stark",
	}

	for _, word := range words {
		if !strings.Contains(summary, word) {
			t.Logf("summary does not contain %s", word)
			t.Fail()
		}
	}

	t.Log(summary)
}

func TestThreadSearch(t *testing.T) {
	thread := &Thread{
		Name: "Test Thread",
		ChatHistory: []openai.ChatMessage{
			{
				Role:    openai.ChatRoleUser,
				Content: "Who is Jon Snow's father? ",
			},
			{
				Role: openai.ChatRoleAssistant,
				Content: "It is revealed in the show that Jon Snow's father is Rhaegar Targaryen, " +
					"making him a true Targaryen heir. However, in the books, it remains a popular " +
					"theory that his father is also Rhaegar, making him the legitimate heir to the Iron Throne.",
			},
			{
				Role:    openai.ChatRoleUser,
				Content: "What is his mother?",
			},
			{
				Role: openai.ChatRoleAssistant,
				Content: "In the TV show, Jon Snow's mother is revealed to be Lyanna Stark. " +
					"She is the younger sister of Ned Stark, who is Jon Snow's adoptive father. " +
					"In the books, it is strongly suggested that the same is true, but it has not yet been explicitly confirmed.",
			},
		},
	}

	matches, err := thread.Search(context.Background(), "Jon Snow")
	if err != nil {
		t.Fatal(err)
	}

	if len(matches) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(matches))
	}

	for _, match := range matches {
		t.Log(match.Message.Content)
	}
}
