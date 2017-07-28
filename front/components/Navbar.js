import {Component} from 'react'
import Link from './Link'

export default class Navbar extends Component {
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
    )
  }
}
