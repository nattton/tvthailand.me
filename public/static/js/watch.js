tvModule.controller('VideoCtrl', function ($scope, $http, $sce) {
  $scope.playYoutube = function(playIndex) {
    var playerInstance = jwplayer("player");
    playerInstance.setup({
      autostart: true,
       playlist: $scope.episode.Playlists
    });
    playerInstance.on('ready', function(){
      playerInstance.playlistItem(playIndex);
    });
    playerInstance.on('play', function(){
      if ($scope.episode.Playlists.length > 1) {
        $scope.episode_title = $scope.episode.Title + " Part " + (playerInstance.getPlaylistIndex()+1) + "/" + $scope.episode.Playlists.length;
      }
    });
  }
  $scope.playDailymotion = function(playIndex) {
    var html = '<iframe class="col-sm-12 col-sm-offset-1 col-md-10 col-md-offset-1" frameborder="0" width="480" height="480" src="' +
                $scope.episode.Playlists[playIndex].sources[0].file + '" ></iframe>';
    $scope.iframeVideo = $sce.trustAsHtml(html);
  }
  $scope.createWebList = function() {
  }
  $scope.watch = function(id, playIndex) {
    $http.get('/api/v1/watch/' + id, { params: { device: "web" }}).
    error(function(){
      $scope.iframeVideo = "Error Loading!"
    }).
    success(function(data){
      $scope.show = data.show;
      $scope.episode = data.episode;
      $scope.show_url = "/show/" + data.show.ID +"/" + data.show.Title;
      $scope.show_title = data.show.Title;

      if (data.episode.Playlists.length > 1) {
        $scope.episode_title = data.episode.Title + " Part " + (playIndex+1) + "/" + data.episode.Playlists.length;
      } else {
        $scope.episode_title = data.episode.Title;
      }

      switch ($scope.episode.SrcType) {
        case  0:
          $scope.playYoutube(playIndex);
          break;
        case 1:
          $scope.playDailymotion(playIndex);
          break;
        default:
            $scope.createWebList();
      }
    });
  };
});
