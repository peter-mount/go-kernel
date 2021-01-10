// Build properties
properties([
  buildDiscarder(logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '', numToKeepStr: '10')),
  disableConcurrentBuilds(),
  disableResume(),
  pipelineTriggers([
    cron('H H * * *')
  ])
])

node( 'Build' ) {
  stage( 'prepare' ) {
    checkout scm
    sh 'docker build -t golib:build --target build .'
  }

  stage( "build" ) {
    sh 'docker build -t golib:test .'
  }
}
