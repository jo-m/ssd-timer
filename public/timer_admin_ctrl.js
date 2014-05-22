angular.module('app').controller('TimerAdminCtrl', function($scope, $timeout) {
    $scope.retryIntervalMs = 2000
    $scope.Config = {}

    $scope.sendCmd = function(cmd) {
        $scope.conn.send(cmd)
    };

    $scope.sendConfig = function() {
        $scope.conn.send(JSON.stringify($scope.Config));
    };

    $scope.openConnection = function() {
        $scope.status = "Retrying connection..."
        if (window["WebSocket"]) {
            $scope.conn = new WebSocket("ws://" + ws_hostname + "/adminws");
            $scope.conn.onclose = function(evt) {
                $scope.status = "No connection"
                $timeout($scope.openConnection, $scope.retryIntervalMs)
            }
            $scope.conn.onmessage = function(evt) {
                $scope.Config = JSON.parse(evt.data).Config;
                $scope.status = "Connection ok"
            }
        } else {
            $scope.status = "Your client does not have web socket support"
        }
    }

    $scope.openConnection()
})
