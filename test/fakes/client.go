package fakes

type Client struct {
	*Conn
}

func NewClient() *Client {
	return &Client{
		Conn: NewConn(),
	}
}

func (c *Client) Connect(l *Listener) {
	l.ch <- c
	return
}
