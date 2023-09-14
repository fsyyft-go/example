package sys

type (
	// Runnable 可运行的应用程序，包含有 Run 方法。
	Runnable interface {
		// Run 运行程序。
		Run()
	}
)
