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
                        openshift.newApp("go-toolset-7-rhel7:latest~https://github.com/stephenhillier/leapi.git")
                    }                   
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

