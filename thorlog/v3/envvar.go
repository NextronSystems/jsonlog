package thorlog

type EnvironmentVariable struct {
	LogObjectHeader

	Variable string `json:"variable" textlog:"var"`
	Value    string `json:"value" textlog:"value"`
}

func (EnvironmentVariable) reportable() {}

const typeEnvironmentVariable = "environment variable"

func init() { AddLogObjectType(typeEnvironmentVariable, &EnvironmentVariable{}) }

func NewEnvironmentVariable(variable string, value string) *EnvironmentVariable {
	return &EnvironmentVariable{
		LogObjectHeader: LogObjectHeader{
			Type: typeEnvironmentVariable,
		},
		Variable: variable,
		Value:    value,
	}
}
