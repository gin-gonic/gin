package main

import "github.com/dustin/go-broadcast"

var roomChannels = make(map[string]broadcast.Broadcaster)

func openListener(roomid string) chan interface{} {
	listener := make(chan interface{})
	room(roomid).Register(listener)
	return listener
}

func closeListener(roomid string, listener chan interface{}) {
	room(roomid).Unregister(listener)
	close(listener)
}

func room(roomid string) broadcast.Broadcaster {
	b, ok := roomChannels[roomid]
	if !ok {
		b = broadcast.NewBroadcaster(10)
		roomChannels[roomid] = b
	}
	return b
}
