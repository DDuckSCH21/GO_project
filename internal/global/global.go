package global

import (
	"bytes"
	"fmt"
	"go_project/internal/models"
	"math"
	"sync"
)

var DB = DBst{DBglobal: make(map[int]models.User)} //{DBglobal : make(map[int]models.User)}

type DBst struct {
	MyMute   sync.Mutex
	DBglobal map[int]models.User
}

func (glob *DBst) Del(id int) bool {

	glob.MyMute.Lock()
	defer glob.MyMute.Unlock()

	delete(glob.DBglobal, id)
	return true
}

func (glob *DBst) Set(id int, usr models.User) bool { //Может лучше возвращать указатель на User??

	glob.MyMute.Lock()
	defer glob.MyMute.Unlock()

	glob.DBglobal[id] = usr
	return true
}

func (glob *DBst) GetAll() bytes.Buffer {

	glob.MyMute.Lock()
	defer glob.MyMute.Unlock()

	var buf bytes.Buffer

	if len(glob.DBglobal) != 0 {
		for ind, val := range glob.DBglobal {
			fmt.Fprintf(&buf, "User ID = %d: %v\n", ind, val)
		}
	}
	return buf
}

func (glob *DBst) Get(id int) (usr models.User, status bool) { //Может лучше возвращать указатель на User??

	glob.MyMute.Lock()
	defer glob.MyMute.Unlock()

	user, ok := glob.DBglobal[id]
	if ok {
		return user, true
	} else {
		return user, false //TODO Проверить, что же я возвращаю, похоже на фигню
	}
}

func (glob *DBst) IsEmpty() bool {

	if len(glob.DBglobal) != 0 {
		return false
	} else {
		return true
	}
}

func (glob *DBst) GetNewKey() int {

	glob.MyMute.Lock()
	defer glob.MyMute.Unlock()
	if len(glob.DBglobal) != 0 {
		maxKey := math.MinInt
		for num := range glob.DBglobal {
			if maxKey < num {
				maxKey = num
			}
		}
		return maxKey + 1
	}
	return 1
}
