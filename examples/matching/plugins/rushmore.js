
implements("OutputRenderer")

function Match(pattern) {
	if (["wes", "jason", "bill"].indexOf(pattern) > -1) {
		return true
	}
	return false
}

function Output() {
	return "She was my Rushmore"
}