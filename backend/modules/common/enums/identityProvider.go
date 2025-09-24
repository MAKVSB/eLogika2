package enums

type IdentityProviderEnum string

const (
	IdentityProviderInternal IdentityProviderEnum = "INT" // Standalone user
	IdentityProviderVSB      IdentityProviderEnum = "VSB" // Imported from VŠB Inbus
)

var IdentityProviderEnumAll = []IdentityProviderEnum{
	IdentityProviderInternal,
	IdentityProviderVSB,
}

func (w IdentityProviderEnum) TSName() string {
	switch w {
	case IdentityProviderInternal:
		return "INT"
	case IdentityProviderVSB:
		return "VSB"
	default:
		return "???"
	}
}
