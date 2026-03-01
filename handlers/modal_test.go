package handlers

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
)

func TestHandleModal_CreatesIssueAndResponds(t *testing.T) {
	cooldownMu.Lock()
	lastActionByKey = map[string]time.Time{}
	cooldownMu.Unlock()

	var gotTitle, gotBody string
	var gotLabels []string
	originalCreate := createIssue
	createIssue = func(title, body string, labels []string) (string, error) {
		gotTitle = title
		gotBody = body
		gotLabels = labels
		return "https://example.com/issue/1", nil
	}
	defer func() { createIssue = originalCreate }()

	responder := &mockResponder{}
	interaction := modalInteraction("user-1", "Bug title", "Bug description")
	HandleModal(responder, interaction)

	if responder.called != 1 {
		t.Fatalf("expected responder to be called once, got %d", responder.called)
	}
	if gotTitle != "Bug title" {
		t.Fatalf("unexpected issue title: %q", gotTitle)
	}
	if !strings.Contains(gotBody, "## User Submission") || !strings.Contains(gotBody, "### Description") || !strings.Contains(gotBody, "Bug description") || !strings.Contains(gotBody, "### Metadata") || !strings.Contains(gotBody, "Reporter:") || !strings.Contains(gotBody, "<@user-1>") || !strings.Contains(gotBody, "Submitted at:") || !strings.Contains(gotBody, "```") {
		t.Fatalf("unexpected issue body: %q", gotBody)
	}
	if len(gotLabels) != 1 || gotLabels[0] != "bug" {
		t.Fatalf("unexpected labels: %v", gotLabels)
	}
	if responder.response == nil || responder.response.Data == nil {
		t.Fatal("expected response data")
	}
	if responder.response.Data.Flags != discordgo.MessageFlagsEphemeral {
		t.Fatalf("expected ephemeral response, got %d", responder.response.Data.Flags)
	}
	if responder.response.Data.Content == "" {
		t.Fatal("expected response content")
	}
}

func TestHandleModal_CooldownBlocks(t *testing.T) {
	cooldownMu.Lock()
	lastActionByKey = map[string]time.Time{}
	lastActionByKey["reportbug_user-1"] = time.Now()
	cooldownMu.Unlock()

	called := false
	originalCreate := createIssue
	createIssue = func(title, body string, labels []string) (string, error) {
		called = true
		return "", nil
	}
	defer func() { createIssue = originalCreate }()

	responder := &mockResponder{}
	interaction := modalInteraction("user-1", "Bug title", "Bug description")
	HandleModal(responder, interaction)

	if called {
		t.Fatal("expected createIssue not to be called during cooldown")
	}
	if responder.called != 1 {
		t.Fatalf("expected responder to be called once, got %d", responder.called)
	}
	if responder.response == nil || responder.response.Data == nil {
		t.Fatal("expected response data")
	}
}

func TestHandleModal_CreateIssueError(t *testing.T) {
	cooldownMu.Lock()
	lastActionByKey = map[string]time.Time{}
	cooldownMu.Unlock()

	originalCreate := createIssue
	createIssue = func(title, body string, labels []string) (string, error) {
		return "", errors.New("boom")
	}
	defer func() { createIssue = originalCreate }()

	responder := &mockResponder{}
	interaction := modalInteraction("user-1", "Bug title", "Bug description")
	HandleModal(responder, interaction)

	if responder.called != 1 {
		t.Fatalf("expected responder to be called once, got %d", responder.called)
	}
	if responder.response == nil || responder.response.Data == nil {
		t.Fatal("expected response data")
	}
	if responder.response.Data.Content == "" {
		t.Fatal("expected response content")
	}
}

func modalInteraction(userID, title, description string) *discordgo.InteractionCreate {
	return &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type:   discordgo.InteractionModalSubmit,
			Member: &discordgo.Member{User: &discordgo.User{ID: userID, Username: "testuser", Discriminator: "1234"}},
			Data: discordgo.ModalSubmitInteractionData{
				CustomID: "modal_reportbug",
				Components: []discordgo.MessageComponent{
					&discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						&discordgo.TextInput{
							CustomID: "bug_title",
							Value:    title,
						},
					}},
					&discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						&discordgo.TextInput{
							CustomID: "bug_description",
							Value:    description,
						},
					}},
				},
			},
		},
	}
}
