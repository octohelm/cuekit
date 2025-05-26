package x

def: #Def & {
	a: string
}

actions: {
	x: #Build & {
		input: def.a
	}
}

#Build: X = {
	input: string

	os: ["linux"]
	arch: ["amd64"]

	_build: {
		for _os in os for _arch in arch {
			"\(_os)/\(_arch)": #Copy & {
				"input": X.input
			}
		}
	}

	_archive: {
		for _os in os for _arch in arch {
			"\(_os)/\(_arch)": {
				_built: _build["\(_os)/\(_arch)"].output

				_copy: #Copy & {
					"input": _built
				}

				output: _copy.output
			}
		}
	}

	output: {
		for _os in os for _arch in arch {
			"\(_os)/\(_arch)": _build["\(_os)/\(_arch)"].output
		}
	}
}

#Copy: "$$type": name: "Copy"
#Copy: "$$task": true
#Copy: {
	input:  string
	output: string @generated()
}

#Def: "$$type": name: "Def"
#Def: "$$task": true
#Def: {
	...
}
