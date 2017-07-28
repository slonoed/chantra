import App from './components/App'
import {render} from 'react-dom'
import {create} from './store'
import {Provider} from 'react-redux'

const store = create()

render(
  <Provider store={store}>
    <App />
  </Provider>,
  document.getElementById('app')
)

if (module.hot) {
  module.hot.accept('./components/App', () => {
    const App = require('./components/App').default
    ReactDOM.render(
      <Provider store={store}>
        <App />
      </Provider>,
      document.getElementById('app')
    )
  })
}
