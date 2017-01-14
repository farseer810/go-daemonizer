package daemonizer

import (
	"os"
	"syscall"
)

func Daemonize() (pid int, err error) {
	os.Setenv(DAEMON_ENV_SIGNATURE, "true")
	execSpec := &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: SharedFileDescriptors(),
	}
	pid, err = syscall.ForkExec(os.Args[0], os.Args, execSpec)
	return
}
