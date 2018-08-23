pipeline {
  agent any
  stages {
    stage('Create new build') {
      when {
        expression {
          script {
            openshift.withCluster() {
              return !openshift.selector("bc", "leapi").exists();
            }
          }
        }
      }
      steps {
        script {
          openshift.withCluster() {
            openshift.newBuild("registry.access.redhat.com/devtools/go-toolset-7-rhel7:latest~https://github.com/stephenhillier/leapi.git")
          }                   
        }
      }
    }
    stage('Test & Build') {
      steps {
        script {
          openshift.withCluster() {
            echo "${env.JOB_BASE_NAME}-${env.BUILD_NUMBER}-${env.CHANGE_ID}"
            openshift.selector("bc", "leapi").startBuild("--wait")
          }
        }
      }
    }
    stage('Promote to dev') {
      steps {
        script {
          openshift.withCluster() {
            openshift.tag("leapi:latest", "leapi:dev")
          }
        }
      }
    }
    stage('Create dev') {
      when {
        expression {
          script {
            openshift.withCluster() {
              return !openshift.selector("dc", "leapi-dev").exists()
            }
          }
        }
      }
      steps {
        script {
          openshift.withCluster() {
            openshift.newApp("leapi:dev", "--name=leapi-dev").narrow("svc").expose("--port=8000")
          }
        }
      }
    }
  }
}

