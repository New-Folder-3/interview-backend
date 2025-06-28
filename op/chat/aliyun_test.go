package chat

import (
	"fmt"
	"github.com/pkg/errors"
	"testing"
)

func TestAliyunChat(t *testing.T) {
	err := errors.WithMessage(errors.New("init"), "test")
	fmt.Println(err.Error())
}
