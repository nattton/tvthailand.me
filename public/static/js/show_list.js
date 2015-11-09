var ShowItem = React.createClass({
  render: function() {
    return (
      <div className="col-xs-12 col-sm-3 placeholder">
        <div className="thumbnail">
        <a href={'/show/' + this.props.data.id + '/' + this.props.data.title}>
        <img src={this.props.data.thumbnail} className="img-responsive" />
        </a>
        <div className="caption"><h4>{this.props.data.title}</h4></div>
        </div>
      </div>
    );
  }
});

var ShowList = React.createClass({
  render : function() {
    return (
      <div className="showList">
      {this.props.shows.map(function(show){
        return <ShowItem key={show.id} data={show} />;
      })}
      </div>
    );
  }
});

var ShowBox = React.createClass({
  loadShows: function() {
    $.ajax({
      url: '/ajax/' + this.props.typeMode + "/" + this.props.typeId,
      data: {"offset": this.state.shows.length},
      dataType: 'json',
      success: function(data) {
        if (data.shows.length == 0) $('#load_more').hide()
        this.setState({shows: this.state.shows.concat(data.shows)});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  },
  componentWillMount: function() {
    this.setState({shows: this.props.shows});
  },
  render: function() {
    return (
      <div className="showBox">
        <div className="row placeholders">
          <ShowList shows={this.state.shows}/>
        </div>
        <ul id="load_more" className="pager">
          <li><a className="btn" onClick={this.loadShows}>Load more ...</a></li>
        </ul>
      </div>
    );
  }
});
if (TYPE_MODE != "") {
  ReactDOM.render(
    <ShowBox typeMode={TYPE_MODE} typeId={TYPE_ID} shows={SHOWS} />,
    document.getElementById('show_content')
  );
} else {
  $('#load_more').hide()
}
