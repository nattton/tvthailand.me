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
  $scope.playEmbed = function(playIndex) {
    if ($scope.isMobile) {
      $scope.embedWidth = 320;
      $scope.embedHeight = 240;
    } else {
      $scope.embedWidth = 800;
      $scope.embedHeight = 460;
    }
    var html = '<iframe class="col-sm-12 col-sm-offset-1 col-md-10 col-md-offset-1" frameborder="0" width="' +
    $scope.embedWidth + '" height="' + $scope.embedHeight +'" src="' +
    $scope.episode.Playlists[playIndex].sources[0].file + '" allowfullscreen></iframe>';
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
          $scope.playEmbed(playIndex);
          break;
        case 14:
          if ($scope.IsURL != true) {
            $scope.playEmbed(playIndex);
          }
          break;
      }
    });
  };
});
