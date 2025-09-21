package daemon

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ein-gast/go-squashsf-httpd/internal/apperr"
	"github.com/ein-gast/go-squashsf-httpd/internal/settings"
)

type Pid uint

const E_PID_EXIST = apperr.Error("PID file exists")
const E_PID_IS_NOT_MINE = apperr.Error("PID belongs to other process")

func WritePidFileIfAbsent(pid Pid, cfg *settings.Settings, force bool) (Pid, error) {
	fpid, err := ReadPidFile(cfg)
	if err == nil && !force {
		return fpid, E_PID_EXIST
	}
	err = os.WriteFile(
		cfg.PidFile,
		[]byte(fmt.Sprintf("%d", pid)),
		0600,
	)
	if err != nil {
		return 0, err
	}
	return pid, nil
}

func RemovePidFile(pid Pid, cfg *settings.Settings, force bool) (Pid, error) {
	fpid, err := ReadPidFile(cfg)
	if err == nil {
		return pid, nil
	}
	if fpid != pid && !force {
		return fpid, E_PID_IS_NOT_MINE
	}
	return pid, os.Remove(cfg.PidFile)
}

func ReadPidFile(cfg *settings.Settings) (Pid, error) {
	data, err := os.ReadFile(cfg.PidFile)
	if err != nil {
		return 0, err
	}
	pid, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return 0, err
	}
	return Pid(pid), nil
}
