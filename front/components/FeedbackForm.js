import {Component} from 'react'
export default class FeedbackForm extends Component {
  constructor(props) {
    super(props)
    bind(this, 'renderItem')
    bind(this, 'addItem')
  }
  onItemChange(i, e) {
    const items = this.props.items.slice()
    items[i] = e.target.value
    this.props.onChange(items)
  }
  addItem() {
    const items = this.props.items.slice()
    items.push('')
    this.props.onChange(items)
  }
  onRemove(i) {
    const items = this.props.items.slice()
    items.splice(i, 1)
    this.props.onChange(items)
  }
  renderItem(item, i) {
    return (
      <div className="input-group input-group-sm">
        <input
          type="text"
          className="form-control"
          placeholder="Answer"
          value={item}
          onChange={e => this.onItemChange(i, e)}
        />
        <span className="input-group-btn">
          <button className="btn btn-default" type="button" onClick={() => this.onRemove(i)}>
            Remove
          </button>
        </span>
      </div>
    )
  }
  render() {
    const {items} = this.props

    return (
      <div>
        <label>Feedback</label>
        {items.map(this.renderItem)}
        <div>
          <br />
          <button type="button" className="btn btn-xs btn-default" onClick={this.addItem}>
            + Add answer
          </button>
        </div>
      </div>
    )
  }
}

function bind(obj, field) {
  obj[field] = obj[field].bind(obj)
}
