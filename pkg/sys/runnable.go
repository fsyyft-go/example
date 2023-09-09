package sys

import (
	"context"
)

type (
	// Runnable 可运行的应用程序，包含有 Run 方法。
	Runnable interface {
		// Run 运行程序。
		Run()
	}

	// RunnableWithContext 可运行的应用程序，包含有 RunContext 方法。
	RunnableWithContext interface {
		// RunContext 运行程序，带有上下文信息。
		RunContext(ctx context.Context)
	}
)
