pipeline {
  agent any
  stages {
    stage('unit test') {
      steps {
        sh 'cd /root/.jenkins/workspace/cddemo_master && make test'
      }
    }
    stage('build and push') {
      steps {
        sh 'cd /root/.jenkins/workspace/cddemo_master && make docker'
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