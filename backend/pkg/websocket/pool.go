package websocket


// máme tu ve struktuře několik kanálů do kterých budu chtít zapisovat
// při zavolání serveWS 
type Pool struct {
	Register	chan *Client
	Unregister	chan *Client
	Clients		map[*Client]bool
	Broadcast  chan Message
}

NewPool()