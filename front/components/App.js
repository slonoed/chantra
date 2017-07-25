import { Component } from "react";
import ChannelPage from "./ChannelPage";
import FeedbackForm from "./FeedbackForm";

class Link extends Component {
  render() {
    const { route, children } = this.props;
    const href = "/" +
      route;

    return (
      <a href={href:href}>
        {children}
      </a>
    );
  }
}

class Navbar extends Component {
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
          </ul>
          <p className="navbar-text navbar-right">
            <a href="/auth/logout" className="navbar-link">
              Logout
            </a>
          </p>
        </div>
      </nav>
    );
  }
}

class BotsPage extends Component {
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
                  ref={i => (this.tokenInput = i)}
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
class ChannelsPage extends Component {
  constructor(props) {
    super(props);
    this.state = { botId: null };
  }
  onSubmit(e) {
    e.preventDefault();

    if (!this.state.botId) {
      return;
    }

    const content = {
      name: this.nameInput.value,
      botId: this.state.botId
    };

    fetch("/api/addChannel", {
      method: "POST",
      body: JSON.stringify(content),
      credentials: "include"
    })
      .then(r => r.json())
      .then(r => {
        if (!r.ok) {
          throw new Error(r.error || "Unknown error");
        }
      })
      .catch(e => {
        alert(e.toString());
      });
  }
  onBotSelect(e) {
    this.setState({ botId: e.target.value });
  }
  renderChannelLink(channel) {
    const link = `/channels/${channel.id}`;
    const style = {
      display: "block"
    };
    return (
      <a href={link} style={style}>
        {channel.title}
      </a>
    );
  }
  renderBot(bot) {
    return (
      <option value={bot.id}>
        {bot.firstName} (@{bot.username})
      </option>
    );
  }
  render() {
    const bots = INITIAL_DATA.bots || [];
    const channels = INITIAL_DATA.channels || [];
    const { botId } = this.state;

    return (
      <div>
        <h1>Channels</h1>
        <form onSubmit={e => this.onSubmit(e)} className="form">
          <div className="form-group">
            <label>Channel name (with @)</label>
            <input
              type="text"
              className="form-control"
              required
              ref={i => (this.nameInput = i)}
            />
          </div>
          <div className="form-group">
            <label htmlFor="bot">Bot</label>
            <select
              className="form-control"
              id="bot"
              value={botId}
              onChange={e => this.onBotSelect(e)}
            >
              <option value="">Select bot</option>;
              {bots.map(c => this.renderBot(c))}
            </select>
          </div>
          <button type="submit" className="btn btn-primary" disabled={!botId}>
            Add new channel
          </button>
        </form>
        <hr />
        {channels.map(c => this.renderChannelLink(c))}
      </div>
    );
  }
}

class NotFoundPage extends Component {
  render() {
    return <div>NotFound page</div>;
  }
}

export default class App extends Component {
  render() {
    let Page = NotFoundPage;
    let channelId = null;
    // for channels page
    const match = /\/channels\/([a-z0-9]+)/.exec(location.pathname);
    if (match && match[1]) {
      channelId = match[1];
      Page = ChannelPage;
    } else {
      Page = {
        "/bots": BotsPage,
        "/channels": ChannelsPage,
        "/": ChannelsPage
      }[location.pathname];
    }

    if (!Page) {
      Page = NotFoundPage;
    }

    return (
      <div>
        <Navbar />
        <div className="container">
          <Page channelId={channelId} />
        </div>
      </div>
    );
  }
}
