angular.module('app', [])

.filter('toMinSec', function() {
  return function(input) {
    if(input == undefined) {
        input = 0
    }
    var sec_num = Math.floor(parseInt(input, 10) / 1000)
    var hours   = Math.floor(sec_num / 3600)
    var minutes = Math.floor((sec_num - (hours * 3600)) / 60)
    var seconds = sec_num - (hours * 3600) - (minutes * 60)

    if (hours   < 10) {hours   = "0"+hours}
    if (minutes < 10) {minutes = "0"+minutes}
    if (seconds < 10) {seconds = "0"+seconds}
    var time    = minutes+':'+seconds
    return time
  };
})

.controller('TimerCtrl', function ($scope, $interval, $timeout) {
    $scope.timerIntervalMs = 200
    $scope.retryIntervalMs = 2000
    $scope.timer = {Config: {}, State: {}}

    $scope.updateTimer = function() {
        if($scope.timer.State.TimerStartAt == undefined) {
            return;
        }

        $scope.calculatedState = {}

        // calculate when the timer will end
        $scope.calculatedState.timerEndAt = new Date($scope.timer.State.TimerStartAt.getTime())
        $scope.calculatedState.timerEndAt.setTime(
            $scope.calculatedState.timerEndAt.getTime() +
            $scope.timer.Config.RoundLength / 1000000 * $scope.timer.Config.NumRounds +
            $scope.timer.Config.PauseLength / 1000000 * ($scope.timer.Config.NumRounds - 1)
        )

        // virtual "current time", might be different when timer is paused
        $scope.calculatedState.currentTimerTime = new Date()

        if($scope.timer.State.TimerPaused) {
            $scope.statusText = "Paused"
            $scope.calculatedState.currentTimerTime = $scope.timer.State.TimerPausedAt

            // move forward the time end time
            $scope.calculatedState.timerEndAt.setTime(
                $scope.calculatedState.timerEndAt.getTime() +
                (new Date()).getTime() -
                $scope.calculatedState.currentTimerTime
            )
        }

        $scope.calculatedState.timeElapsed =
            $scope.calculatedState.currentTimerTime -
            $scope.timer.State.TimerStartAt

        if($scope.calculatedState.currentTimerTime < $scope.timer.State.TimerStartAt) {
            // not started yet
            $scope.calculatedState.round = 1
            $scope.calculatedState.timerTitle = $scope.timer.Config.NotStartedText
            $scope.calculatedState.delta = $scope.timer.Config.RoundLength / 1000000

        } else if($scope.calculatedState.currentTimerTime > $scope.calculatedState.timerEndAt){
            // already finished
            $scope.calculatedState.round = $scope.timer.Config.NumRounds
            $scope.calculatedState.timerTitle = $scope.timer.Config.FinishedText
            $scope.calculatedState.delta = 0

        } else {
            // running or pause

            // calculate round number
            $scope.calculatedState.round = 1 + Math.floor(
                $scope.calculatedState.timeElapsed /
                (($scope.timer.Config.RoundLength + $scope.timer.Config.PauseLength) / 1000000)
            )

            // calculate when the current round ends
            $scope.calculatedState.currentRoundEnds = new Date($scope.timer.State.TimerStartAt.getTime())
            $scope.calculatedState.currentRoundEnds.setTime(
                $scope.calculatedState.currentRoundEnds.getTime() +
                $scope.timer.Config.RoundLength / 1000000 * $scope.calculatedState.round +
                $scope.timer.Config.PauseLength / 1000000 * ($scope.calculatedState.round)
            )

            $scope.calculatedState.delta =
                $scope.calculatedState.currentRoundEnds -
                $scope.calculatedState.currentTimerTime

            // are we in the pause?
            if($scope.calculatedState.delta < $scope.timer.Config.PauseLength / 1000000) {
                $scope.calculatedState.pause = true
                // $scope.calculatedState.delta = $scope.timer.Config.RoundLength / 1000000
            } else {
                $scope.calculatedState.pause = false
                $scope.calculatedState.delta -= $scope.timer.Config.PauseLength / 1000000
            }

            $scope.calculatedState.timerTitle = $scope.calculatedState.pause ? $scope.timer.Config.PauseText : $scope.timer.Config.RoundText
        }
    }

    $interval($scope.updateTimer, $scope.timerIntervalMs);

    $scope.openConnection = function() {
        if (window["WebSocket"]) {
            $scope.conn = new WebSocket("ws://" + ws_hostname + "/timer");
            $scope.conn.onclose = function(evt) {
                $scope.timer.State.TimerPaused = true
                $scope.statusText = "No connection"
                $timeout($scope.openConnection, $scope.retryIntervalMs)
            }
            $scope.conn.onmessage = function(evt) {
                $scope.timer = JSON.parse(evt.data)
                // parse dates
                $scope.timer.State.TimerStartAt = new Date($scope.timer.State.TimerStartAt)
                $scope.timer.State.TimerPausedAt = new Date($scope.timer.State.TimerPausedAt)
                $scope.updateTimer()
            }
        } else {
            $scope.timer.State.TimerPaused = true
            $scope.statusText = "Your client does not have web socket support"
        }
    }

    $scope.openConnection()
})
