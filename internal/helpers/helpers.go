package helpers

func In[T comparable](value T, checklist []T) bool {
	for i := range checklist {
		if value == checklist[i] {
			return true
		}
	}
	return false
}
