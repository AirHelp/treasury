#!groovy

def label = "treasury-${UUID.randomUUID().toString()}"

podTemplate(label: label, containers: [
  containerTemplate(
    name: 'golang',
    image: 'golang:1.11-alpine',
    ttyEnabled: true,
    command: 'cat',
    resourceRequestCpu: '100m',
    resourceRequestMemory: '128Mi',
    envVars: [
        envVar(key: 'GO111MODULE', value: 'on'),
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

  }
}
