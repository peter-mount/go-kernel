package strings

type StringHandler func(string) error

// With starts a new StringHandler
func With(b StringHandler) StringHandler {
	return b
}

// Then joins two string handlers so the left-hand side runs first then the right-hand side.
func (a StringHandler) Then(b StringHandler) StringHandler {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	return func(s string) error {
		err := a(s)
		if err != nil {
			return err
		}
		return b(s)
	}
}

// Do invokes a StringHandler with a specific string
func (a StringHandler) Do(s string) error {
	if a == nil {
		return nil
	}
	return a(s)
}
