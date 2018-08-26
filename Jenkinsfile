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
            // create a new build config if one does not already exist
            if ( !openshift.selector("bc", "leapi-${PR_NUM}-builder").exists() ) {
              echo "Creating a new build config for pull request ${PR_NUM}"
              openshift.newBuild("https://github.com/stephenhillier/leapi.git#pull/${env.CHANGE_ID}/head", "--strategy=docker", "--name=leapi-${PR_NUM}-builder")
            }

            if ( !openshift.selector("bc", "leapi-${PR_NUM}").exists() ) {
              openshift.newBuild("alpine:3.8", "--source-image=leapi-${PR_NUM}-builder", "--name=leapi-${PR_NUM}", "--source-image-path=/go/bin/leapi:.", """--dockerfile='FROM alpine:3.8
              RUN mkdir -p /app
              COPY leapi /app/leapi
              ENTRYPOINT [\"/app/leapi\"]
              '""")
            } else {
              echo "Starting build from pull request ${PR_NUM}"
              openshift.selector("bc", "leapi-${PR_NUM}").startBuild("--wait")
            }

            def builds = openshift.selector("bc", "leapi-${PR_NUM}").related("builds")

            timeout(10) {
              builds.untilEach {
                return it.object().status.phase == "Complete"
              }
            }

            // the dev deployment will automatically run as soon as a new image is tagged as `dev`
            echo "Successfully built image: tagging as new dev image"
            openshift.tag("leapi-${PR_NUM}:latest", "leapi-${PR_NUM}:dev")

          }
        }
      }
    }

    // Deployment to dev happens automatically when a new image is tagged `dev`.
    // This stage monitors the newest deployment for pods/containers to report back as ready.
    stage('Deploy to dev') {
      steps {
        script {
          openshift.withCluster() {

            // if a deployment config does not exist for this pull request, create one
            if ( !openshift.selector("dc", "leapi-dev-${PR_NUM}").exists() ) {
              echo "Creating a new deployment config for pull request ${PR_NUM}"
              openshift.newApp("leapi-${PR_NUM}:dev", "--name=leapi-dev-${PR_NUM}").narrow("dc").expose("--port=8000")
            }

            echo "Waiting for deployment to dev..."
            def newVersion = openshift.selector("dc", "leapi-dev-${PR_NUM}").object().status.latestVersion

            // find the pods for the newest deployment
            def pods = openshift.selector('pod', [deployment: "leapi-dev-${PR_NUM}-${newVersion}"])

            // wait until each container in this deployment's pod reports as ready
            timeout(10) {
              pods.untilEach(1) {
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
}

