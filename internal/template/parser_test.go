package template

import (
	"testing"
)

func TestTemplateFillOut(t *testing.T) {
	t.Run("Test PullRequestTemplate FillOut", func(t *testing.T) {
		pullRequestTemplate, err := NewPullRequestTemplate()
		if err != nil {
			t.Fatalf("failed to create PullRequestTemplate: %v", err)
		}

		filled, err := pullRequestTemplate.FillOut("memberName", "postTitle", "publishedAt", "postUrl", "summary")

		if err != nil {
			t.Fatalf("failed to execute PullRequestTemplate: %v", err)
		}

		t.Log(filled)
	})

	t.Run("Test SummaryPromptTemplate FillOut", func(t *testing.T) {
		summaryPromptTemplate, err := NewSummaryPromptTemplate()

		if err != nil {
			t.Fatalf("failed to create SummaryPromptTemplate: %v", err)
		}

		filled, err := summaryPromptTemplate.FillOut("타이틀!", "컨텐츠!")

		if err != nil {
			t.Fatalf("failed to fill filled SummaryPromptTemplate: %v", err)
		}

		t.Log(filled)
	})
}
