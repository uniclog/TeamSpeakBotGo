package ts

import (
    "UnicBotGo/vlc"
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

func ServerEventsListener(client *ts3.Client) {
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
                //log.Println(notification)
                clientInfo := GetClientInfo(notification.Data)
                log.Printf("Client joined: ID=%d, Name=%s", clientInfo.ID, clientInfo.Nickname)
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
                            vlc.NextStation()
                        case strings.HasPrefix(msg.Message, VCPrev):
                            vlc.PrevStation()
                        case strings.HasPrefix(msg.Message, VCInfo):
                            {
                                title, nowPlaying := vlc.GetTrackInfo()
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

func UpdateTimeChannel(client *ts3.Client) {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            tt := time.Now().In(time.FixedZone("Moscow", 3*60*60)).Format("15:04")
            localTime := fmt.Sprintf("[spacer.time]   TIME : %s", tt)
            ChangeChannelName(client, 22, localTime)
        }
    }
}

func UpdateChillOutChannel(client *ts3.Client) {
    ticker := time.NewTicker(15 * time.Second)
    defer ticker.Stop()

    var titleActual, nowPlayingActual string
    re := regexp.MustCompile(`(?i)(emp|[^0-9A-Za-zА-я_ -])`)

    for {
        select {
        case <-ticker.C:
            title, nowPlaying := vlc.GetTrackInfo()
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
