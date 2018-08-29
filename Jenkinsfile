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
    // 
    // todo: output unit tests to jenkins output. Currently they will fail the pipeline run
    // but without any output.
    stage('Test & build image') {
      steps {
        script {
          openshift.withCluster() {
            // create a new build config if one does not already exist
            // unit tests are run during build (application will not be built if unit tests fail)
            if ( !openshift.selector("bc", "leapi-${PR_NUM}-builder").exists() ) {
              echo "Creating a new build config for pull request ${PR_NUM}"
              openshift.newBuild("https://github.com/stephenhillier/leapi.git#pull/${env.CHANGE_ID}/head", "--strategy=docker", "--name=leapi-${PR_NUM}-builder", "-l pr=${PR_NUM}")
            } else {
              echo "Starting build from pull request ${PR_NUM}"
              openshift.selector("bc", "leapi-${PR_NUM}-builder").startBuild("--wait")
            }

            echo "Waiting for builds from buildconfig leapi-${PR_NUM}-builder to finish"
            def lastBuildNumber = openshift.selector("bc", "leapi-${PR_NUM}-builder").object().status.lastVersion
            def lastBuild = openshift.selector("build", "leapi-${PR_NUM}-builder-${lastBuildNumber}")
            timeout(10) {
              lastBuild.untilEach(1) {
                return it.object().status.phase == "Complete"
              }
            }

            // start building an application image.
            // this is a chained build; only the application binary will be brought forward from the builder image.
            // the image can only be used as an executable
            if ( !openshift.selector("bc", "leapi-${PR_NUM}").exists() ) {
              echo "Creating new application build config"
              openshift.newBuild("alpine:3.8", "--source-image=leapi-${PR_NUM}-builder", "--allow-missing-imagestream-tags", "--name=leapi-${PR_NUM}", "-l pr=${PR_NUM}", "--source-image-path=/go/bin/leapi:.", """--dockerfile='FROM alpine:3.8
              RUN mkdir -p /app
              COPY leapi /app/leapi
              ENTRYPOINT [\"/app/leapi\"]
              '""")
            } else {
              echo "Creating application image"
              openshift.selector("bc", "leapi-${PR_NUM}").startBuild("--wait")
            }

            echo "Waiting for application build from leapi-${PR_NUM} to finish"
            def lastAppBuildNumber = openshift.selector("bc", "leapi-${PR_NUM}").object().status.lastVersion
            def lastAppBuild = openshift.selector("build", "leapi-${PR_NUM}-${lastAppBuildNumber}")
            timeout(10) {
              lastAppBuild.untilEach(1) {
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
              openshift.newApp(
                "leapi-${PR_NUM}:dev openshift/postgresql:9.6",
                "-l pr=${PR_NUM}",
                "-e POSTGRESQL_DATABASE=testdb"
              )

              openshift.selector("dc", "leapi-dev-${PR_NUM}").expose("--port=8000")
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

