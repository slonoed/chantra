import parse from '../parser'

const testData = [
  {input: 'Hello', output: ['Hello']},
  {input: 'Hello *Vova*', output: ['Hello', {type: 'bold', content: ['Vova']}]},
  {input: 'Hello *_Vova_*', output: ['Hello', {type: 'bold', content: ['_Vova_']}]},
  {input: 'Hello _*Vova*_', output: ['Hello', {type: 'italic', content: ['*Vova*']}]},
  {input: 'Рандомный текст *Жирный* _Italic_ лишняя * лишний _ и снова _italic_', output: 'ERROR'},
  {input: 'Hello _*Vova**_', output: ['Hello', {type: 'italic', content: ['*Vova**']}]},
  {input: 'Hello *_Vova__*', output: ['Hello', {type: 'bold', content: ['_Vova__']}]},
  {
    input: '[text](http://www.example.com)',
    output: ['[text]', {type: 'link', content: 'text', href: 'http://www.example.com'}],
  },
  {
    input: '*[text]*(http://www.example.com)',
    output: [
      {type: 'bold', content: ['[text]']},
      {type: 'link', content: 'http://www.example.com', href: 'http://www.example.com'},
    ],
  },
  {
    input: '_[text]_(http://www.example.com)',
    output: [
      {type: 'italic', content: ['[text]']},
      {type: 'link', content: 'http://www.example.com', href: 'http://www.example.com'},
    ],
  },
  {input: '_[text](http://www.example.com)_', output: [{type: 'italic', content: ['[text]']}]},
  {
    input: '*[text](http://www.example.com)*',
    output: [{type: 'bold', content: ['[text](http://www.example.com)']}],
  },
  {
    input: '```[text](http://www.example.com)```',
    output: [{type: 'pre', content: ['[text](http://www.example.com)']}],
  },
  {input: '```*толстый текст*```', output: [{type: 'pre', content: ['текст*']}]},
  {input: '```*текст*```', output: [{type: 'pre', content: ['*текст*']}]},
  {input: '```* текст*```', output: [{type: 'pre', content: ['текст*']}]},
  {input: '```_ текст_```', output: [{type: 'pre', content: ['текст_']}]},
  {input: '```_текст_```', output: [{type: 'pre', content: ['_текст_']}]},
  {input: '```_текст   _```', output: [{type: 'pre', content: ['   _']}]},
  {input: '```*текст   *```', output: [{type: 'pre', content: ['   *']}]},
  {input: '```*_текст_*```', output: [{type: 'pre', content: ['*_текст_*']}]},
  {
    input: '```текст``` *жирный текст*  _italic text_',
    output: [
      {type: 'pre', content: ['текст']},
      {type: 'bold', content: ['жирный текст']},
      {type: 'italic', content: ['italic text']},
    ],
  },
  {input: '*```*толстый текст*```*', output: ['```толстый текст```']},
  {input: '_```_толстый текст_```_', output: ['```толстый текст```']},
  {input: '_```_толстый текст_```_', output: ['```толстый текст```']},
  {
    input: '_```*толстый текст*```_',
    output: [{type: 'italic', content: ['```*толстый текст*```']}],
  },
  {input: '*```_толстый текст_```*', output: [{type: 'bold', content: ['```_толстый текст_```']}]},
]

// Create test for each item
testData.forEach(({input, output}) => {
  test(`Check text << ${input} >> `, () => {
    expect(parse(input)).toEqual(output)
  })
})
