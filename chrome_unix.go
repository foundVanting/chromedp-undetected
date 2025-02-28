//go:build linux

// Package chromedpundetected provides a chromedp context with an undetected
// Chrome browser.
package chromedpundetected

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/chromedp/chromedp"
)

func headlessOpts() (opts []chromedp.ExecAllocatorOption, cleanup func() error, err error) {
	// Create virtual display
	frameBuffer, err := newFrameBuffer("1920x1080x24")
	if err != nil {
		return nil, nil, err
	}
	cleanup = frameBuffer.Stop
	opt := chromedp.ModifyCmdFunc(func(cmd *exec.Cmd) {
		cmd.Env = append(cmd.Env, "DISPLAY=:"+frameBuffer.Display)
		cmd.Env = append(cmd.Env, "XAUTHORITY="+frameBuffer.AuthPath)

		// Default modify command per chromedp
		if _, ok := os.LookupEnv("LAMBDA_TASK_ROOT"); ok {
			// do nothing on AWS Lambda
			return
		}

		if cmd.SysProcAttr == nil {
			cmd.SysProcAttr = new(syscall.SysProcAttr)
		}

		// When the parent process dies (Go), kill the child as well.
		cmd.SysProcAttr.Pdeathsig = syscall.SIGKILL
	})
	return []chromedp.ExecAllocatorOption{opt}, cleanup, nil
}
