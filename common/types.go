package common

type Guild struct {
	ID     string
	Prefix string
}

func NewGuild(id string) Guild {
	return Guild{
		ID:     id,
		Prefix: ".",
	}
}

type Keyword struct {
	ID       int
	GuildID  string
	Key      string
	Reaction string
}

func NewKeyword(guildID string, key string, reaction string) Keyword {
	return Keyword{
		GuildID:  guildID,
		Key:      key,
		Reaction: reaction,
	}
}
