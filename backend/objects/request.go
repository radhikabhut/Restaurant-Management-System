package objects

import "context"

type RequestObject struct {
	Ctx      context.Context
	Request  interface{}
	Response interface{}
}
