class MessageForm extends React.Component {
  constructor(props) {
    super(props);
    const { channels} = props;
    this.state = {
      channelId: channels[0].id,
      botId: null,
      sentDate: moment(),
      feedback: null
    };
  }
  onSubmit(e) {
    e.preventDefault();
    const date = moment(this.state.sentDate);
    const hours = parseInt(this.hours.value, 10);
    const minutes = parseInt(this.minutes.value, 10);
    date.hours(hours);
    date.minutes(minutes);
    this.props.onAdd({
      channel_id: this.state.channelId,
      sentDate: date,
      text: this.text.value,
      parseMode: this.md.checked ? 'Markdown' : '',
      withPotification: !this.disableNotification.checked,
      withPreview: !this.disableWebPagePreview.checked,
      answers: this.state.feedback ? this.state.feedback.map(a => ({
        text: a
      })) : []
    });
  }
  onDateChange(date) {
    this.setState({ sentDate: date });
  }
  onChannelSelect(e) {
    this.setState({
      channelId: e.target.value.id,
    });
  }
  renderChannel(channel) {
    return (
      <option value={channel.id}>{channel.title} (@{channel.username})</option>
    );
  }
  renderBot(bot) {
    return <option value={bot.id}>{bot.firstName} (@{bot.username})</option>;
  }
  render() {
    const { channels} = this.props;
    const { channelId, sentDate, feedback } = this.state;
    const channel = channels.find(c => c.id === channelId);
    const hours = moment().hours();
    const minutes = moment().minutes();

    const feedbackEl = feedback == null
      ? <button
          className="btn btn-link"
          onClick={() => this.setState({ feedback: [] })}
        >
          + Add feedback
        </button>
      : <FeedbackForm
          items={feedback}
          onChange={f => this.setState({ feedback: f })}
        />;

    return (
      <form onSubmit={e => this.onSubmit(e)} className="form">
        <div className="row">
          <div className="col-xs-12 col-md-6">
            <div className="form-group">
              <label htmlFor="channel">Channel</label>
              <select
                className="form-control"
                id="channel"
                value={channelId}
                onChange={e => this.onChannelSelect(e)}
              >
                {channels.map(c => this.renderChannel(c))}
              </select>
            </div>
          </div>
        </div>
        <div className="form-group">
          <label htmlFor="text">Text</label>&nbsp;
          <label>
            <input type="checkbox" ref={i => this.md = i}/> Use markdown
          </label>
          <textarea className="form-control" id="" name="" rows="10"
              ref={i => this.text = i} />
        </div>
        <div className="form-group">
          <div className="checkbox">
            <label>
              <input type="checkbox" ref={i => this.disableNotification = i}/>
              Disable notification
            </label>
          </div>
          <div className="checkbox">
            <label>
              <input type="checkbox" ref={i => this.disableWebPagePreview = i}/>
              Disable web page preview
            </label>
          </div>
        </div>
        <div className="form-group">
          {feedbackEl}
        </div>
        <div className="form-group">
          <label>Schedule sending</label>
          <div className="form-control-static row">
            <div className="col-xs-4">
              <DatePicker
                className="form-control"
                selected={sentDate}
                dateFormat="YYYY-MM-DD"
                minDate={moment()}
                locale="en-gb"
                onChange={d => this.onDateChange(d)}
              />
              <span className="help-block">Date</span>
            </div>
            <div className="col-xs-4">
              <input
                className="form-control"
                type="text"
                defaultValue={hours}
                ref={i => this.hours = i}
              />
              <span className="help-block">Hours</span>
            </div>
            <div className="col-xs-4">
              <input
                className="form-control"
                type="text"
                defaultValue={minutes}
                ref={i => this.minutes = i}
              />
              <span className="help-block">Minutes</span>
            </div>
          </div>
        </div>

        <button type="submit" className="btn btn-primary">
          Schedule message
        </button>
      </form>
    );
  }
}
