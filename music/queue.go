package music

func (v *VCHandler) AddSongToQueue(conn *vcConnection, s Song) {
	conn.queue = append(conn.queue, s)
}
