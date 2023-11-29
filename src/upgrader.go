package lib

type Upgrader struct{}

func NewUpgrader() *Upgrader {
	return &Upgrader{}
}

func (u *Upgrader) Upgrade(client *Client) {

}
