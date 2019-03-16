#!groovy

node {
    def TIDB_TEST_BRANCH = "master"
    def TIKV_BRANCH = "master"
    def PD_BRANCH = "master"

    checkout scm
    stage("build"){
        stage("build tidb"){
            dir("tidb"){
                checkout scm
                sh "make"
            }
        }

        stage("build tikv"){
            dir("tikv"){
                git url: 'https://github.com/pingcap/tikv'
                sh "make"
            }
        }

        stage("build tikv"){
            dir("pd"){
                git url: 'https://github.com/pingcap/pd'
                sh "make"
            }
        }
    }

    stage("test"){
            dir("tidb"){
                sh "make test"
            }

            dir("pd"){
                sh "make test"
            }

            dir("tikv"){
                sh "make test"
            }
    }
    stage("create docker images"){

    }
    stage("integration test"){
    }
}
