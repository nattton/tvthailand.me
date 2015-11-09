var ShowItem = React.createClass({
  render: function () {
    return React.createElement(
      "div",
      { className: "col-xs-12 col-sm-3 placeholder" },
      React.createElement(
        "div",
        { className: "thumbnail" },
        React.createElement(
          "a",
          { href: '/show/' + this.props.data.id + '/' + this.props.data.title },
          React.createElement("img", { src: this.props.data.thumbnail, className: "img-responsive" })
        ),
        React.createElement(
          "div",
          { className: "caption" },
          React.createElement(
            "h4",
            null,
            this.props.data.title
          )
        )
      )
    );
  }
});

var ShowList = React.createClass({
  render: function () {
    return React.createElement(
      "div",
      { className: "showList" },
      this.props.shows.map(function (show) {
        return React.createElement(ShowItem, { key: show.id, data: show });
      })
    );
  }
});

var ShowBox = React.createClass({
  loadShows: function () {
    $.ajax({
      url: '/ajax/' + this.props.typeMode + "/" + this.props.typeId,
      data: { "offset": this.state.shows.length },
      dataType: 'json',
      success: (function (data) {
        if (data.shows.length == 0) $('#load_more').hide();
        this.setState({ shows: this.state.shows.concat(data.shows) });
      }).bind(this),
      error: (function (xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }).bind(this)
    });
  },
  componentWillMount: function () {
    this.setState({ shows: this.props.shows });
  },
  render: function () {
    return React.createElement(
      "div",
      { className: "showBox" },
      React.createElement(
        "div",
        { className: "row placeholders" },
        React.createElement(ShowList, { shows: this.state.shows })
      ),
      React.createElement(
        "ul",
        { id: "load_more", className: "pager" },
        React.createElement(
          "li",
          null,
          React.createElement(
            "a",
            { className: "btn", onClick: this.loadShows },
            "Load more ..."
          )
        )
      )
    );
  }
});
if (TYPE_MODE != "") {
  ReactDOM.render(React.createElement(ShowBox, { typeMode: TYPE_MODE, typeId: TYPE_ID, shows: SHOWS }), document.getElementById('show_content'));
} else {
  $('#load_more').hide();
}