package util

type Validator struct {
	Error map[string]string
}

func NewValidator() *Validator {
	return &Validator{}
}
func (v *Validator) Check() {

}
