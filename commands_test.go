package main

import (
	"conduit/i18n"
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
)

type mockRegistrar struct {
	createAppID string
	createGuild string
	createCmd   *discordgo.ApplicationCommand
	createErr   error
	deleteAppID string
	deleteGuild string
	deleteCmdID string
	deleteErr   error
}

func (m *mockRegistrar) ApplicationCommandCreate(appID, guildID string, cmd *discordgo.ApplicationCommand, _ ...discordgo.RequestOption) (*discordgo.ApplicationCommand, error) {
	m.createAppID = appID
	m.createGuild = guildID
	m.createCmd = cmd
	if m.createErr != nil {
		return nil, m.createErr
	}
	return &discordgo.ApplicationCommand{ID: "cmd-1"}, nil
}

func (m *mockRegistrar) ApplicationCommandDelete(appID, guildID, cmdID string, _ ...discordgo.RequestOption) error {
	m.deleteAppID = appID
	m.deleteGuild = guildID
	m.deleteCmdID = cmdID
	return m.deleteErr
}

func TestRegisterReportBugCommand(t *testing.T) {
	registrar := &mockRegistrar{}
	cmd, err := registerReportBugCommand(registrar, "app-1", "guild-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd == nil || cmd.ID != "cmd-1" {
		t.Fatalf("unexpected command: %#v", cmd)
	}
	if registrar.createAppID != "app-1" || registrar.createGuild != "guild-1" {
		t.Fatalf("unexpected create args: %s %s", registrar.createAppID, registrar.createGuild)
	}
	expectedName := i18n.T("reportbug_command_name")
	if registrar.createCmd == nil || registrar.createCmd.Name != expectedName {
		t.Fatalf("unexpected command details: expected name %s, got %s", expectedName, registrar.createCmd.Name)
	}
}

func TestDeleteReportBugCommand(t *testing.T) {
	registrar := &mockRegistrar{}
	if err := deleteReportBugCommand(registrar, "app-1", "guild-1", "cmd-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if registrar.deleteAppID != "app-1" || registrar.deleteGuild != "guild-1" || registrar.deleteCmdID != "cmd-1" {
		t.Fatalf("unexpected delete args: %s %s %s", registrar.deleteAppID, registrar.deleteGuild, registrar.deleteCmdID)
	}
}

func TestRegisterReportBugCommand_Error(t *testing.T) {
	registrar := &mockRegistrar{createErr: errors.New("boom")}
	if _, err := registerReportBugCommand(registrar, "app-1", "guild-1"); err == nil {
		t.Fatal("expected error")
	}
}
