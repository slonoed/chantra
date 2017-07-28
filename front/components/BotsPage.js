import {Component} from 'react'

export default class BotsPage extends Component {
  constructor(props) {
    super(props)
    this.fieldId = 'bots-page-form-input-' + Math.floor(Math.random() * 100000)
    this.state = {
      bots: INITIAL_DATA.bots || [],
      loading: false,
    }
  }
  onSubmit(e) {
    e.preventDefault()
    const token = this.tokenInput.value
    this.tokenInput.value = ''
    this.setState({loading: true})
    fetch('/api/addBot', {
      method: 'POST',
      credentials: 'include',
      body: JSON.stringify({token}),
    })
      .then(r => r.json())
      .then(data => {
        const bots = this.state.bots.filter(b => b.id !== data.bot.id)
        bots.unshift(data.bot)
        this.setState({bots})
        this.setState({loading: false})
      })
      .catch(e => {
        this.setState({loading: true})
      })
  }
  renderBot(bot) {
    return (
      <div className="col-xs-12 col-md-6 col-lg-4 panel panel-default">
        <div className="panel-body">
          {bot.firstName}
        </div>
      </div>
    )
  }
  render() {
    const {bots, loading} = this.state

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
              <button type="submit" className="btn btn-primary" disabled={loading}>
                Add new bot
              </button>
            </form>
          </div>
        </div>
      </div>
    )
  }
}
