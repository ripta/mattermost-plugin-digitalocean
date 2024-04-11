package main

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"
	"time"

	"github.com/digitalocean/godo"
	"github.com/mattermost/mattermost/server/public/model"

	"github.com/phillipahereza/mattermost-plugin-digitalocean/server/client"
)

func (p *Plugin) listDatabaseClustersCommandFunc(client client.DigitalOceanService, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {

	opts := &godo.ListOptions{}

	databases, response, err := client.ListDatabaseClusters(context.TODO(), opts)

	if err != nil {
		p.API.LogError("failed to fetch databases", "response", response, "Err", err.Error())
		return p.responsef(args, "Error while fetching databases list"),
			&model.AppError{Message: err.Error()}
	}

	if len(databases) == 0 {
		return p.responsef(args, "You don't have any databases configured"), nil
	}

	buffer := new(bytes.Buffer)

	w := new(tabwriter.Writer)

	w.Init(buffer, 8, 8, 0, '\t', 0)
	fmt.Fprintf(w, "\n |%s|%s|%s|%s|%s|%s|", "ID", "Name", "Engine", "Status", "Size", "Region")
	fmt.Fprintf(w, "\n |%s|%s|%s|%s|%s|%s|", "------", "----", "----", "----", "----", "----")

	for _, db := range databases {

		fmt.Fprintf(w, "\n |%s|%s|%s|%s|%s|%s|", db.ID, db.Name, db.EngineSlug, db.Status, db.SizeSlug, db.RegionSlug)
	}

	w.Flush()

	return p.responsef(args, buffer.String()), nil
}

func (p *Plugin) listDatabaseClusterBackupsCommandFunc(client client.DigitalOceanService, args *model.CommandArgs, id string) (*model.CommandResponse, *model.AppError) {

	backups, response, err := client.ListDatabaseClusterBackups(context.TODO(), id, nil)

	if err != nil {
		p.API.LogError("failed to get backups for database", "id", id, "response", response, "Err", err.Error())
		return p.responsef(args, "Error while fetching backups list for database cluster %s", id),
			&model.AppError{Message: err.Error()}
	}

	if len(backups) == 0 {
		return p.responsef(args, "You don't have any cluster backups"), nil
	}

	buffer := new(bytes.Buffer)

	w := new(tabwriter.Writer)

	w.Init(buffer, 8, 8, 0, '\t', 0)
	fmt.Fprintf(w, "\n |%s|%s|", "Created At", "Size (GB)")
	fmt.Fprintf(w, "\n |%s|%s|", "------", "----")

	for _, backup := range backups {

		fmt.Fprintf(w, "\n |%s|%f|", backup.CreatedAt.Format(time.RFC822), backup.SizeGigabytes)
	}

	w.Flush()

	return p.responsef(args, buffer.String()), nil
}

func (p *Plugin) addUserToDatabaseClusterCommandFunc(client client.DigitalOceanService, args *model.CommandArgs, id, name string) (*model.CommandResponse, *model.AppError) {
	dbUserCreateReq := &godo.DatabaseCreateUserRequest{
		Name:          name,
		MySQLSettings: nil,
	}
	user, response, err := client.CreateDatabaseUser(context.TODO(), id, dbUserCreateReq)
	if err != nil {
		p.API.LogError("failed to create user for database", "id", id, "response", response, "Err", err.Error())
		return p.responsef(args, "Error while creating a user on database %s because %s", id, err.Error()),
			&model.AppError{Message: err.Error()}
	}
	return p.responsef(args, "Name: `%s`\t Password: `%s`\t Role: `%s`", user.Name, user.Password, user.Role), nil
}

func (p *Plugin) listDatabaseClusterUsersCommandFunc(client client.DigitalOceanService, args *model.CommandArgs, id string) (*model.CommandResponse, *model.AppError) {

	users, response, err := client.ListDatabaseClusterUsers(context.TODO(), id, nil)

	if err != nil {
		p.API.LogError("failed to get users for database", "id", id, "response", response, "Err", err.Error())
		return p.responsef(args, "Error while fetching users for database cluster %s", id),
			&model.AppError{Message: err.Error()}
	}

	if len(users) == 0 {
		return p.responsef(args, "You don't have any cluster backups"), nil
	}

	usersList := ""
	for _, user := range users {
		usersList += fmt.Sprintf("- Name: `%s`\t Role: `%s`\n", user.Name, user.Role)
	}
	return p.responsef(args, usersList), nil
}

func (p *Plugin) deleteDatabaseClusterUserCommandFunc(client client.DigitalOceanService, args *model.CommandArgs, id, userName string) (*model.CommandResponse, *model.AppError) {

	response, err := client.DeleteDatabaseClusterUser(context.TODO(), id, userName)

	if err != nil {
		p.API.LogError("failed to delete user for database", "id", id, "name", userName, "response", response, "Err", err.Error())
		return p.responsef(args, "Error while deleting user %s for database cluster %s because %s", userName, id, err.Error()),
			&model.AppError{Message: err.Error()}
	}

	return p.responsef(args, "Successfully deleted user %s from database %s", userName, id), nil
}

func (p *Plugin) listDatabasesInClusterCommandFunc(client client.DigitalOceanService, args *model.CommandArgs, id string) (*model.CommandResponse, *model.AppError) {

	dbs, response, err := client.ListDatabasesInCluster(context.TODO(), id, nil)
	if err != nil {
		p.API.LogError("failed to list databases in this cluster", "id", id, "response", response, "Err", err.Error())
		return p.responsef(args, "Error listing databases in the cluster %s because %s", id, err.Error()),
			&model.AppError{Message: err.Error()}
	}

	if len(dbs) == 0 {
		return p.responsef(args, "You don't have any cluster backups"), nil
	}

	dbList := ""
	for i, db := range dbs {
		dbList += fmt.Sprintf("%d. Name: `%s`\n", i+1, db.Name)
	}
	return p.responsef(args, dbList), nil
}
