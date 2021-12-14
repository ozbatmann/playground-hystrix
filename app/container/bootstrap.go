package container

import (
	"playground/app/clients"
)


func Init() *clients.LimitApiClientStruct{
	return clients.NewTestApiClient()

}