package x

actions: {
	_def: #Def & {
		cwd: ""
	}

	action: #Action & {
		cwd: _def.cwd
	}
}

#Action: {
	x: string

	#Build
}

#Build: {
	cwd: string

	manifests: {
		_build: #Exec & {
			"cwd": cwd
		}
	}
}

#Exec: {
	cwd: string
	_bin: #Bin & {
		"cwd": cwd
	}
}

#Bin: {
	cwd: string

	_sys_info: #SysInfo & {
		"cwd": cwd
	}

	version: string | *"v3.17.0"
	os:      string | *"\(_sys_info.platform.os)"
	arch:    string | *"\(_sys_info.platform.arch)"

	_fetch: #Fetch & {
		url: "https://get.helm.sh/helm-\(version)-\(os)-\(arch).tar.gz"
	}

	result: _fetch.file
}

#SysInfo: $$type: name: "SysInfo"
#SysInfo: {
	cwd: string

	platform!: string @generated()
}

#Platform: $$type: name: "#Platform"
#Platform: {
	os!:   string
	arch!: string
}

#File: $$type: name: "File"
#File: $$file: id!:  string

#Fetch: $$type: name: "Fetch"
#Fetch: {
	url:  string
	file: #File @generated()
}

#Def: $$type: name: "Def"
#Def: {
	...
}
