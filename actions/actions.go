package actions

type Group interface {
	AddAction(name, description string, svgInput []byte, runFunc interface{})
	AddManualTrigger(name, description string, svgIcon []byte, runFunc interface{})
}

type Provider interface {
	AddGroup(name, description string, svgIcon []byte) Group
}

type ProviderRegistry interface {
	AddProvider(name, description string, svgIcon []byte) Provider
}

var Registry ProviderRegistry
