package intermediate

import (
	"github.com/go_concurrency/error-propagation/lowlevel"
	"github.com/go_concurrency/error-propagation/error-prop"
	"os/exec"
)

type IntermediateErr struct {
	error
}


func RunJob(id string) error {
	const jobBinPath = "/bad/job/binary"
	isExecutable, err := lowlevel.IsGloballyExec(jobBinPath)
	if err != nil {
		return err
	} else if !isExecutable {
		return error_prop.WrapError(nil, "job binary is not executable")
	}
	return exec.Command(jobBinPath, "--id="+id).Run()
}