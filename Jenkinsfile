pipeline {
    agent {
        label 'go-builder'
    }
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
    }
}
