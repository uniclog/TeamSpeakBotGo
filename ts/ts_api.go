package ts

import (
    "UnicBotGo/config"
    "fmt"
	"github.com/multiplay/go-ts3"
	"log"
	"strconv"
)

func InitNewClient(config *config.Config) *ts3.Client {
	client, err := ts3.NewClient(config.Address)
	handleError(err, "Failed to connect to TeamSpeak server")
	log.Println(fmt.Sprintf("Connected to %s", config.Address))
	return client
}

func CloseClient(client *ts3.Client) {
	if err := client.Close(); err != nil {
		log.Printf("Failed to close TeamSpeak client: %v", err)
	}
}

func Login(client *ts3.Client, config *config.Config) {
	handleError(client.Login(config.Username, config.Password), fmt.Sprintf("Login failed: '%s'", config.Username))
	log.Println(fmt.Sprintf("Login %s", config.Username))
}

func UseVirtualServer(client *ts3.Client) {
	handleError(client.Use(1), "Failed to select virtual server")
	log.Printf("Select virtual server: 1")
}

func SetNick(client *ts3.Client, config *config.Config) {
	connectionInfo, err := client.Whoami()
	handleError(err, "Failed to retrieve connection info")
	if connectionInfo.ClientName != config.BotName {
		handleError(client.SetNick(config.BotName), fmt.Sprintf("Failed to set nickname '%s'", config.BotName))
	}
	log.Println(fmt.Sprintf("Set nickname '%s'", config.BotName))
}

func RegisterServerEvents(client *ts3.Client) {
	registerEvent(client, ts3.ServerEvents)
}

func RegisterTextChannelEvents(client *ts3.Client) {
	registerEvent(client, ts3.TextChannelEvents)
}

func RegisterTextPrivateEvents(client *ts3.Client) {
	registerEvent(client, ts3.TextPrivateEvents)
}

func registerEvent(client *ts3.Client, event ts3.NotifyCategory) {
	handleError(client.Register(event), string("Failed to register: "+event))
	log.Printf("Register events: %s", event)
}

func GetActiveClientMap(client *ts3.Client) map[int]*ts3.OnlineClient {
	clients, err := client.Server.ClientList()
	handleError(err, "Failed to get client list")

	clientMap := make(map[int]*ts3.OnlineClient)
	for _, client := range clients {
		clientMap[client.ID] = client
	}
	return clientMap
}

func ClientMoveRequest(client *ts3.Client, id uint) {
	me, _ := client.Whoami()
	cmd := ts3.NewCmd("clientmove").WithArgs(
		ts3.NewArg("clid", strconv.Itoa(me.ClientID)),
		ts3.NewArg("cid", id),
	)
	_, err := client.ExecCmd(cmd)
	handleError(err, fmt.Sprintf("Failed to execute command: %v", err))
	log.Printf("%s moved to channel %d", me.ClientLoginName, id)
}

func SendMessageToClient(client *ts3.Client, userId int, message string) {
	cmd := ts3.NewCmd("sendtextmessage").WithArgs(
		ts3.NewArg("targetmode", 1),
		ts3.NewArg("target", strconv.Itoa(userId)),
		ts3.NewArg("msg", message),
	)
	_, _ = client.ExecCmd(cmd)
}

func SendMessageToChannel(client *ts3.Client, channelID int, message string) {
	cmd := ts3.NewCmd("sendtextmessage").WithArgs(
		ts3.NewArg("targetmode", 2),
		ts3.NewArg("target", strconv.Itoa(channelID)),
		ts3.NewArg("msg", message),
	)
	_, _ = client.ExecCmd(cmd)
}

func ChangeChannelName(client *ts3.Client, channelID int, tittle string) {
	cmd := ts3.NewCmd("channeledit").WithArgs(
		ts3.NewArg("cid", strconv.Itoa(channelID)),
		ts3.NewArg("channel_name", tittle),
	)
	_, err := client.ExecCmd(cmd)
	if err != nil {
		log.Fatalf("editchannel err: %v", err)
	}
}

func handleError(err error, message ...string) {
	if err != nil {
		log.Fatalf("%s: %v", message[0], err)
	}
	if len(message) > 1 {
		log.Printf(message[1])
	}
}
