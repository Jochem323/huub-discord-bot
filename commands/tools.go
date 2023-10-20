package commands

import "github.com/bwmarrin/discordgo"

func IsAdmin(d *discordgo.Session, m *discordgo.MessageCreate) bool {
	perms, err := d.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		return false
	}

	return perms&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator
}
