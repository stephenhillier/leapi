pipeline {
  agent any
  stages {
    stage('Start pipeline') {
      steps {
        script {
          abortAllPreviousBuildInProgress(currentBuild)
        }
      }
    }
    stage('Create new build') {
      when {
        expression {
          script {
            openshift.withCluster() {
              return !openshift.selector("bc", "leapi-${env.JOB_BASE_NAME}".toLowerCase()).exists();
            }
          }
        }
      }
      steps {
        script {
          openshift.withCluster() {
            checkout scm
            openshift.newBuild("--docker-image=registry.access.redhat.com/devtools/go-toolset-7-rhel7:latest", ".", "--name=leapi-${env.JOB_BASE_NAME}".toLowerCase())
          }                   
        }
      }
    }
    stage('Test & build image') {
      steps {
        script {
          openshift.withCluster() {
            echo "${env.JOB_BASE_NAME}".toLowerCase()
            openshift.selector("bc", "leapi-${env.JOB_BASE_NAME}".toLowerCase()).startBuild("--wait")
          }
        }
      }
    }
    stage('Promote to dev') {
      steps {
        script {
          openshift.withCluster() {
            openshift.tag("leapi-${env.JOB_BASE_NAME}:latest", "leapi:dev")
          }
        }
      }
    }
    stage('Create new dev deployment') {
      when {
        expression {
          script {
            openshift.withCluster() {
              return !openshift.selector("dc", "leapi-dev-${env.JOB_BASE_NAME}".toLowerCase()).exists()
            }
          }
        }
      }
      steps {
        script {
          openshift.withCluster() {
            openshift.newApp("leapi:dev", "--name=leapi-dev-${env.JOB_BASE_NAME}".toLowerCase()).narrow("svc").expose("--port=8000")
          }
        }
      }
    }
  }
}

