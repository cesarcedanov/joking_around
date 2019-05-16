class App extends React.Component {
  render(){
    if (this.loggedIn){
      return <loggedIn />
    } else {
      return <Home />
    }
  }
}

class Home extends React.Component {
  render(){
    return (
        <div className="container">
          <div className="col-xs-8 col-xs-offset-2 jumbotron text-center">
            <h1>Joking Around</h1>
            <p>A list of Good Jokes!</p>
            <p>Sign in to have a btter Experience.</p>
            <a 
              onClick={this.authenticate}
              className="btn btn-primary btn-lg btn-Login btn-block"
            >
              Sign In
            </a>
          </div>
        </div>
    );
  }
}


class LoggedIn extends React.Component {
  state = {
    jokes: []
  }

  render(){
    return (
      <div className="container">
        <div className="col-lg-12">
          <br />
          <span className="pull-right">
            <a onClick={this.logout}>
              Log out  
            </a>
          </span>

          <h1>Joking Around</h1>
          <p>Let me told you some jokes:</p>
          <div>
            {this.state.jokes.map( ({joke, i}) => {
              return <Joke key={i} joke={joke} />
            })}
          </div>
        </div>
      </div>
    );
  }
}


class Joke extends React.Component {
  state = {
    liked: ""
  }

  like = () => {
    // TODO
  }
  
  render() {
    const { joke } = this.props;
    return(
      <div className="col-xs-4">
        <div className="panel panel-default">
          <div className="panel-heading">
            #{joke.id}
            <span className="pull-right">
              {this.state.liked}
            </span>
          </div>
          <div className="panel-footer">
            {joke.likes} Likes &nbsp;
            <a 
              className="btn btn-default"
              onClick={this.like}
            >
              <span className="glyphicon glyphicon-thumbs-up" />
            </a>
          </div>
        </div>
      </div>
    );
  }
}

ReactDOM.render(
  <App />,
  document.querySelector('#root')
);