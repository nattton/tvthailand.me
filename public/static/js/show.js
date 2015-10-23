angular.module('infinite-scroll').value('THROTTLE_MILLISECONDS', 250)
app.filter('urlencode', function() {
  return function(input) {
    return window.encodeURIComponent(input);
  }
});
app.controller('VideoCtrl', function ($scope, $http) {
  $scope.episodes = [];
  $scope.busy = false;
  $scope.last_length = -1;
  $scope.nextPage = function() {
    if ($scope.busy || $scope.last_length >= $scope.episodes.length) return;
      $scope.busy = true;
    $scope.last_length = $scope.episodes.length;
    $http.get("/api/v1/show/" + $scope.showID + "/" + $scope.episodes.length).
      error(function(){

      }).
      success(function(data){
          angular.forEach(data.Episodes, function(episode) {
            $scope.episodes.push(episode)
          });
          $scope.busy = false;
      });
  };
});
