package main

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/pkg/errors"
	cron "github.com/robfig/cron/v3"

	"github.com/phillipahereza/mattermost-plugin-digitalocean/server/client"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	store Store

	BotUserID string

	cron *cron.Cron
}

// OnActivate is
func (p *Plugin) OnActivate() error {
	p.API.RegisterCommand(getCommand())

	botID, err := p.API.EnsureBotUser(&model.Bot{
		Username:    "do",
		DisplayName: "DO",
		Description: "Created by the DigitalOcean plugin.",
	}) //, plugin.ProfileImagePath(profileImage))

	if err != nil {
		p.API.LogError("Failed to ensure digitalOcean bot", "Err", err.Error())
		return errors.Wrap(err, "failed to ensure digitalOcean bot")
	}
	p.BotUserID = botID

	if img, err := os.ReadFile(filepath.Join("assets", "do.png")); err == nil {
		_ = p.API.SetProfileImage(botID, img)
	}

	store := CreateStore(p)

	// Add an empty subscription where channels will be kept
	// Only add it if this was never done before
	sub, _ := store.LoadSubscription()
	if sub == nil {
		store.StoreSubscription(Subscription{})
	}

	p.store = store

	// Register cron
	p.cron = newCron()

	// start jobs
	p.StartCronJobs()

	return nil
}

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "do",
		DisplayName:      "do",
		Description:      "Integration with DigitalOcean.",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: help",
		AutoCompleteHint: "[command]",
	}
}

// GetClient returns a digital ocean client with configured token
func (p *Plugin) GetClient(mmUser string) (*client.DigitalOceanClient, error) {
	token, err := p.store.LoadUserDOToken(mmUser)
	if err != nil {
		p.API.LogError("Failed to load DO token", "user", mmUser, "Err", err.Error())
		return nil, err
	} else if token == "" {
		p.API.LogError("Failed to load DO token", "user", mmUser, "Err", err)
		return nil, errors.New("Missing DigitalOcean token: User `/do token` to get instructions on how to add a token")
	}

	tokenSource := &client.TokenSource{
		AccessToken: token,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	godoClient := godo.NewClient(oauthClient)
	return &client.DigitalOceanClient{Client: godoClient}, nil
}
