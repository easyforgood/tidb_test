#!groovy

node {
    def TIDB_TEST_BRANCH = "master"
    // def TIKV_BRANCH = "master"
    // def PD_BRANCH = "master"

    environment {
        PATH = "$PATH:/usr/local/bin"
    }

    checkout scm
    dir("SRE"){
        git url: 'https://github.com/easyforgood/SRE_test.git'
        if (!fileExists("bin") ){
            sh "mkdir bin"
        }
    }
    
    stage("build"){
        try{
            stage("build tidb"){
                dir("tidb"){
                    checkout scm
                    sh "make server"
                    sh "cp bin/*  ../SRE/bin/"
                }
            }
            stage("build tikv"){
                dir("tikv"){
                git url: 'https://github.com/pingcap/tikv'
                    sh "make build"
                    sh "cp target/debug/tikv-server  ../SRE/bin/"
                }
            }

            stage("build pd"){
                dir("pd"){
                    git url: 'https://github.com/pingcap/pd'
                    sh "make build"
                    sh "cp bin/*  ../SRE/bin/"
                }
            }
        }
        catch(e){
            emailext body: "项目构建失败.\r\n 构建地址:${BUILD_URL}\r\n Pull-Request:${GITHUB_PR_URL} \r\n 异常消息：${e.toString()}", subject: 'Pull-Request 构建结果通知【失败】', to: "${GITHUB_PR_AUTHOR_EMAIL}"
            sh "exit 1"
        }
    }

    stage("test"){
         try{
            dir("tidb"){
               sh "make test"
            }
            dir("pd"){
               // sh "make test"
            }
            dir("tikv"){
               sh "make test"
            }
         }catch(e){
            echo "单元测试失败"
            emailext body: "单元测试失败.\r\n 构建地址:${BUILD_URL}\r\n Pull-Request:${GITHUB_PR_URL} \r\n错误消息：${e.toString()}", subject: 'Pull-Request 构建结果通知【失败】', to: "${GITHUB_PR_AUTHOR_EMAIL}"
               sh "exit 1"
        }
    }
    stage("create docker images"){
        dir("SRE"){
            try{
                docker.build("tidb_test", "-f tidb/Dockerfile .")
                docker.build("tikv_test", "-f tikv/Dockerfile .")
                docker.build("pd_test", "-f pd/Dockerfile .")
            }catch(e){
                emailext body: "创建测试容器镜像失败.\r\n 构建地址:${BUILD_URL}\r\n Pull-Request:${GITHUB_PR_URL} \r\n 异常消息：${e.toString()}", subject: 'Pull-Request 构建结果通知【失败】', to: "${GITHUB_PR_AUTHOR_EMAIL}"
                sh "exit 1"
            }
        }
    }
    stage("integration test"){
        dir("SRE"){
            sh "docker-compose up -d tidb"
            // or http health api?
            try{
                sleep 40
                    def retCode, retMsg = sh (
                        script: "go run integration/main.go",
                        returnStatus: true,
                        returnStdout: true
                    ).trim()
                    echo integration_test_result
            }catch(e){
                emailext body: "集成测试失败.\r\n 构建地址:${BUILD_URL}\r\n Pull-Request:${GITHUB_PR_URL} \r\n 异常消息：${e.toString()}", subject: 'Pull-Request 构建结果通知【失败】', to: "${GITHUB_PR_AUTHOR_EMAIL}"
                sh "exit 1"
            }finally{
                sh "docker-compose down"
            }
        }
    }
    emailext body: "构建成功.\r\n 构建地址:${BUILD_URL}\r\n Pull-Request:${GITHUB_PR_URL} \r\n ", subject: 'Pull-Request 构建结果通知【成功】', to: "${GITHUB_PR_AUTHOR_EMAIL}"
}
