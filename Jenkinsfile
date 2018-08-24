pipeline {
  environment {
    // PR_NUM is the pull request number e.g. 'pr-4'
    PR_NUM = "${env.JOB_BASE_NAME}".toLowerCase()
  }
  agent any
  stages {
    stage('Start pipeline') {
      steps {
        script {
          abortAllPreviousBuildInProgress(currentBuild)
        }
      }
    }

    // create a new build config only if one does not already exist for this pull request.
    // todo: may be able to share build configs across pull requests
    stage('Create new build') {
      when {
        expression {
          script {
            openshift.withCluster() {
              return !openshift.selector("bc", "leapi-${PR_NUM}").exists();
            }
          }
        }
      }
      steps {
        script {
          openshift.withCluster() {
            checkout scm
            openshift.newBuild("--docker-image=registry.access.redhat.com/devtools/go-toolset-7-rhel7:latest", ".", "--name=leapi-${PR_NUM}")
          }                   
        }
      }
    }

    // start a new build for this pull request.
    // unit tests are run inside the s2i builder image (unit test dependencies
    // will not be available in the final application image).  See .s2i/bin/assemble
    // 
    // todo: output unit tests to jenkins output. Currently they will fail the pipeline run
    // but without any output.
    stage('Test & build image') {
      steps {
        script {
          openshift.withCluster() {
            echo "Starting build from pull request ${PR_NUM}"
            openshift.selector("bc", "leapi-${PR_NUM}").startBuild("--wait")
          }
        }
      }
    }

    // upon successful build, tag the image `dev`.
    // the dev deployment will automatically run as soon as a new dev image is ready.
    stage('Promote to dev') {
      steps {
        script {
          openshift.withCluster() {
            openshift.tag("leapi-${PR_NUM}:latest", "leapi-${PR_NUM}:dev")
          }
        }
      }
    }

    // create a new deployment for dev.
    // this stage is skipped if a deployment config already exists for this pull request
    // otherwise, a new app will be created using the built image
    stage('Create new dev deployment') {
      when {
        expression {
          script {
            openshift.withCluster() {
              return !openshift.selector("dc", "leapi-dev-${PR_NUM}").exists()
            }
          }
        }
      }
      steps {
        script {
          openshift.withCluster() {
            openshift.newApp("leapi-${PR_NUM}:dev", "--name=leapi-dev-${PR_NUM}").narrow("svc").expose("--port=8000")
          }
        }
      }
    }
    stage('Verify deployment') {
      steps {
        script {
          openshift.withCluster() {
            echo "Waiting for deployment to dev..."
            def dc = openshift.selector("dc", "leapi-dev-${PR_NUM}")

            // wait until each container in this deployment's pod reports as ready
            dc.related("pods").untilEach(1) {
              return it.object().status.containerStatuses.every {
                it.ready
              }
            }
            echo "Deployment successful!"
          }
        }
      }
    }
  }
}

