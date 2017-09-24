const Paragraph = ({text}) => {
  const lines = text.split('\n')
  const elements = lines.reduce((all, line) => all.concat([line, <br />]), [])

  elements.pop()

  return (
    <div>
      {elements}
    </div>
  )
}

export default Paragraph
