package router

import (
	"github.com/vallen-genar/router/entity"
)

type Config struct {
	Prefix string
	UriVarExp string
	HttpNotFound entity.IController
}
