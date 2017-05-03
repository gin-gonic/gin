package rest

import (
	"github.com/gin-gonic/gin"
)

// All of the methods are the same type as HandlerFunc
// if you don't want to support any methods of CRUD, then don't implement it
type CreateSupported interface {
	CreateHandler(*gin.Context)
}
type ListSupported interface {
	ListHandler(*gin.Context)
}
type TakeSupported interface {
	TakeHandler(*gin.Context)
}
type UpdateSupported interface {
	UpdateHandler(*gin.Context)
}
type DeleteSupported interface {
	DeleteHandler(*gin.Context)
}

// It defines
//   POST: /path
//   GET:  /path
//   PUT:  /path/:id
//   POST: /path/:id
func CRUD(group *gin.RouterGroup, path string, resource interface{}) {
	if resource, ok := resource.(CreateSupported); ok {
		group.POST(path, resource.CreateHandler)
	}
	if resource, ok := resource.(ListSupported); ok {
		group.GET(path, resource.ListHandler)
	}
	if resource, ok := resource.(TakeSupported); ok {
		group.GET(path+"/:id", resource.TakeHandler)
	}
	if resource, ok := resource.(UpdateSupported); ok {
		group.PUT(path+"/:id", resource.UpdateHandler)
	}
	if resource, ok := resource.(DeleteSupported); ok {
		group.DELETE(path+"/:id", resource.DeleteHandler)
	}
}
