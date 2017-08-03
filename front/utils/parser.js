/*
 * Returns list of data for markdown
 */
export default function parse(markdownText) {
  const md = markdownText
  // TODO implement this
  let result = []
  let acc = ''
  let mode = null
  let reg = /\*.+?\*/
  for (let char of md) {
    if (char !== '*' && char !== '_') {
      acc += char
    } else if (mode == null && md.match(reg) != null) {
      result.push(acc)
      acc = ''
      mode = char === '*' ? 'bold' : 'italic'
    } else {
      result.push({type: mode, content: acc})
      acc = ''
      mode = null
    }
  }
  result.push(acc)
  return result
}
