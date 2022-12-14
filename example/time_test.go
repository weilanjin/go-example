package test

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeNew(t *testing.T) {
	fmt.Println(time.Now().Location())
}
