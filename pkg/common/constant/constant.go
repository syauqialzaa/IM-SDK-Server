package constant

const (
	HEART_BEAT		= "heartbeat"
	PONG       		= "pong"
	NULL_ID	int32	= 0
	
	// paths
	BASE_URL		= "http://localhost:8080"
	GET_ALL			= "/students"
	
	// methods
	METHOD_GET		= "GET"
	METHOD_POST		= "POST"
	METHOD_PUT		= "PUT"
	METHOD_PATCH	= "PATCH"
	METHOD_DELETE	= "DELETE"

	// message type, single chat or group chat
	MESSAGE_TYPE_USER  = 1
	MESSAGE_TYPE_GROUP = 2

	// uuid tags type
	USER_TAG_UUID  = "[usr]"
	GROUP_TAG_UUID = "[grp]"

	// delete status
	DELETE_FOR_ALL	= "all"
	DELETE_FOR_SELF	= "self"
	REINTERACT		= "reinteract"

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