package main

const (
	Name    = "UnicBotGo"
	Author  = "D.V."
	Version = "0.0.0"
)

func main() {

	config := loadConfig()
	tsQueryClient := InitNewClient(&config)
	defer closeClient(tsQueryClient)
	Login(tsQueryClient, &config)
	UseVirtualServer(tsQueryClient)
	SetNick(tsQueryClient, &config)
	ClientMoveRequest(tsQueryClient, 32)

	RegisterServerEvents(tsQueryClient)
	RegisterTextChannelEvents(tsQueryClient)
	RegisterTextPrivateEvents(tsQueryClient)

	go serverEventsListener(tsQueryClient)
	go updateTimeChannel(tsQueryClient)
	go updateChillOutChannel(tsQueryClient)

	select {}
}
