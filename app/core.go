package app

import (
	"UnicBotGo/config"
    "UnicBotGo/ts"
)

func Run() {

	cfg := config.Get()
	tsQueryClient := ts.InitNewClient(&cfg)
	defer ts.CloseClient(tsQueryClient)
	ts.Login(tsQueryClient, &cfg)
	ts.UseVirtualServer(tsQueryClient)
	ts.SetNick(tsQueryClient, &cfg)
	ts.ClientMoveRequest(tsQueryClient, 32)

	ts.RegisterServerEvents(tsQueryClient)
	ts.RegisterTextChannelEvents(tsQueryClient)
	ts.RegisterTextPrivateEvents(tsQueryClient)

	go ts.ServerEventsListener(tsQueryClient)
	go ts.UpdateTimeChannel(tsQueryClient)
	go ts.UpdateChillOutChannel(tsQueryClient)

	select {}
}
