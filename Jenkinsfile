#!groovy

def label = "treasury-${UUID.randomUUID().toString()}"

podTemplate(label: label, containers: [
  containerTemplate(
    name: 'golang',
    image: 'golang:1.11-alpine',
    ttyEnabled: true,
    command: 'cat',
    resourceRequestCpu: '100m',
    resourceRequestMemory: '128Mi'
  )
  ]) {
  node(label) {
    stage('go test') {
      container('golang'){
        withEnv('GO111MODULE=on'){
          sh 'apk add --no-cache git'
          sh 'go mod download'
          sh 'go test -cover -v ./...'
        }
      }
    }
  }
}
