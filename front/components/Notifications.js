import {Component} from 'react'
import {connect} from 'react-redux'

const Notification = ({text}) =>
  <div className="Notification">
    {text}
  </div>

class Notifications extends Component {
  render() {
    const {notifications} = this.props

    return (
      <div>
        {notifications.map(n => <Notification key={n} text={n} />)}
      </div>
    )
  }
}

export default connect(state => ({
  notifications: Object.values(state.notifications).filter(x => x),
}))(Notifications)
