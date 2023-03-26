package cli

type validate struct {
	templatePath string
}

type plan struct {
	templatePath string
	paramFile    string
}

type apply struct {
	autoApprove  bool
	templatePath string
	paramFile    string
}

type destroy struct {
	autoApprove  bool
	templatePath string
}

type output struct{}

type version struct {
	version string
}

type drift struct {
	stackName string
}
