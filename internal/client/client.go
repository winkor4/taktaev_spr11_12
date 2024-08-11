package client

import (
	"errors"
)

type Client struct {
	login      string
	password   string
	runAddress string
}

type Config struct {
	Login      string
	Password   string
	RunAddress string
}

var dataTypes = map[string]bool{
	"logpass": true,
}

const (
	regCommand    = "reg"
	authCommand   = "auth"
	addCommand    = "add"
	getCommand    = "get"
	delCommand    = "del"
	updateCommand = "update"
	listCommand   = "list"
)

func NewClient(cfg Config) *Client {
	return &Client{
		login:      cfg.Login,
		password:   cfg.Password,
		runAddress: cfg.RunAddress,
	}
}

func (c *Client) Do(command string) error {
	switch command {
	case regCommand:
		return c.register()
	case authCommand:
		return c.auth()
	case addCommand:
		reqData, err := readAddFlags()
		if err != nil {
			return err
		}
		return c.addContent(reqData)
	case getCommand:
		name, err := readGetDelFlags()
		if err != nil {
			return err
		}
		return c.getContent(name)
	case delCommand:
		name, err := readGetDelFlags()
		if err != nil {
			return err
		}
		return c.delContent(name)
	case updateCommand:
		reqData, err := readAddFlags()
		if err != nil {
			return err
		}
		return c.updateContent(reqData)
	case listCommand:
		return c.listContent()
	default:
		return errors.New("неизвестная команда")
	}
}
