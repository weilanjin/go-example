package er

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[LogID: %d]: ", key))
	log.Printf("%#v", err)
	log.Printf("[%v] %v", key, message)
}

func TestErr(t *testing.T) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)
	err := runJob("1")
	if err != nil {
		msg := "There was an unexpected issue; please report this as a bug."
		if _, ok := err.(IntermediateErr); ok {
			msg = err.Error()
		}
		handleError(1, err, msg)
	}
}
