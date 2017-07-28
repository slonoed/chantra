import moment from 'moment'
import {Component} from 'react'

const Answer = ({answer}) =>
  <tr>
    <td>
      {answer.text}
    </td>
    <td>
      {answer.votes}
    </td>
  </tr>

export default class Post extends Component {
  renderAnswers(answers) {
    return (
      <table className="table table-bordered table-condensed">
        <thead>
          <tr>
            <th>Answer</th>
            <th>Votes</th>
          </tr>
        </thead>
        <tbody>
          {answers.map(a => <Answer answer={a} />)}
        </tbody>
      </table>
    )
  }
  render() {
    const {sentDate, text, answers} = this.props.post
    const date = moment.unix(sentDate).format('YYYY-MM-DD HH:mm')
    const isNew = moment.unix(this.props.post.sentDate).isAfter(moment().subtract(10, 'seconds'))

    const style = {
      backgroundColor: isNew ? '#fffcda' : null,
      borderBottom: '1px solid #ddd',
      padding: '5px',
    }

    return (
      <div style={style}>
        <b>
          {date}
        </b>
        <p>
          {text}
        </p>
        {answers.length > 0 ? this.renderAnswers(answers) : null}
      </div>
    )
  }
}
