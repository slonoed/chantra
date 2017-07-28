import {Component} from 'react'
import Post from './Post'
import Alert from './Alert'
import MessageForm from './MessageForm'
import {connect} from 'react-redux'
import {scheduleNewPost} from '../reducers/posts'

class ChannelPage extends Component {
  onMessageAdd(message) {
    this.props.scheduleNewPost(message, this.props.channelId)
  }

  render() {
    const {channels, bots} = INITIAL_DATA
    const {channelId, posts} = this.props
    const channel = channels.find(c => c.id === channelId)
    const {title} = channel
    const filteredPosts = Object.values(posts).filter(p => p.channel_id === channelId)
    const byDate = (b, a) => a.sentDate - b.sentDate

    const columnStyle = {flex: 1}

    return (
      <div>
        <h2>
          {title}
        </h2>
        <div className="row">
          <div className="col-md-6">
            <MessageForm bots={bots} onAdd={m => this.onMessageAdd(m)} />
          </div>
          <div className="col-md-6">
            <h2>Posts</h2>
            {filteredPosts.sort(byDate).map(p => <Post key={p.id} post={p} />)}
          </div>
        </div>
      </div>
    )
  }
}

const mapStateToProps = state => ({
  posts: state.posts.posts,
})
const mapDispatchToProps = {scheduleNewPost}

export default connect(mapStateToProps, mapDispatchToProps)(ChannelPage)
