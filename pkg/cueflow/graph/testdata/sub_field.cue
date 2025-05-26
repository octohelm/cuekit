package x

actions: {
	_pull: #Pull & {}

	result: _pull.output.rootfs
}

#Container: "$$type": name: "Container"
#Container: {
	$$container: id!: string

	rootfs!:   #Fs
	platform!: string
}

#Fs: $$type: name: "Fs"
#Fs: $$fs: id!:    string

#Pull: $$type: name: "Pull"
#Pull: {
	output: #Container @generated()
}
