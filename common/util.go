package common

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func CheckAdmin(session *discordgo.Session, channelID string, userID string) bool {
	return checkPermission(session, discordgo.PermissionAdministrator, channelID, userID)
}

func checkPermission(session *discordgo.Session, permission int, channelID string, userID string) bool {
	permissions, err := session.State.UserChannelPermissions(userID, channelID)
	if err != nil {
		fmt.Printf("%s\n", err)
		return false
	}
	return permissions&permission == permission
}
