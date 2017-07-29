import {Component} from 'react'
import ChannelPage from './ChannelPage'
import FeedbackForm from './FeedbackForm'
import ChannelsPage from './ChannelsPage'
import Navbar from './Navbar'
import NotFoundPage from './NotFoundPage'
import Notifications from './Notifications'
import BotsPage from './BotsPage'

export default class App extends Component {
  render() {
    let Page = NotFoundPage
    let channelId = null
    // for channels page
    const match = /\/channels\/([a-z0-9]+)/.exec(location.pathname)
    if (match && match[1]) {
      channelId = match[1]
      Page = ChannelPage
    } else {
      Page = {
        '/bots': BotsPage,
        '/channels': ChannelsPage,
        '/': ChannelsPage,
      }[location.pathname]
    }

    if (!Page) {
      Page = NotFoundPage
    }

    return (
      <div>
        <Navbar />
        <div className="container">
          <Page channelId={channelId} />
        </div>
        <Notifications />
      </div>
    )
  }
}
