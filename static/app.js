const DatePicker = window.DatePicker.default;

class Link extends React.Component {
  render() {
    const { route, children } = this.props;
    const href = "/" + route;

    return <a href={href}>{children}</a>;
  }
}

class Navbar extends React.Component {
  render() {
    return (
      <nav className="navbar navbar-inverse">
        <div className="container">
          <ul className="nav navbar-nav">
            <li>
              <Link route="bots">Bots</Link>
            </li>
            <li>
              <Link route="channels">Channels</Link>
            </li>
            <li>
              <Link route="messages">Messages</Link>
            </li>
          </ul>
          <p className="navbar-text navbar-right">
            <a href="/auth/logout" className="navbar-link">Logout</a>
          </p>
        </div>
      </nav>
    );
  }
}

class BotsPage extends React.Component {
  constructor(props) {
    super(props);
    this.fieldId = "bots-page-form-input-" + Math.floor(Math.random() * 100000);
    this.state = {
      bots: INITIAL_DATA.bots || [],
      loading: false
    };
  }
  onSubmit(e) {
    e.preventDefault();
    const token = this.tokenInput.value;
    this.tokenInput.value = "";
    this.setState({ loading: true });
    fetch("/api/addBot", {
      method: "POST",
      credentials: "include",
      body: JSON.stringify({ token })
    })
      .then(r => r.json())
      .then(data => {
        const bots = this.state.bots.filter(b => b.id !== data.bot.id);
        bots.unshift(data.bot);
        this.setState({ bots });
        this.setState({ loading: false });
      })
      .catch(e => {
        this.setState({ loading: true });
      });
  }
  renderBot(bot) {
    return (
      <div className="col-xs-12 col-md-6 col-lg-4 panel panel-default">
        <div className="panel-body">
          {bot.firstName}
        </div>
      </div>
    );
  }
  render() {
    const { bots, loading } = this.state;

    return (
      <div>
        <div className="row">
          <h1>Bots</h1>
        </div>
        <div className="row">
          <div className="col-xs-12 col-sm-8">
            {bots.map(b => this.renderBot(b))}
          </div>
          <div className="col-xs-12 col-sm-4 well">
            <form onSubmit={e => this.onSubmit(e)} className="form">
              <div className="form-group">
                <label htmlFor={this.fieldId}>Bot token</label>
                <input
                  type="text"
                  className="form-control"
                  required
                  ref={i => this.tokenInput = i}
                  disabled={loading}
                />
              </div>
              <button
                type="submit"
                className="btn btn-primary"
                disabled={loading}
              >
                Add new bot
              </button>
            </form>
          </div>
        </div>
      </div>
    );
  }
}
class ChannelsPage extends React.Component {
  constructor(props) {
    super(props)
    this.state = {botId: null}
  }
  onSubmit(e) {
    e.preventDefault()

    if (!this.state.botId) {
      return
    }

    const content = {
      name: this.nameInput.value,
      botId: this.state.botId
    }

    fetch('/api/addChannel', {
      method: 'POST',
      body: JSON.stringify(content),
      credentials: 'include'
    })
    .then(r => r.json())
    .then(r => {
      if (!r.ok) {
        throw new Error(r.error || 'Unknown error');
      }
    })
    .catch(e => {
      alert(e.toString());
    });
  }
  onBotSelect(e) {
    this.setState({botId: e.target.value})
  }
  renderBot(bot) {
    return <option value={bot.id}>{bot.firstName} (@{bot.username})</option>;
  }
  render() {
    const bots = INITIAL_DATA.bots || []
    const channels = INITIAL_DATA.channels || []
    const {botId} = this.state

    return <div>
      <h1>Channels</h1>
      <form onSubmit={e => this.onSubmit(e)} className="form">
        <div className="form-group">
          <label>Channel name (with @)</label>
          <input
            type="text"
            className="form-control"
            required
            ref={i => this.nameInput = i}
          />
        </div>
        <div className="form-group">
          <label htmlFor="bot">Bot</label>
          <select
            className="form-control"
            id="bot"
            value={botId}
            onChange={e => this.onBotSelect(e)}>
            <option value="">Select bot</option>;
            {bots.map(c => this.renderBot(c))}
          </select>
        </div>
        <button
          type="submit"
          className="btn btn-primary"
          disabled={!botId}
        >
          Add new channel
        </button>

      </form>
      <hr/>
      {channels.map(c => (<div>{c.title} - {c.username}</div>))}
    </div>
  }
}
class MessagesPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      bots: INITIAL_DATA.bots || [],
      posts: INITIAL_DATA.posts || [],
      channels: INITIAL_DATA.channels || []
    };
  }
  onMessageAdd(message) {
    const content = Object.assign(message, {
      sentDate: message.sentDate.unix()
    });
    fetch('/api/addPost', {
      method: 'POST',
      body: JSON.stringify(content),
      credentials: 'include'
    })
    .then(r => r.json())
    .then(r => {
      if (!r.ok) {
        throw new Error(r.error || 'Unknown error');
      }
    })
    .catch(e => {
      alert(e.toString());
    });
  }
  renderAnswer(a) {
    return <div>
      {a.text} — {a.votes}
    </div>;
  }
  renderPost(post) {
    const date = moment.unix(post.sentDate).format('YYYY-MM-DD HH:mm');
    return <div className="panel panel-default" key={post.id}>
      <div className="panel-heading">
        {date}
      </div>
      <div className="panel-body">
        {post.text}
        <hr/>
        {post.answers.map(a => this.renderAnswer(a))}
      </div>
    </div>;
  }
  render() {
    const { bots, channels, posts } = this.state;
    const hasChannels = channels.length > 0;
    const byDate = (b, a) => a.sentDate - b.sentDate;
    // TODO if no channels don't render form
    return (
      <div className="row">
        <div className="col-xs-12 col-md-6">
          {hasChannels &&
            <MessageForm
              bots={bots}
              channels={channels}
              onAdd={m => this.onMessageAdd(m)}
            />}
          {hasChannels || <div>No channels to send message</div>}
        </div>
        <div className="col-xs-12 col-md-6">
          {posts.sort(byDate).map(p => this.renderPost(p))}
        </div>
      </div>
    );
  }
}
class NotFoundPage extends React.Component {
  render() {
    return <div>NotFound page</div>;
  }
}

class App extends React.Component {
  render() {
    let Page = {
      "/": MessagesPage,
      "/bots": BotsPage,
      "/channels": ChannelsPage,
      "/messages": MessagesPage
    }[location.pathname];
    if (!Page) {
      Page = NotFoundPage;
    }

    return (
      <div>
        <Navbar />
        <div className="container">
          <Page />
        </div>
      </div>
    );
  }
}

ReactDOM.render(<App />, document.getElementById("app"));
