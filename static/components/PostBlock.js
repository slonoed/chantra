class PostBlock extends React.Component {
render() {
  const {sentDate, text} = this.props.post
  const date = moment.unix(sentDate).format('YYYY-MM-DD HH:mm');
  return <div className="panel">
    <b>{date}</b>
    <p>{text}</p>
     <hr/>
    </div>
  }
}
