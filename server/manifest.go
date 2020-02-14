// This file is automatically generated. Do not modify it manually.

package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

var manifest *model.Manifest

const manifestStr = `
{
  "id": "com.mattermost.digitalocean",
  "name": "DigitalOcean Plugin",
  "description": "A DigitalOcean plugin for Mattermost",
  "version": "0.1.0",
  "min_server_version": "5.12.0",
  "server": {
    "executables": {
      "linux-amd64": "server/dist/plugin-linux-amd64",
      "darwin-amd64": "server/dist/plugin-darwin-amd64",
      "windows-amd64": "server/dist/plugin-windows-amd64.exe"
    },
    "executable": ""
  },
  "webapp": {
    "bundle_path": "webapp/dist/main.js"
  },
  "settings_schema": {
    "header": "Mattermost plugin for DigitalOcean Teams.",
    "footer": "",
    "settings": [
      {
        "key": "DOTeamID",
        "display_name": "Unique DigitalOcean Team Identifier",
        "type": "text",
        "help_text": "",
        "placeholder": "",
        "default": null
      },
      {
        "key": "DOAdmins",
        "display_name": "Users that are not system admins on Mattermost but have advanced plugin privileges",
        "type": "text",
        "help_text": "",
        "placeholder": "",
        "default": null
      }
    ]
  }
}
`

func init() {
	manifest = model.ManifestFromJson(strings.NewReader(manifestStr))
}
