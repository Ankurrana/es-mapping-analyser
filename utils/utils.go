package utils

func Some(usecases, list []string) bool {
	// If atleast one element is same in usecases and list, then return true!
	// otherwise false
	for _, use_case := range usecases {
		for _, item := range list {
			if use_case == item {
				return true
			}
		}
	}
	return false
}

func Contains(subject string, arrList []string) bool {
	for _, item := range arrList {
		if item == subject {
			return true
		}
	}
	return false
}
