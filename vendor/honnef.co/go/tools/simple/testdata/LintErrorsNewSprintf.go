package pkg

import (
	"errors"
	"fmt"
)

func fn() {
	_ = fmt.Errorf("%d", 0)
	_ = errors.New("")
	_ = errors.New(fmt.Sprintf("%d", 0)) // MATCH "should use fmt.Errorf"
}
