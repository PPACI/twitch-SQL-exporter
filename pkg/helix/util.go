package helix

func getCleanMap(input map[string]string) map[string]string {
	output := map[string]string{}
	for k, v := range input {
		if v != "" {
			output[k] = v
		}
	}
	return output
}
