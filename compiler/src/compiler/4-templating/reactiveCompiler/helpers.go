package reactiveCompiler

import "hudson-newey/2web/src/models"

func getUniqueSelectors(input []*models.ReactiveProperty) []*models.ReactiveProperty {
	seen := make(map[string]bool)
	result := []*models.ReactiveProperty{}

	for _, item := range input {
		if !seen[item.Node.Selector] {
			seen[item.Node.Selector] = true
			result = append(result, item)
		}
	}

	return result
}
