package main

action: {
	_local: #Local & {}

	_env: #Env & {
		KEY: string
	}

	_exec: #Exec & {
		cwd: _local.dir
		cmd: [
			"env",
		]
		env: {
			KEY: _env.KEY
		}
	}

	result: _exec.stdout
}
