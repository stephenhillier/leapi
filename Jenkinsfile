pipeline {
    agent any
    stages {
        stage('Create new application') {
            steps {
                script {
                    openshift.withCluster() {
                        openshift.newApp("go-toolset-7-rhel7:latest~https://github.com/stephenhillier/leapi.git")
                    }
                }
            }
        }
    }
}
