pipeline{
    agent any
    environment {
            DOCKER_IMAGE = "cjburchell/restmock"
            DOCKER_TAG = "${env.BRANCH_NAME}"
            PROJECT_PATH = "/go/src/github.com/cjburchell/restmock"
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

        stage('Lint') {
                    steps {
                        script{
                        docker.withRegistry('https://390282485276.dkr.ecr.us-east-1.amazonaws.com', 'ecr:us-east-1:redpoint-ecr-credentials') {
                                docker.image('cjburchell/goci:latest').inside("-v ${env.WORKSPACE}:${PROJECT_PATH}"){
                                    sh """cd ${PROJECT_PATH} && go list ./... | grep -v /vendor/ > projectPaths"""
                                    def paths = sh returnStdout: true, script:"""awk '{printf "/go/src/%s ",\$0} END {print ""}' projectPaths"""

                                    def vetResults = sh returnStdout: true, script:"""go tool vet ${paths}"""
                                    writeFile file: 'vet_results.txt', text: vetResults
                                    echo vetResults

                                    def lintResults = sh returnStdout: true, script:"""golint ${paths}"""
                                    writeFile file: 'lint_results.txt', text: lintResults
                                    echo lintResults
                                }
                            }
                        }
                    }
                }

                /*stage('Tests') {
                    steps {
                        script{
                            docker.withRegistry('https://390282485276.dkr.ecr.us-east-1.amazonaws.com', 'ecr:us-east-1:redpoint-ecr-credentials') {
                                docker.image('cjburchell/goci:latest').inside("-v ${env.WORKSPACE}:${PROJECT_PATH}"){
                                    sh """cd ${PROJECT_PATH} && go list ./... | grep -v /vendor/ > projectPaths"""
                                    def paths = sh returnStdout: true, script:"""awk '{printf "/go/src/%s ",\$0} END {print ""}' projectPaths"""

                                    def testResults = sh returnStdout: true, script:"""go test -v ${paths}"""
                                    writeFile file: 'test_results.txt', text: testResults
                                    echo testResults
                                    sh """go2xunit -input test_results.txt > tests.xml"""
                                    sh """cd ${PROJECT_PATH} && ls"""

                                }
                            }
                        }
                    }
                }*/

        stage('Build') {
            steps {
                script {
                    if( env.BRANCH_NAME == "master")
                    {
                        docker.build("${DOCKER_IMAGE}").tag("latest")
                    }
                    else {
                        docker.build("${DOCKER_IMAGE}").tag("${DOCKER_TAG}")
                    }
                }
            }
        }

        stage ('Push') {
            steps {
                script {
                    docker.withRegistry('https://390282485276.dkr.ecr.us-east-1.amazonaws.com', 'ecr:us-east-1:redpoint-ecr-credentials') {
                            if( env.BRANCH_NAME == "master")
                            {
                                docker.image("${DOCKER_IMAGE}").push("latest")
                            }
                            else {
                                docker.image("${DOCKER_IMAGE}").push("${DOCKER_TAG}")
                            }
                        }
                    }
                }
        }
    }

    post {
                always {
                      archiveArtifacts '*results.txt'
                      //archiveArtifacts 'tests.xml'
                      //junit allowEmptyResults: true, testResults: 'tests.xml'
                      script{
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