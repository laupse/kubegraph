package main

import (
	"dagger.io/dagger"
	"dagger.io/dagger/core"
	"universe.dagger.io/docker"
	"universe.dagger.io/alpha/kubernetes/kapp"
	"universe.dagger.io/alpha/kubernetes/kustomize"
)

registry_url: string | *"localhost:5001"
app: string | *"kubegraph"
tag: string | *"ci"

#PrepareIntegrationTest: {   
    srcTest: dagger.#FS

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
        dest: "\(registry_url)/it-test:ci"
    }

    done: bool | *true
}

#DoIntegrationTest: {
    job: dagger.#FS
    kc: dagger.#Secret
    imageReady: bool

    deployment: kapp.#Deploy & {
        app:        "it-test"
        fs:         job
        kubeConfig: kc
        file:       "./it-test-job.yaml"
    }
}

#Deploy: {
    imageRef: string

    manifest: dagger.#FS

    kustomizationFs: dagger.#FS

    outputFile: string | *"result.yaml"

    kubeConfig: dagger.#Secret

    _template: kustomize.#Kustomize & {
        source:        manifest
        kustomization: kustomizationFs
    }

    _file: core.#WriteFile & {
        input:    dagger.#Scratch
        path:     outputFile
        contents: _template.output
    }

    deploy: kapp.#Deploy & {
        app:          "kubegraph"
        fs:           _file.output
        "kubeConfig": kubeConfig
        file:         outputFile
    }

    output: _file.output
}

dagger.#Plan & {
	actions: {

		build: docker.#Dockerfile & {
			source: client.filesystem.".".read.contents
		}

        push: docker.#Push & {
            image: build.output
            dest: "\(registry_url)/\(app):\(tag)"
        }
         
        deployInt: #Deploy & {
            imageRef: push.result
            manifest: client.filesystem."./manifests".read.contents
            kustomizationFs: client.filesystem."./it-test".read.contents
            kubeConfig: client.commands.kc.stdout
        }

        test: {    
            _prepare: #PrepareIntegrationTest & {
                srcTest: client.filesystem."./it-test/pytest".read.contents
            }   

            integrationTest: #DoIntegrationTest & {
                imageReady: _prepare.done
                job: client.filesystem."./it-test".read.contents
                kc: client.commands.kc.stdout
                
            }
        }
	}
	client: {
        commands: kc: {
            name: "cat"
            args: ["kind-ci.yaml"]
            stdout: dagger.#Secret
        }
        filesystem: {
            ".": read: {
                contents: dagger.#FS
            }

            "./it-test": read: {
                contents: dagger.#FS
                include: ["*.yaml"]
            }

            "./it-test/pytest": read: {
                contents: dagger.#FS
                exclude: ["*cache*"]
            }

            "./manifests": read: {
                contents: dagger.#FS
            }

            "./_tmp": write: {
                contents: actions.deployInt.output
            }
        }
    }
}