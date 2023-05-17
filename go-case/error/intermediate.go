package er

import "os/exec"

type IntermediateErr struct {
	error
}

func runJob(id string) error {
	const jobBinPath = "bad/job/path"
	isExec, err := isGloballyExec(jobBinPath)
	if err != nil {
		return IntermediateErr{WrapError(
			err,
			"cannot run job %q: requisite binaries not available",
			id,
		)}
	} else if !isExec {
		return WrapError(
			nil,
			"cannot run job %q: requisite binaries are not executable",
			id,
		)
	}
	return exec.Command(jobBinPath, "--id="+id).Run()
}
