pipeline{
    agent any
    environment {
            DOCKER_IMAGE = "cjburchell/loki"
            DOCKER_TAG = "${env.BRANCH_NAME}-${env.BUILD_NUMBER}"
            PROJECT_PATH = "/go/src/github.com/cjburchell/loki"
    }

    stages{
        stage('Setup') {
            steps {
                script{
                    slackSend color: "good", message: "Job: ${env.JOB_NAME} with build number ${env.BUILD_NUMBER} started"
                }
             /* Let's make sure we have the repository cloned to our workspace */
             checkout scm
             }
        }
        stage('Static Analysis') {
            parallel {
                stage('Vet') {
                    agent {
                        docker {
                            image 'cjburchell/goci:1.14'
                            args '-v $WORKSPACE:$PROJECT_PATH'
                        }
                    }
                    steps {
                        script{
                                sh """go vet ./..."""

                                def checkVet = scanForIssues tool: [$class: 'GoVet']
                                publishIssues issues:[checkVet]
                        }
                    }
                }

                stage('Lint') {
                    agent {
                        docker {
                            image 'cjburchell/goci:1.14'
                            args '-v $WORKSPACE:$PROJECT_PATH'
                        }
                    }
                    steps {
                        script{
                            sh """golint ./..."""

                            def checkLint = scanForIssues tool: [$class: 'GoLint']
                            publishIssues issues:[checkLint]
                        }
                    }
                }
            }
        }

        stage('Build image') {
            steps {
                script {

                    def image = docker.build("${DOCKER_IMAGE}")
                    image.tag("${DOCKER_TAG}")
                    if( env.BRANCH_NAME == "master"){
                        image.tag("latest")
                    }
                }
            }
        }

        stage ('Push image') {
            steps {
                script {
                    docker.withRegistry('', 'dockerhub') {
                        def image = docker.image("${DOCKER_IMAGE}")
                        image.push("${DOCKER_TAG}")
                        if( env.BRANCH_NAME == "master") {
                            image.push("latest")
                        }
                    }
                }
            }
        }
    }

    post {
        always {
              script{
			      sh "docker system prune -f || true"
			      sh "docker image prune -af || true"
				  
                  if ( currentBuild.currentResult == "SUCCESS" ) {
                    slackSend color: "good", message: "Job: ${env.JOB_NAME} with build number ${env.BUILD_NUMBER} was successful"
                  }
                  else if( currentBuild.currentResult == "FAILURE" ) {
                    slackSend color: "danger", message: "Job: ${env.JOB_NAME} with build number ${env.BUILD_NUMBER} was failed"
                  }
                  else if( currentBuild.currentResult == "UNSTABLE" ) {
                    slackSend color: "warning", message: "Job: ${env.JOB_NAME} with build number ${env.BUILD_NUMBER} was unstable"
                  }
                  else {
                    slackSend color: "danger", message: "Job: ${env.JOB_NAME} with build number ${env.BUILD_NUMBER} its result (${currentBuild.currentResult}) was unclear"
                  }
              }
        }
    }

}