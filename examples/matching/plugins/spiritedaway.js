
implements("OutputRenderer")

function Match(pattern) {
	if (pattern == "miyazaki" || pattern == "ghibli") {
		return true
	}
	return false
}

function Output() {
	return "The tunnel led Chihiro to a mysterious town..."
}