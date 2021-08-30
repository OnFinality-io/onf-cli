package models

type Action string

const (
	ADD     Action = "add"
	REPLACE Action = "replace"
	REMOVE  Action = "remove"
	UPSERT  Action = "upsert"
)

type Category string

const (
	ARG Category = "arg"
	ENV Category = "env"
	VAR Category = "var"
)

type VarValueType string

const (
	String      VarValueType = "string"
	StringArray VarValueType = "array"
	Empty       VarValueType = "empty"
	File        VarValueType = "file"
)

type RuleValueInputType string

const (
	Variable RuleValueInputType = "variable"
	Template RuleValueInputType = "template"
)

type NodeType string

const (
	Full      NodeType = "full"
	Validator NodeType = "validator"
	Collator  NodeType = "collator"
	Archive   NodeType = "archive"
	Light     NodeType = "light"
)

type RuleOperation struct {
	Action Action `json:"action"`
}

type Config struct {
	Operations map[NodeType]*RuleOperationCollection `json:"operations"`
}

type Metadata struct {
	Chainspec    *string `json:"chainspec"`
	ImageVersion *string `json:"imageVersion"`
}

type Vars []*Var

func (vars *Vars) Merge(newVars []*Var) {
	m := make(map[string]*Var)
	for _, v := range *vars {
		m[v.Payload.Key] = v
	}
	for _, newVar := range newVars {
		if newVar.Action == REMOVE {
			delete(m, newVar.Payload.Key)
		} else {
			m[newVar.Payload.Key] = newVar
		}
	}
	values := make([]*Var, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	*vars = values
}

type Args []*Arg

func (args *Args) Merge(newArgs []*Arg) {
	oriSectionMap := make(map[*int]*[]*Arg)
	for _, v := range *args {
		if argsGroupSecion, ok := oriSectionMap[v.Payload.Section]; ok {
			*argsGroupSecion = append(*argsGroupSecion, v)
		} else {
			oriSectionMap[v.Payload.Section] = &[]*Arg{v}
		}

	}

	genKey := func(args *Arg) string {
		//key:=fmt.Sprintf("%s_%d",args.Payload.Key,args.Payload.Section)
		key := args.Payload.Key
		return key
	}

	values := make([]*Arg, 0, len(*args)+len(newArgs))
	for section, sectionArgs := range oriSectionMap {
		m := make(map[string]*Arg)
		for _, v := range *sectionArgs {
			m[genKey(v)] = v
		}
		for _, newVar := range newArgs {
			if newVar.Payload.Section == nil || section == nil {
				if newVar.Action == REMOVE {
					delete(m, genKey(newVar))
				} else {
					m[genKey(newVar)] = newVar
				}
			} else {
				if *newVar.Payload.Section == *section {
					if newVar.Action == REMOVE {
						delete(m, genKey(newVar))
					} else {
						m[genKey(newVar)] = newVar
					}
				}
			}
		}
		for _, v := range m {
			values = append(values, v)
		}
	}

	*args = values
}

type Envs []*Env

func (envs *Envs) Merge(newEnvs []*Env) {
	m := make(map[string]*Env)
	for _, v := range *envs {
		m[v.Payload.Key] = v
	}
	for _, newVar := range newEnvs {
		if newVar.Action == REMOVE {
			delete(m, newVar.Payload.Key)
		} else {
			m[newVar.Payload.Key] = newVar
		}
	}
	values := make([]*Env, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	*envs = values
}

type RuleOperationCollection struct {
	Arg Args `json:"arg"`
	Var Vars `json:"var"`
	Env Envs `json:"env"`
}

type Options struct {
	Overwritable bool `json:"overwritable"`
}

type FileTypeValue struct {
	Source      *string `json:"source"`
	Destination string  `json:"destination"`
}

type VarValue struct {
	Payload   interface{}  `json:"payload"`
	ValueType VarValueType `json:"valueType"`
}
type VarModel struct {
	Key      string    `json:"key"`
	Value    *VarValue `json:"value"`
	Options  Options   `json:"options"`
	Category Category  `json:"category" default:"var"`
}
type Var struct {
	RuleOperation
	Payload *VarModel `json:"payload"`
}

type ValueModel struct {
	Payload   interface{}        `json:"payload"`
	InputType RuleValueInputType `json:"inputType"`
}
type ArgModel struct {
	Key      string      `json:"key"`
	Value    *ValueModel `json:"value"`
	Options  Options     `json:"options"`
	Section  *int        `json:"section,omitempty"`
	Category Category    `json:"category" default:"arg"`
}
type Arg struct {
	RuleOperation
	Payload *ArgModel `json:"payload"`
}

type EnvModel struct {
	Key      string      `json:"key"`
	Value    *ValueModel `json:"value"`
	Options  Options     `json:"options"`
	Category Category    `json:"category" default:"env"`
}
type Env struct {
	RuleOperation
	Payload EnvModel `json:"payload"`
}
