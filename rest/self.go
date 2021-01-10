package rest

// Self returns the full request uri based on the headers of the request and
// the supplied path.
func (r *Rest) Self(path string) string {
	scheme := r.GetHeader("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http"
	}

	host := r.GetHeader("X-Forwarded-Host")
	if host == "" {
		host = r.Request().Host
	} else {
		// We have a host so check the port
		port := r.GetHeader("X-Forwarded-Port")
		if port != "" && port != "80" && port != "443" {
			host = host + ":" + port
		}
	}
	return scheme + "://" + host + path
}

// SelfRequest returns the full request uri based on the headers of the request
// and the RequestURI
func (r *Rest) SelfRequest() string {
	return r.Self(r.Request().RequestURI)
}
