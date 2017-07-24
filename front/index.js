import App from './components/App'
import {render} from 'react-dom'
render(<App />, document.getElementById('app'))

if (module.hot) {
  module.hot.accept('./components/App', () => {
    const App = require('./components/App').default
    ReactDOM.render(<App />, document.getElementById('app'))
  })
}
