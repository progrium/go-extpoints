
implements("ProgramObserver")

function ProgramStarted() {
	console.log("Yay, builtin plugin! It's starting!")
}

function ProgramEnded() {
	console.log("Yay, builtin plugin! It's over!")
}