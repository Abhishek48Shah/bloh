package util

type Validator struct {
	Error map[string]string
}

func NewValidator() *Validator {
	return &Validator{Error: map[string]string{}}
}
func (v *Validator) createError(key string, message string) {
	if _, ok := v.Error[key]; !ok {
		v.Error[key] = message
	}
}
func (v *Validator) Ok() bool {
	if len(v.Error) == 0 {
		return true
	}
	return false
}
func (v *Validator) Check(key string, ok bool, message string) {
	if !ok {
		v.createError(key, message)
	}
}
