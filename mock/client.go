package mock

import (
	"context"

	"github.com/MrSametBurgazoglu/enterprise/client"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

var _ client.DatabaseClient = (*Client)(nil)

type Client struct {
	pgxmock.PgxPoolIface
}

func (c *Client) NewTransaction(context.Context, ...pgx.TxOptions) (client.DatabaseTransactionClient, error) {
	return c, nil
}

func (c *Client) Exit(){

}

func NewMockClient() *Client {
	conn, _ := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	return &Client{PgxPoolIface: conn}
}
