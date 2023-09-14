package goroutine

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	oeTesting "github.com/fsyyft-go/example/pkg/testing"
)

func TestGetGoID(t *testing.T) {
	assertions := assert.New(t)
	t.Run("测试获取 GoroutineID", func(t *testing.T) {

		var wg sync.WaitGroup
		var idOuter, idInternal int64
		wg.Add(1)
		idOuter = GetGoID()
		go func() {
			idInternal = GetGoID()
			wg.Done()
		}()
		wg.Wait()
		// 值每次都不一样，有需要的情况可以打印出来查看。
		assertions.NotEqual(idOuter, idInternal)
		// 在没有复用的情况下，里的一般会比外的大。
		assertions.LessOrEqual(idOuter, idInternal)
		oeTesting.Println(idInternal, idOuter)
	})
}
