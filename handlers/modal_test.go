package handlers

import (
	"errors"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
)

func TestHandleModal_CreatesIssueAndResponds(t *testing.T) {
	reset := setCooldownStateForTest(2 * time.Second)
	defer reset()

	var gotTitle, gotBody string
	originalCreate := createIssue
	createIssue = func(title, body string) (string, error) {
		gotTitle = title
		gotBody = body
		return "https://example.com/issue/1", nil
	}
	defer func() { createIssue = originalCreate }()

	responder := &mockResponder{}
	interaction := modalInteraction("user-1", "Bug title", "Bug description")
	HandleModal(responder, interaction)

	if responder.called != 1 {
		t.Fatalf("expected responder to be called once, got %d", responder.called)
	}
	if gotTitle != "Bug title" || gotBody != "Bug description" {
		t.Fatalf("unexpected issue content: %q / %q", gotTitle, gotBody)
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
	reset := setCooldownStateForTest(30 * time.Second)
	defer reset()

	cooldownMu.Lock()
	lastReportByUser["user-1"] = time.Now()
	cooldownMu.Unlock()

	called := false
	originalCreate := createIssue
	createIssue = func(title, body string) (string, error) {
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
	reset := setCooldownStateForTest(2 * time.Second)
	defer reset()

	originalCreate := createIssue
	createIssue = func(title, body string) (string, error) {
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
			Member: &discordgo.Member{User: &discordgo.User{ID: userID}},
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

func setCooldownStateForTest(duration time.Duration) func() {
	cooldownMu.Lock()
	previousMap := lastReportByUser
	previousDuration := cooldownDuration
	lastReportByUser = map[string]time.Time{}
	cooldownDuration = duration
	cooldownMu.Unlock()

	return func() {
		cooldownMu.Lock()
		lastReportByUser = previousMap
		cooldownDuration = previousDuration
		cooldownMu.Unlock()
	}
}
