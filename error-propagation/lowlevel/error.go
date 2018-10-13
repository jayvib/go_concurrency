package lowlevel

import (
	"os"
	"github.com/go_concurrency/error-propagation/error-prop"
)

type LowLevelErr struct {
	error
}

func IsGloballyExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, LowLevelErr{
			error: error_prop.WrapError(err, err.Error()),
		}
	}
	return info.Mode().Perm()&0100 == 0100, nil
}
