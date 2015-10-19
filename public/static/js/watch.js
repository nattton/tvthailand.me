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
  $scope.playDailymotion = function(playIndex) {
    var html = '<iframe class="col-sm-12 col-sm-offset-1 col-md-10 col-md-offset-1" frameborder="0" width="480" height="480" src="' +
                $scope.episode.Playlists[playIndex].sources[0].file + '" ></iframe>';
    $scope.iframeVideo = $sce.trustAsHtml(html);
  }
  $scope.watch = function(id, playIndex) {
    $http.get('/api/v1/watch/' + id, { params: { device: "web" }}).
    error(function(){
      $scope.iframeVideo = "Error Loading!"
    }).
    success(function(data){
      $scope.episode = data.Episode;
      switch ($scope.episode.SrcType) {
        case  0:
          $scope.playYoutube(playIndex);
          break;
        case 1:
          $scope.playDailymotion(playIndex);
          break;
      }
    });
  };
});
