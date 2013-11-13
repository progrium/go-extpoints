
implements("PluggableBackend")

String.prototype.repeat = function(num) {
    return new Array(num + 1).join(this);
}

function Name() {
	return "powerhorse"
}

function Process(a, b) {
	return ([a, " ", b, " "].join("")).toUpperCase().repeat(5)
}