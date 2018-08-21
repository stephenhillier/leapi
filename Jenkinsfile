pipeline {
    agent any
    stages {
        stage('Pipeline test') {
            steps {
                script {
                    openshift.withCluster() {
                        openshift.withProject() {
                            echo "Hello from project ${openshift.project()} in cluster ${openshift.cluster()}"
                        }
                    }
                }
            }
        }
        stage('Create new application') {
            steps {
                script {
                    openshift.withCluster() {
                        openshift.newApp("registry.access.redhat.com/devtools/go-toolset-7-rhel7:latest~https://github.com/stephenhillier/leapi.git")
                    }
                }
            }
        }
    }
}
