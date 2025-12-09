pipeline {
    agent any
    
    environment {
        REPO = 'https://github.com/yakovlevruslan/kbot'
        BRANCH = 'main'
        GITHUB_TOKEN = credentials('github-token')
    }
    parameters {
  choice choices: ['linux', 'darwin', 'windows'], description: 'Target operating system', name: 'OS'
  choice choices: ['amd64', 'arm64'], description: 'Target architecture', name: 'ARCH'
}    
    stages {

        stage('clone') {
            steps {
                echo 'Clone Repository'
                git branch: "${BRANCH}", url: "${REPO}"
            }
        }

        stage('test') {
            steps {
                echo 'Testing started'
                sh "make test"
            }
        }

        stage('build') {
            steps {
                echo "Building binary started"
                sh "make build TARGETOS=${params.OS} TARGETARCH=${params.ARCH}"
            }
        }

        stage('image') {
            steps {
                echo "Building image started"
                sh "make image TARGETOS=${params.OS} TARGETARCH=${params.ARCH}"
            }
        }        

        stage('login to GHCR') {
            steps {
                sh "echo $GITHUB_TOKEN_PSW | docker login ghcr.io -u $GITHUB_TOKEN_USR --password-stdin"
            }
        }
        
        stage('push image') {
            steps {
              sh "make push TARGETOS=${params.OS} TARGETARCH=${params.ARCH}"
            }
        } 
    }
}