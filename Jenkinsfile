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
  stage( 'checkout' ) {
    checkout scm
  }

  stage( 'go mod' ) {
    // Seems go 1.16+ have changed things so this is required otherwise modules are not handled correctly with go.sum breaking
    // via https://github.com/golang/go/issues/44129#issuecomment-860060061
    sh 'go env -w GOFLAGS=-mod=mod'

    sh 'go mod download'
  }

// Run each package separately so they don't interfere with each other.
// Also some like moduletest MUST be separate as it can work only once.
// Hence we have a stage for each compatible set
  stage( "test util" ) {
    sh 'CGO_ENABLED=0 go test -v ./util'
    sh 'CGO_ENABLED=0 go test -v ./util/walk'
  }

  stage( "test" ) {
    sh 'CGO_ENABLED=0 go test -v .'
    sh 'CGO_ENABLED=0 go test -v ./test'
  }

  stage( "test modules" ) {
    sh 'CGO_ENABLED=0 go test -v ./test/moduletest'
  }

  stage( "test interfaces" ) {
    sh 'CGO_ENABLED=0 go test -v ./test/interfaces'
  }
}
