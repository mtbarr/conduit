package handlers

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

type mockResponder struct {
	called      int
	interaction *discordgo.Interaction
	response    *discordgo.InteractionResponse
	err         error
}

func (m *mockResponder) InteractionRespond(i *discordgo.Interaction, resp *discordgo.InteractionResponse, _ ...discordgo.RequestOption) error {
	m.called++
	m.interaction = i
	m.response = resp
	return m.err
}

func TestHandleCommand_ShowsModal(t *testing.T) {
	responder := &mockResponder{}
	interaction := &discordgo.Interaction{
		Type: discordgo.InteractionApplicationCommand,
		Data: discordgo.ApplicationCommandInteractionData{
			Name: "reportbug",
		},
	}
	HandleCommand(responder, &discordgo.InteractionCreate{Interaction: interaction})

	if responder.called != 1 {
		t.Fatalf("expected responder to be called once, got %d", responder.called)
	}
	if responder.response == nil {
		t.Fatal("expected a response")
	}
	if responder.response.Type != discordgo.InteractionResponseModal {
		t.Fatalf("expected modal response, got %d", responder.response.Type)
	}
	if responder.response.Data == nil || responder.response.Data.CustomID != "modal_reportbug" {
		t.Fatalf("expected modal_reportbug, got %#v", responder.response.Data)
	}
}

func TestHandleCommand_IgnoresOtherCommands(t *testing.T) {
	responder := &mockResponder{}
	interaction := &discordgo.Interaction{
		Type: discordgo.InteractionApplicationCommand,
		Data: discordgo.ApplicationCommandInteractionData{
			Name: "other",
		},
	}
	HandleCommand(responder, &discordgo.InteractionCreate{Interaction: interaction})

	if responder.called != 0 {
		t.Fatalf("expected responder not called, got %d", responder.called)
	}
}
