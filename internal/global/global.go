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
	fmt.Printf("Test DB.Del; id=%d\n", id)

	glob.MyMute.Lock()
	defer glob.MyMute.Unlock()

	// _, ok := glob.DBglobal[id] //Лишняя какая-то проверка
	// if ok {
	delete(glob.DBglobal, id)
	return true
	// } else {
	// 	return false
	// }
}

func (glob *DBst) Set(id int, usr models.User) bool {
	fmt.Printf("Test DB.Set; id=%d\n", id)

	glob.MyMute.Lock()
	defer glob.MyMute.Unlock()

	// _, ok := glob.DBglobal[id] //
	// if ok {
	glob.DBglobal[id] = usr
	fmt.Printf("Test DB.Set_2; true\n", id)

	return true
	// } else {
	// 	fmt.Printf("Test DB.Set_2; false\n", id)

	// 	return false
	// }
}

func (glob *DBst) GetAll() bytes.Buffer {
	fmt.Printf("Test DB.GetAll\n")

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

func (glob *DBst) Get(id int) (usr models.User, status bool) {
	fmt.Printf("Test DB.Get; id=%d\n", id)
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
	fmt.Printf("Test DB.IsEmpty\n")

	if len(glob.DBglobal) != 0 {
		return false
	} else {
		return true
	}
}

func (glob *DBst) GetNewKey() int {
	fmt.Printf("Test DB.GetNewKey_in\n")

	glob.MyMute.Lock()
	defer glob.MyMute.Unlock()
	if len(glob.DBglobal) != 0 {
		maxKey := math.MinInt
		for num := range glob.DBglobal {
			if maxKey < num {
				maxKey = num
			}
		}
		fmt.Printf("Test DB.GetNewKey_out ret=%d\n", maxKey+1)
		return maxKey + 1
	}
	fmt.Printf("Test DB.GetNewKey_out ret=1\n")
	return 1
}
