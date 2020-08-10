package router

import "github.com/g8y3e/router/entity"

type Config struct {
	Prefix string
	UriVarExp string
	HttpNotFound entity.IController
}
