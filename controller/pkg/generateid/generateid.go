package generateid

import (
	"math/rand"
	"time"
)

func GenerateID() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(100)
}
