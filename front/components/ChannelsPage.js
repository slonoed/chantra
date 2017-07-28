import {Component} from 'react'

export default class ChannelsPage extends Component {
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
      botId: this.state.botId,
    }

    fetch('/api/addChannel', {
      method: 'POST',
      body: JSON.stringify(content),
      credentials: 'include',
    })
      .then(r => r.json())
      .then(r => {
        if (!r.ok) {
          throw new Error(r.error || 'Unknown error')
        }
      })
      .catch(e => {
        alert(e.toString())
      })
  }
  onBotSelect(e) {
    this.setState({botId: e.target.value})
  }
  renderChannelLink(channel) {
    const link = `/channels/${channel.id}`
    const style = {
      display: 'block',
    }
    return (
      <a href={link} style={style}>
        {channel.title}
      </a>
    )
  }
  renderBot(bot) {
    return (
      <option value={bot.id}>
        {bot.firstName} (@{bot.username})
      </option>
    )
  }
  render() {
    const bots = INITIAL_DATA.bots || []
    const channels = INITIAL_DATA.channels || []
    const {botId} = this.state

    return (
      <div>
        <h1>Channels</h1>
        <form onSubmit={e => this.onSubmit(e)} className="form">
          <div className="form-group">
            <label>Channel name (with @)</label>
            <input type="text" className="form-control" required ref={i => (this.nameInput = i)} />
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
    )
  }
}
