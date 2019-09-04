#!groovy

properties([disableConcurrentBuilds()])

pipeline {
  agent any

  triggers {
    pollSCM('*/30 * * * *')
  }
  options {
    buildDiscarder(logRotator(numToKeepStr: '10', artifactNumToKeepStr: '10'))
    timestamps()
  }
  stages {
    stage("docker login") {
      steps {
        echo " ============== docker login =================="
        withCredentials([usernamePassword(credentialsId: 'dockerhub_keep2share', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD')]) {
          sh """
            docker login -u $USERNAME -p $PASSWORD
          """
        }
      }
    }
    stage("Create docker image") {
      steps {
        dir ('.') {
          sh "docker build -f Dockerfile -t keep2share/${env.JOB_NAME}:latest ."
          sh "docker push keep2share/${env.JOB_NAME}:latest"
        }
      }
    }
  }
  post {
    success {
      slackSend(
        color: 'good',
        message: "Successfully build: ${env.JOB_NAME} #${env.BUILD_NUMBER} :tada:",
      )
    }
    failure {
      slackSend(
        color: 'danger',
        message: "Build FAILED: ${env.JOB_NAME} #${env.BUILD_NUMBER} :facepalm_mem:",
      )
    }
  }
}