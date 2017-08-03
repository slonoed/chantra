import {Component} from 'react'

export default class Link extends Component {
  render() {
    const {route, children} = this.props
    const href = `/${route}`

    return (
      <a href={href}>
        {children}
      </a>
    )
  }
}
