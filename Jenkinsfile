#!groovy
@Library('jenkins-pipeline-library') _

def label = "treasury-${UUID.randomUUID().toString()}"

golang(label: label) {
  node(label) {
    stage('github checkout') {
      checkout scm
    }

    ci {
      golangTestPipeline()
    }
  }
}
