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
                        openshift.newBuild("go-toolset-7-rhel7:latest~https://github.com/stephenhillier/leapi.git")
                    }                   
                }
            }
        }
        stage('Build new image') {
            steps {
                script {
                    openshift.withCluster() {
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
                            !openshift.selector("dc", "leapi-dev")
                        }
                    }
                }
            }
            steps {
                script {
                    openshift.withCluster() {
                        openshift.newApp("leapi:dev", "--name=leapi-dev").narrow("svc").expose()
                    }
                }
            }
        }
    }
}

