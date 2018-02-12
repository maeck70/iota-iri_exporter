package main




func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func btof(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func actify(b float64) string {
	if b == 1{
		return "Active"
	}
	return "Inactive"
}