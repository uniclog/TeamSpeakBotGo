package ts

import (
	"github.com/multiplay/go-ts3"
	"strconv"
	"strings"
)

func GetClientInfo(data map[string]string) ts3.OnlineClient {
	client := ts3.OnlineClient{}

	client.ID = parseInt(data["clid"])
	client.ChannelID = parseInt(data["cid"])
	client.DatabaseID = parseInt(data["client_database_id"])
	client.Nickname = data["client_nickname"]
	client.Type = parseInt(data["client_type"])
	client.Away = parseBool(data["client_away"])
	client.AwayMessage = data["client_away_message"]

	client.OnlineClientExt = &ts3.OnlineClientExt{
		UniqueIdentifier: stringPointer(data["client_unique_identifier"]),
		Country:          stringPointer(data["client_country"]),
		IP:               stringPointer(data["connection_client_ip"]),
		Badges:           stringPointer(data["client_badges"]),
		IconID:           parseIntPointer(data["client_icon_id"]),
	}
	client.OnlineClientExt.OnlineClientGroups = &ts3.OnlineClientGroups{
		ChannelGroupID:                 parseIntPointer(data["client_channel_group_id"]),
		ChannelGroupInheritedChannelID: parseIntPointer(data["client_channel_group_inherited_channel_id"]),
		ServerGroups:                   parseIntSlicePointer(data["client_servergroups"]),
	}
	client.OnlineClientExt.OnlineClientVoice = &ts3.OnlineClientVoice{
		FlagTalking:        parseBoolPointer(data["client_flag_talking"]),
		IsChannelCommander: parseBoolPointer(data["client_is_channel_commander"]),
	}
	client.OnlineClientExt.OnlineClientTimes = &ts3.OnlineClientTimes{
		IdleTime:      parseIntPointer(data["client_idle_time"]),
		Created:       parseIntPointer(data["client_created"]),
		LastConnected: parseIntPointer(data["client_lastconnected"]),
	}
	client.OnlineClientExt.OnlineClientInfo = &ts3.OnlineClientInfo{
		Version:  stringPointer(data["client_version"]),
		Platform: stringPointer(data["client_platform"]),
	}

	return client
}

type MessageStruct struct {
	InvokerId        int    `json:"invokerid"`
	InvokerName      string `json:"invokername"`
	UniqueIdentifier string `json:"invokeruid"`
	Message          string `json:"msg"`
	Target           int    `json:"target"`
	TargetMode       int    `json:"targetmode"`
}

func GetMessageInfo(data map[string]string) MessageStruct {
	message := MessageStruct{}
	message.InvokerId = parseInt(data["invokerid"])
	message.InvokerName = data["invokername"]
	message.UniqueIdentifier = data["invokeruid"]
	message.Message = data["msg"]
	message.Target = parseInt(data["target"])
	message.TargetMode = parseInt(data["targetmode"])
	return message
}

func parseInt(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return value
}

func parseIntPointer(s string) *int {
	value, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &value
}

func stringPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func parseBool(s string) bool {
	value, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return value
}

func parseBoolPointer(s string) *bool {
	value, err := strconv.ParseBool(s)
	if err != nil {
		return nil
	}
	return &value
}

func parseIntSlicePointer(s string) *[]int {
	if s == "" {
		return nil
	}
	stringSlice := strings.Split(s, ",")
	intSlice := make([]int, len(stringSlice))
	for i, v := range stringSlice {
		intSlice[i], _ = strconv.Atoi(v)
	}
	return &intSlice
}
