package kernel

type Daemon struct {
	daemon bool
}

func (d *Daemon) SetDaemon() {
	d.daemon = true
}

func (d *Daemon) IsDaemon() bool {
	return d.daemon
}
