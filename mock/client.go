package mock

import (
	"github.com/MrSametBurgazoglu/enterprise/client"
	"github.com/pashagolub/pgxmock/v4"
)

var _ client.DatabaseClient = (*Client)(nil)

type Client struct {
	pgxmock.PgxPoolIface
}

func (c *Client) BeginHook() {
	return
}

func (c *Client) EndHook() {
	return
}

func NewMockClient() *Client {
	conn, _ := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	return &Client{PgxPoolIface: conn}
}
