// Package tempclient
// Mode Name: temporal连接器
// Mode Desc: 模块功能描述
package tempclient

import "go.temporal.io/sdk/client"

type TClient struct {
	c client.Client
	o client.Options
}

func NewClient() *TClient {
	return &TClient{}
}

func (tc *TClient) Init(opts *client.Options) error {
	tc.o = *opts
	var err error
	tc.c, err = client.Dial(tc.o)
	if err != nil {
		return err
	}
	return nil
}

func (tc *TClient) Close() {
	tc.c.Close()
}

func (tc *TClient) GetConn() client.Client {
	return tc.c
}
