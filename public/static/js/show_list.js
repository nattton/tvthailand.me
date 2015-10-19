angular.module('infinite-scroll').value('THROTTLE_MILLISECONDS', 250)
app.filter('urlencode', function() {
  return function(input) {
    return window.encodeURIComponent(input);
  }
});
app.controller('VideoCtrl', function ($scope, $http) {
  $scope.shows = [];
  $scope.busy = false;
  $scope.last_length = -1;
  $scope.nextPage = function() {
    if ($scope.busy || $scope.last_length >= $scope.shows.length) return;
      $scope.busy = true;
    if ($scope.apiPath == "") {
      $scope.busy = true;
      return;
    }
    $scope.last_length = $scope.shows.length
    $http.get("/api/v1" + $scope.apiPath + $scope.shows.length).
      error(function(){

      }).
      success(function(data){
          angular.forEach(data.Shows, function(show) {
            $scope.shows.push(show)
          });
          $scope.busy = false;
      });
  };
});
