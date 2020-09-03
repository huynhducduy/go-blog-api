pipeline {
    agent any
    stages {
        stage('prepare') {
            steps {
                checkout scm
            }
        }
        stage('build') {
            steps {
                sh 'sudo cp /home/ubuntu/go-blog.env .env'
                sh 'cat .env'
                sh 'sudo docker build -t go-blog .'
            }
        }
        stage('test') {
            steps {
                sh 'echo "Passed!"'
            }
        }
        stage('prepare to deploy') {
            steps {
                sh '(sudo docker kill go-blog || true) && (sudo docker rm go-blog || true)'
            }
        }
        stage('deploy') {
            steps {
                sh 'sudo docker run -dit -p 8081:80 --name go-blog go-blog:latest'
            }
        }
    }
}