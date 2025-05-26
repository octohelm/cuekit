package main

action: {
	_local: #Local & {}

	_target: {
		os: ["linux"]
		arch: ["amd64", "arm64"]
	}

	steps: [
		#Copy & {},
	]

	_build: {
		for _os in _target.os for _arch in _target.arch {
			"\(_os)/\(_arch)": #Build & {
				input: _local.dir

				"steps": [
					for s in steps {
						s
					},

					#InstallPackage & {
						packages: {
							"git": _
						}
					},

					{
						input: _

						_exec: #Exec & {
							cwd: input
							cmd: [
								"env",
							]
						}

						output: _exec.cwd
					},
				]
			}
		}
	}

	_archive: {
		for _os in _target.os for _arch in _target.arch {
			"\(_os)/\(_arch)": {
				_built: _build["\(_os)/\(_arch)"].output

				_copy: #Copy & {
					input: _built
				}

				output: _copy.output
			}
		}
	}

	result: _archive["linux/arm64"].output
}

#InstallPackage: {
	input: #WorkDir
	packages: [pkgName=string]: _

	_install: #Build & {
		"input": input

		"steps": [
			if len(packages) > 0 {
				for p, v in packages {
					{
						input: _

						_exec: #Exec & {
							cwd: input
							cmd: ["echo", p]
						}

						output: _exec.cwd
					}
				}
			},
		]
	}

	output: _install.output
}
