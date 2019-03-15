pipeline {
  agent any
  stages {
    stage('build docker image') {
      steps {
        sh 'echo "build docker image"'
      }
    }
    stage('push image') {
      steps {
        sh 'echo "push image to docker hub"'
      }
    }
    stage('deploy app') {
      steps {
        sh 'echo "deploy app"'
      }
    }
    stage('test app') {
      steps {
        sh 'echo "test app"'
      }
    }
  }
}