package commands

import "github.com/bwmarrin/discordgo"

func IsAdmin(d *discordgo.Session, m *discordgo.MessageCreate) bool {
	perms, err := d.State.MessagePermissions(m.Message)
	if err != nil {
		return false
	}

	return perms&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator
}
