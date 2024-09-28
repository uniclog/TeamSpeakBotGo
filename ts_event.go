package main

import (
	"fmt"
	"github.com/multiplay/go-ts3"
	"log"
	"regexp"
	"strings"
	"time"
)

const (
	Join           string = "cliententerview"
	Left           string = "clientleftview"
	ChannelMessage string = "textmessage"
)
const (
	VCNext string = "!next"
	VCPrev string = "!prev"
	VCGoto string = "!goto"
	VCInfo string = "!info"
)

const (
	Welcome string = "Привет! Хорошего настроения! И приятно провести время!"
)

var onlineClients map[int]*ts3.OnlineClient
var me *ts3.ConnectionInfo

func serverEventsListener(client *ts3.Client) {
	// get active clients
	onlineClients = GetActiveClientMap(client)
	me, _ = client.Whoami()

	printOnlineClients(onlineClients)

	// process events
	log.Println("Starting event listener")
	for notification := range client.Notifications() {
		//log.Println(notification)
		switch notification.Type {
		case Join:
			{
				clientInfo := GetClientInfo(notification.Data)
				log.Printf("Client joined: UniqueIdentifier=%s ID=%d, Name=%s", *clientInfo.UniqueIdentifier,
					clientInfo.ID, clientInfo.Nickname)
				// add to active users list
				onlineClients[clientInfo.ID] = &clientInfo
				//printOnlineClients(onlineClients)
				// send welcome message
				SendMessageToClient(client, clientInfo.ID, Welcome)
			}
		case Left:
			{
				clientInfo := GetClientInfo(notification.Data)
				log.Printf("Client left: ID=%d", clientInfo.ID)

				// delete from active users list
				delete(onlineClients, clientInfo.ID)
				//printOnlineClients(onlineClients)
			}
		case ChannelMessage:
			{
				//log.Println(notification)
				msg := GetMessageInfo(notification.Data)
				if msg.InvokerId == me.ClientID {
					break
				}
				switch msg.TargetMode {
				case 1: // Private
					{
						log.Printf("Private --> %d:%s - %s", msg.InvokerId, msg.InvokerName, msg.Message)
					}
				case 2: // Channel
					{
						log.Printf("Channel --> %d:%s - %s", msg.InvokerId, msg.InvokerName, msg.Message)
						switch {
						case strings.HasPrefix(msg.Message, VCNext):
							NextStation()
						case strings.HasPrefix(msg.Message, VCPrev):
							PrevStation()
						case strings.HasPrefix(msg.Message, VCInfo):
							{
								title, nowPlaying := GetTrackInfo()
								if title != "" {
									SendMessageToChannel(client, 32, "Station: "+title)
								}
								if nowPlaying != "" {
									SendMessageToChannel(client, 32, "Track: "+nowPlaying)
								}
							}
						}
					}
				}

				// SendMessageToChannel(client, 32, Welcome)
			}
		}
	}
}

func printOnlineClients(onlineClients map[int]*ts3.OnlineClient) {
	log.Printf("Active clients:")
	for id, client := range onlineClients {
		log.Printf("--> %d : %s", id, client.Nickname)
	}
}

func updateTimeChannel(client *ts3.Client) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tt := time.Now().In(time.FixedZone("Moscow", 3*60*60)).Format("15:04:05")
			localTime := fmt.Sprintf("[spacer.time]   TIME : %s", tt)
			СhangeChannelName(client, 22, localTime)
		}
	}
}

func updateChillOutChannel(client *ts3.Client) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	var titleActual, nowPlayingActual string
	re := regexp.MustCompile(`(?i)(emp|[^0-9A-Za-zа-яА-Я_ -])`)

	for {
		select {
		case <-ticker.C:
			title, nowPlaying := GetTrackInfo()
			re.ReplaceAllString(title, "")
			re.ReplaceAllString(nowPlaying, "")
			if title != titleActual && title != "" {
				titleActual = title
				SendMessageToChannel(client, 32, "Station: "+title)
			}
			if nowPlaying != nowPlayingActual && nowPlaying != "" {
				nowPlayingActual = nowPlaying
				SendMessageToChannel(client, 32, "Track: "+nowPlaying)
			}
		}
	}
}
