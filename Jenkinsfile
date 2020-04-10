def imageName = 'openapi-linter'
def currentBranch = env.BRANCH_NAME
def projectName = ('jenkins-' + currentBranch).replaceAll('\'', '_')

node('slave') {
    try {
        stage('Checkout') {
            checkout scm
        }
        ansiColor('xterm') {
            stage('linters') {
                sh "make lint"
            }

            stage('test') {
                sh "make test"
            }


            node('master') {
                stage('Docs') {
                    checkout scm
                    def tag = versions.tag()
                    // tagged release should already have changelog generated
                    if ( !(currentBranch == 'master' && tag) ) {
                        sh "make changelog"
                    }
                    sh "make documentation"
                    hostDocs('openapi-linter', 'build/html')
                }
            }
        }
    } finally {
        stage('cleanup') {
            sh "echo 'Good bye'"
        }
    }
}
