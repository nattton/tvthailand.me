app.controller('VideoCtrl', function ($scope, $http, $sce) {
  $scope.playYoutube = function(playIndex) {
    var playerInstance = jwplayer("player");
    playerInstance.setup({
      autostart: true,
       playlist: $scope.episode.Playlists
    });
    playerInstance.on('ready', function(){
      playerInstance.playlistItem(playIndex);
    });
  }
  $scope.watch = function(id, playIndex) {
    $http.get('/api/v1/watch/' + id, { params: { device: "web" }}).
    error(function(){
      $scope.logPlayer = "Error Loading!"
    }).
    success(function(data){
      $scope.episode = data.Episode;
      if ($scope.episode.SrcType == 0) {
        $scope.playYoutube(playIndex);
      }
    });
  };
});
