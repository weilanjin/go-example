package err

import (
	"log"
	"testing"

	_errors "errors"

	"github.com/pkg/errors"
)

func TestErrors(t *testing.T) {
	err := errors.New("test error")
	log.Printf("%+v", err)
}

func TestWebServer(t *testing.T) {
	err := Api()
	log.Printf("%+v", err)
}

func Api() error {
	return Service()
}

func Service() error {
	return Dao()
}

func Dao() error {
	err := _errors.New("test xxxx error")
	return errors.WithStack(err)
}