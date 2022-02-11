package kernel

type Daemon struct {
	daemon    bool
	webserver bool
}

func (d *Daemon) SetDaemon() {
	d.daemon = true
}

func (d *Daemon) ClearDaemon() {
	d.daemon = false
}

func (d *Daemon) IsDaemon() bool {
	return d.webserver || d.daemon
}

func (d *Daemon) SetWebserver() {
	d.webserver = true
}

func (d *Daemon) IsWebserver() bool {
	return d.webserver
}
