package threading

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoSafe(t *testing.T) {
	// 设置log输出不做任何操作，直接返回success
	log.SetOutput(ioutil.Discard)

	i := 0
	defer func() {
		assert.Equal(t, 1, i)
	}()

	ch := make(chan struct{})
	GoSafe(func() {
		defer func() {
			ch <- struct{}{}
		}()
		panic("panic!")
	})

	<-ch
	i++
}
