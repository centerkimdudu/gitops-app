pipeline {
    agent {
        kubernetes {
            yaml '''
            apiVersion: v1
            kind: Pod
            spec:
              containers:
              - name: kaniko
                image: gcr.io/kaniko-project/executor:debug
                command:
                - sleep
                args:
                - "9999999"
                volumeMounts:
                - name: docker-config
                  mountPath: /kaniko/.docker
              volumes:
              - name: docker-config
                secret:
                  secretName: docker-credentials
                  items:
                  - key: .dockerconfigjson
                    path: config.json
            '''
        }
    }

    environment {
        DOCKER_IMAGE = 'kimtaejung/gitops-demo-app'
        IMAGE_TAG = "${BUILD_NUMBER}"
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: 'main',
                    url: 'https://github.com/centerkimdudu/gitops-app.git',
                    credentialsId: 'github-credentials'
            }
        }

        stage('Build & Push Image') {
            steps {
                container('kaniko') {
                    sh """
                        /kaniko/executor \
                          --context=dir://\${WORKSPACE} \
                          --destination=\${DOCKER_IMAGE}:\${IMAGE_TAG} \
                          --destination=\${DOCKER_IMAGE}:latest
                    """
                }
            }
        }
    }

    post {
        success { echo "✅ Build #\${BUILD_NUMBER} succeeded" }
        failure { echo "❌ Build #\${BUILD_NUMBER} failed" }
    }
}
