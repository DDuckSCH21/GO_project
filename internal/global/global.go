package global

import (
	"go_project/internal/models"
	"sync"
)

var DB = make(map[int]models.User)
var MyMute sync.Mutex
