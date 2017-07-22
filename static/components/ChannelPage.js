class Alert extends React.Component {
  render() {
    return <h1>Done!</h1>
  }
}

class ChannelPage extends React.Component {
  constructor(props) { 
    super(props); 
    const {posts} = INITIAL_DATA
    this.state = {posts,showMe: false};
  }

  onMessageAdd(message) {
   const {posts} = this.state
   const content = Object.assign(message, {
      sentDate: message.sentDate.unix(),
      channel_id: this.props.channelId
    
  })
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
     else {
      const newPosts = posts.concat(r.post);
      this.setState({ posts: newPosts , showMe:true});
      setTimeout(() => this.setState({showMe: false}), 3000)
      }
    })
    .catch(e => {
      alert(e.toString());
    }); 
  }

 
  render() {
    const {channels, bots} = INITIAL_DATA
    const {channelId} = this.props;
    const channel = channels.find(c => c.id === channelId)
    const {posts, showMe} = this.state
    const filteredPosts = this.state.posts.filter(p => p.channel_id === channelId)
    const byDate = (b, a) => a.sentDate - b.sentDate;
    const {title} = channel
      return <div>
         {showMe===true? <Alert/>: ""}
        <h2>{title}</h2>
        <MessageForm
        bots={bots}
        onAdd={m => this.onMessageAdd(m)}
      />
     <hr/>
      <h2>Posts</h2>
      {filteredPosts.sort(byDate).map(p => <PostBlock post={p}/>)}
    </div>
  }
}