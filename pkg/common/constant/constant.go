package constant

const (
	HEART_BEAT = "heartbeat"
	PONG       = "pong"

	// message type, single chat or group chat
	MESSAGE_TYPE_USER  = 1
	MESSAGE_TYPE_GROUP = 2

	// message content type
	TEXT         = 1
	FILE         = 2
	IMAGE        = 3
	AUDIO        = 4
	VIDEO        = 5
	AUDIO_ONLINE = 6
	VIDEO_ONLINE = 7

	// messsage queue type
	GO_CHANNEL = "gochannel"
	KAFKA      = "kafka"
)