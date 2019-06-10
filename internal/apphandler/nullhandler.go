package apphandler

// NullHandler do nothing
type NullHandler struct {
	Handler AppHandler
	App     *App
}

// ManageApp - don't do anything
func (nh *NullHandler) ManageApp() {
	return
}
