package main

import (
	"dagger.io/dagger"
	"dagger.io/dagger/core"
	"universe.dagger.io/docker"
	"universe.dagger.io/bash"
	"universe.dagger.io/alpha/kubernetes/kapp"
)

app: string | *"kubegraph"

#PrepareIntegrationTest: {
	srcTest: dagger.#FS

	registry_url: string

	build: docker.#Build & {
		steps: [
			docker.#Pull & {
				source: "python:3.9-alpine"
			},
			docker.#Copy & {
				contents: srcTest
				dest:     "/test"
			},
			docker.#Run & {
				command: {
					name: "pip"
					args: ["install", "-r", "requirements.txt"]
				}
				workdir: "/test"
			},
			docker.#Set & {
				config: {
					cmd: ["pytest"]
					workdir: "/test"
				}
			},
		]
	}

	push: docker.#Push & {
		image: build.output
		dest:  "\(registry_url)/it-test:ci"
	}

	output: push.result
}

#DoIntegrationTest: {
	job:        dagger.#FS
	kc:         dagger.#Secret
	imageRef:   core.#Ref

	deployment: kapp.#Deploy & {
		app:        "it-test"
		fs:         job
		kubeConfig: kc
		file:       "./it-test-job.yaml"
	}
}

#Deploy: {
	imageRef: string

	source: dagger.#FS

	kc: dagger.#Secret

	_template: bash.#RunSimple & {
        script: contents: "sed s#IMAGE_REF#$IMAGE_REF#g /source/deploy.template.yaml > /tmp/deploy.yaml"
        env: IMAGE_REF: imageRef
        mounts: "/source": {
			dest:     "/source"
			contents: source
		}
	}

	apply: kapp.#Deploy & {
		app:        "kubegraph"
		file: "/tmp/deploy.yaml"
		kubeConfig: kc
		fs:         _template.output.rootfs
	}
}

dagger.#Plan & {
	actions: {

		registry_url: client.env.REGISTRY_URL
		tag: client.env.IMAGE_TAG

		build: docker.#Dockerfile & {
			source: client.filesystem.".".read.contents
		}

		push: docker.#Push & {
			image: build.output
			dest:  "\(registry_url)/\(app):\(tag)"
			if client.env.REGISTRY_NEED_AUTH == "yes" {
				username: client.env.REGISTRY_USERNAME
				password: client.env.REGISTRY_PASSWORD
			}
		}

		deploy: #Deploy & {
            imageRef: push.result
            source: client.filesystem."./manifests".read.contents
            kc: client.commands.kc.stdout
        }

		test: {
			_prepare: #PrepareIntegrationTest & {
				srcTest: client.filesystem."./it-test/pytest".read.contents
				"registry_url": registry_url
			}

			integrationTest: #DoIntegrationTest & {
				imageRef: _prepare.output
				job:        client.filesystem."./it-test".read.contents
				kc:         client.commands.kc.stdout
			}
		}
	}
	client: {
		env: REGISTRY_PASSWORD: dagger.#Secret | ""
		env: REGISTRY_USERNAME: string | ""
		env: REGISTRY_URL: string | *"localhost:5001"
		env: IMAGE_TAG: string | *"ci"
		env: REGISTRY_NEED_AUTH: string | *"no"
		commands: kc: {
			name: "cat"
			args: ["kind-ci.yaml"]
			stdout: dagger.#Secret
		}
		filesystem: {
			".": read: contents: dagger.#FS

			"./it-test": read: {
				contents: dagger.#FS
				include: ["*.yaml"]
			}

			"./manifests": read: {
				contents: dagger.#FS
				include: ["*.yaml"]
			}

			"./it-test/pytest": read: {
				contents: dagger.#FS
				exclude: ["*cache*"]
			}
		}
	}
}
