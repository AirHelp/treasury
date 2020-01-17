#!groovy
@Library('jenkins-pipeline-library@feature/golang-podTemplate') _

def label = "treasury-${UUID.randomUUID().toString()}"

golang(label: label) {
  node(label) {
    stage('github checkout') {
      checkout scm
    }

    ci {
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
}
