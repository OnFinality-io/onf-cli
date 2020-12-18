package printer

type JSONYamlPrintFlags struct {
}

func (f *JSONYamlPrintFlags) AllowedFormats() []string {
	if f == nil {
		return []string{}
	}
	return []string{"json", "yaml"}
}

func NewJSONYamlPrintFlags() *JSONYamlPrintFlags {
	return &JSONYamlPrintFlags{}
}
