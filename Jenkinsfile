pipeline {
  def buildVersion = "${env.JOB_BASE_NAME}-${env.BUILD_NUMBER}-${env.CHANGE_ID}"
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
              return !openshift.selector("bc", "leapi-${buildVersion}").exists();
            }
          }
        }
      }
      steps {
        script {
          openshift.withCluster() {
            checkout scm
            openshift.newBuild("--from-image=registry.access.redhat.com/devtools/go-toolset-7-rhel7:latest", ".", "--name=leapi-${buildVersion}")
          }                   
        }
      }
    }
    stage('Test & build image') {
      steps {
        script {
          openshift.withCluster() {
            echo "${buildVersion}"
            openshift.selector("bc", "leapi-${buildVersion}").startBuild("--wait")
          }
        }
      }
    }
    stage('Promote to dev') {
      steps {
        script {
          openshift.withCluster() {
            openshift.tag("leapi-${buildVersion}:latest", "leapi:dev")
          }
        }
      }
    }
    stage('Create new dev deployment') {
      when {
        expression {
          script {
            openshift.withCluster() {
              return !openshift.selector("dc", "leapi-dev-${buildVersion}").exists()
            }
          }
        }
      }
      steps {
        script {
          openshift.withCluster() {
            openshift.newApp("leapi:dev", "--name=leapi-dev-${buildVersion}").narrow("svc").expose("--port=8000")
          }
        }
      }
    }
  }
}

