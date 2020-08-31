package resource

type Process struct{}

func (m *Process) GetInfo() (*ResourceResult, error) {
	cmdRes, err := execSysCmd(5, "top", "-n", "1")
	if err != nil {
		return nil, err
	}
	data, _ := parseTOP(cmdRes, `^tasks:.*`, `(\d+)`)

	n, _ := rtToName(RTPRO)
	return NewResourceResult(n, data), nil

}
