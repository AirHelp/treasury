#!groovy

def label = "treasury-${UUID.randomUUID().toString()}"

podTemplate(label: label, containers: [
  containerTemplate(
    name: 'golang',
    image: 'golang:1.12.6-alpine',
    ttyEnabled: true,
    command: 'cat',
    resourceRequestCpu: '100m',
    resourceRequestMemory: '128Mi',
    envVars: [
        envVar(key: 'GO111MODULE', value: 'on'),
        envVar(key: 'CGO_ENABLED', value: '0'),
    ]
  )
  ]) {
  node(label) {
    stage('github checkout') {
      checkout scm
    }

    stage('download Go deps') {
      container('golang'){
        sh 'apk add --no-cache git'
        sh 'go mod download'
      }
    }

    stage('go test') {
      container('golang'){
        sh 'go test -cover -v ./...'
      }
    }

    stage('go formatting') {
      container('golang'){
        sh 'gofmt -s -w .'
      }
    }

    stage('go vet') {
      container('golang'){
        sh 'go vet -v ./...'
      }
    }

  }
}
