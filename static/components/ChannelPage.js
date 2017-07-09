class Post extends React.Component {
  render() {
    const {sentDate, text} = this.props.post
    const date = moment.unix(sentDate).format('YYYY-MM-DD HH:mm');
    return <div className="panel">
      <b>{date}</b>
      <p>{text}</p>
       <hr/>
    </div>
  }
}
class ChannelPage extends React.Component {
  onMessageAdd(message) {
    const content = Object.assign(message, {
      sentDate: message.sentDate.unix(),
      channel_id: this.props.channelId
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
  render() {
    const {channels, bots, posts} = INITIAL_DATA
    const {channelId} = this.props;
    const channel = channels.find(c => c.id === channelId)
    const {title} = channel
    const filteredPosts = posts.filter(p => p.channel_id === channelId)
    return <div>
      <h2>{title}</h2>
      <MessageForm
        bots={bots}
        onAdd={m => this.onMessageAdd(m)}
      />
      <hr/>
      <h2>Posts</h2>
      {filteredPosts.map(p => <Post post={p}/>)}
    </div>
  }
}
