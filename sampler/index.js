const Listener = require('./Listener')

const l = new Listener('../app/res/priv/model.tflite', '../app/res/priv/huge-vocabulary.scorer', 'computer')
l.on('command', async (result) => {
  console.log(result)
})
l.start()
