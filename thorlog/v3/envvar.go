package thorlog

type EnvironmentVariable struct {
	LogObjectHeader

	File    string   `json:"file,omitempty" textlog:"file,omitempty"`
	Process *Process `json:"process,omitempty" textlog:"process,expand,omitempty"`

	Variable string `json:"variable" textlog:"var"`
	Value    string `json:"value" textlog:"value"`
}

const typeEnvironmentVariable = "environment variable"

func init() { AddLogObjectType(typeEnvironmentVariable, &EnvironmentVariable{}) }

func NewEnvironmentVariable(variable string, value string) *EnvironmentVariable {
	return &EnvironmentVariable{
		LogObjectHeader: LogObjectHeader{
			Type:    typeEnvironmentVariable,
			Summary: variable + "=" + value,
		},
		Variable: variable,
		Value:    value,
	}
}
