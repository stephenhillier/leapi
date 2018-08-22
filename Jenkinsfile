pipeline {
    agent any
    stages {
        stage('Create new build') {
            steps {
                when {
                    expression {
                        openshift.withCluster() {
                            return !openshift.selector("bc", "leapi").exists();
                        }
                    }
                }
                openshift.withCluster() {
                    openshift.newApp("go-toolset-7-rhel7:latest~https://github.com/stephenhillier/leapi.git")
                }
            }
        }
        stage('Build new image') {
            steps {
                script {
                    openshift.withCluster() {
                        openshift.selector("bc", "leapi").newBuild()
                    }
                }
            }
        }
    }
}

