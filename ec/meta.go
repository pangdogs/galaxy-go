package ec

type Meta map[string]any

func (m Meta) Get(k string) any {
	if m == nil {
		return nil
	}
	return m[k]
}

func (m Meta) Range(fun func(k string, v any) bool) {
	if fun == nil {
		return
	}

	for k, v := range m {
		if !fun(k, v) {
			return
		}
	}
}
