tvModule.controller('VideoCtrl', function ($scope, $http, $sce, $location, $anchorScroll) {
  $scope.watch = function(id, playIndex) {
    $http.get('/api/v1/watch_otv/' + id, { params: { device: "web" }}).
    error(function(){
      $scope.iframeVideo = "Error Loading!"
    }).
    success(function(data){
      $scope.data = data;
      $scope.show_url = "/show_otv/" + data.season_detail.content_season_id;
      $scope.show_title = data.episode_detail.part_items[playIndex].name_th;
      $scope.episode_detail = data.episode_detail.detail;
      var iframe = data.episode_detail.part_items[playIndex].stream_url.replace("&lt;", "<").replace("&gt;", ">");
      $scope.iframeVideo = $sce.trustAsHtml(iframe);
    });
  };
});
