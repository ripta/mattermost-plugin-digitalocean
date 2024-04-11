package main

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/digitalocean/godo"
	"github.com/mattermost/mattermost/server/public/model"

	"github.com/phillipahereza/mattermost-plugin-digitalocean/server/client"
)

func drawKeysTable(keys []godo.Key) string {
	buffer := new(bytes.Buffer)

	w := new(tabwriter.Writer)

	w.Init(buffer, 8, 8, 0, '\t', 0)
	fmt.Fprintf(w, "\n |%s|%s|%s|", "ID", "Name", "Public Key")
	fmt.Fprintf(w, "\n |%s|%s|%s|", "------", "----", "------")

	for _, key := range keys {

		fmt.Fprintf(w, "\n |%d|%s|%s|", key.ID, key.Name, key.PublicKey)
	}

	w.Flush()
	return buffer.String()
}

func (p *Plugin) listSSHKeysCommandFunc(client client.DigitalOceanService, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {

	opts := &godo.ListOptions{}

	keys, response, err := client.ListSSHKeys(context.TODO(), opts)

	if err != nil {
		p.API.LogError("failed to fetch ssh keys", "response", response, "Err", err.Error())
		return p.responsef(args, "Error while fetching ssh keys list"),
			&model.AppError{Message: err.Error()}
	}

	if len(keys) == 0 {
		return p.responsef(args, "You don't have any ssh keys defined. Use `/do create-key <name> <publicKey>` to create one"), nil
	}

	return p.responsef(args, drawKeysTable(keys)), nil
}

func (p *Plugin) createSSHKeysCommandFunc(client client.DigitalOceanService, args *model.CommandArgs, name, publicKey string) (*model.CommandResponse, *model.AppError) {

	keyCreateRequest := &godo.KeyCreateRequest{
		Name:      name,
		PublicKey: publicKey,
	}

	key, response, err := client.CreateSSHKey(context.TODO(), keyCreateRequest)

	if err != nil {
		p.API.LogError("failed to create ssh key", "response", response, "Err", err.Error())
		return p.responsef(args, "Error while creating SSH key. %s", err.Error()),
			&model.AppError{Message: err.Error()}
	}

	return p.responsef(args, "Successfully created SSH key %s", key.Name), nil

}

func (p *Plugin) retrieveSSHKeyCommandFunc(client client.DigitalOceanService, args *model.CommandArgs, id int) (*model.CommandResponse, *model.AppError) {

	key, response, err := client.GetSSHKeyByID(context.TODO(), id)
	if err != nil {
		p.API.LogError("Failed to retrieve SSH key", "id", id, "response", response, "Err", err.Error())
		return p.responsef(args, "Failed to retrieve SSH key with ID `%d`", id),
			&model.AppError{Message: err.Error()}
	}
	return p.responsef(args, drawKeysTable([]godo.Key{*key})), nil

}

func (p *Plugin) deleteSSHKeyCommandFunc(client client.DigitalOceanService, args *model.CommandArgs, id int) (*model.CommandResponse, *model.AppError) {

	response, err := client.DeleteSSHKeyByID(context.TODO(), id)
	if err != nil {
		p.API.LogError("Failed to delete SSH key", "id", id, "response", response, "Err", err.Error())
		return p.responsef(args, "Failed to delete SSH key with ID `%d`", id),
			&model.AppError{Message: err.Error()}
	}

	return p.responsef(args, "Successfully delete SSH key with ID `%d`", id), nil
}
