const AUTH0_CLIENT_ID="gW7qWkx0DrG2cTp2Dxa1ZE8RqWX8DWEy"
const AUTH0_DOMAIN="joking-around.auth0.com"
const AUTH0_API_AUDIENCE="https://joking-around.auth0.com/api/v2/"
const AUTH0_CALLBACK_URL="http://localhost:9000"

class App extends React.Component {
  parseHash() {
    this.auth0 = new auth0.WebAuth({
      domain: AUTH0_DOMAIN,
      clientID: AUTH0_CLIENT_ID
    });

    this.auth0.parseHash(window.location.hash, (err, authResult) => {
      err => console.log(err);
      if (
        authResult !== null &&
        authResult.accessToken !== null &&
        authResult.idToken !== null
      ) {
        localStorage.setItem("access_token", authResult.accessToken);
        localStorage.setItem("id_token", authResult.idToken);
        localStorage.setItem("profile", JSON.stringify(authResult.idTokenPayload));
        window.location = window.location.href.substr(
          0,
          window.location.href.indexOf("#")
        );
      }
    });
  }

  setup(){
    $.ajaxSetup({
      beforeSend: (r) => {
        if (localStorage.getItem("access_token")) {
          r.setRequestHeader(
            "Authorization",
            `Bearer ${localStorage.getItem("access_token")}`
          );
        }
      }
    });
  }

  setState(){
    if (localStorage.getItem("id_token")) {
      this.loggedIn = true;
    } else {
      this.loggedIn = false;
    }
  }

  componentWillMount(){
    this.setup();
    this.parseHash();
    this.setState();
  }

  render(){
    if (this.loggedIn){
      return <LoggedIn />
    } else {
      return <Home />
    }
  }
}

class Home extends React.Component {
  constructor(props){
    super(props);
    this.authenticate = this.authenticate.bind(this);
  }

  authenticate() {
    this.WebAuth = new auth0.WebAuth({
      domain: AUTH0_DOMAIN,
      clientID: AUTH0_CLIENT_ID,
      scope: "openid profile",
      audience: AUTH0_API_AUDIENCE,
      responseType: "token id_token",
      redirectUri: AUTH0_CALLBACK_URL
    });
    this.WebAuth.authorize();
  }

  render(){
    return (
        <div className="container">
          <div className="col-xs-8 col-xs-offset-2 jumbotron text-center">
            <h1>Joking Around</h1>
            <p>A list of good Jokes!</p>
            <p>Sign in to have a better Experience.</p>
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

  logout = () => {
    localStorage.removeItem("id_token");
    localStorage.removeItem("access_token");
    localStorage.removeItem("profile");
    location.reload();
  }

  fetchJokes = () => {
    $.get("http://localhost:9000/api/jokes", res => {
      this.setState({
        jokes: res
      });
    });
  }
  
  componentDidMount() {
    this.fetchJokes();
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
            {this.state.jokes.map( (joke, i) => {
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
    liked : false,
    joke: this.props.joke || {}
  }

  like = () => {
    this.likeJoke(this.state.joke);
  }

  likeJoke = (joke) => {
    $.post(`http://localhost:9000/api/jokes/like/${joke.id}`, joke => {
      this.setState( { joke, liked: true });
    })
  }
  
  render() {
    const { joke } = this.state;
    return(
      <div className="col-xs-4">
        <div className="panel panel-default">
          <div className="panel-heading">
            #{joke.id}
            <span className="pull-right">
              {this.state.liked ? 'Liked' : ""}
            </span>
          </div>
          <div className="panel-body joke-hld">{this.props.joke.content}</div>
          <div className="panel-footer">
            {joke.likes} times- It made me laugh &nbsp;
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