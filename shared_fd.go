package daemonizer

import (
	"os"
	"strconv"
)

const (
	DAEMON_ENV_SIGNATURE          = "_DAEMONIZER__IS_DAEMON"
	SHARED_FILE_DESCRIPTOR_NUMBER = "_DAEMONIZER__SHARED_FD_NUMBER"
)

var (
	sharedFds []uintptr
	IsDaemon  bool
)

func init() {
	if os.Getenv(DAEMON_ENV_SIGNATURE) == "true" {
		IsDaemon = true
	} else {
		IsDaemon = false
	}

	if IsDaemon {
		fd_number, err := strconv.Atoi(os.Getenv(SHARED_FILE_DESCRIPTOR_NUMBER))
		if err != nil {
			fd_number = 0
		}
		sharedFds = make([]uintptr, fd_number)
		for i := 0; i < fd_number; i++ {
			sharedFds[i] = uintptr(i)
		}
	} else {
		sharedFds = make([]uintptr, 0, 5)
	}
}

func AddSharedFileDescriptor(fd uintptr) {
	if !IsDaemon {
		sharedFds = append(sharedFds, fd)
		os.Setenv(SHARED_FILE_DESCRIPTOR_NUMBER, strconv.Itoa(len(sharedFds)))
	}
}

func SharedFileDescriptors() []uintptr {
	return sharedFds
}
