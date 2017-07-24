const BabiliPlugin = require('babili-webpack-plugin')
const webpack = require('webpack')
const path = require('path')
const env = process.env.NODE_ENV
const isProd = env === 'production'

module.exports = {
  devtool: isProd ? 'source-map' : 'eval-source-map',
  entry: {
    app: isProd
      ? './front/index'
      : [
          'webpack-dev-server/client?http://localhost:3001',
          'webpack/hot/only-dev-server',
          './front/index',
        ],
  },
  output: {
    path: path.resolve(__dirname, './static'),
    filename: '[name].js',
    publicPath: isProd ? '' : 'http://localhost:3001/',
  },
  module: {
    rules: [{test: /\.js$/, exclude: /node_modules/, loader: 'babel-loader'}],
  },
  plugins: [
    new webpack.ProvidePlugin({
      React: 'react',
      ReactDOM: 'react-dom',
    }),
    isProd ? false : new webpack.HotModuleReplacementPlugin(),
    isProd && new BabiliPlugin(),
  ].filter(x => x),
  devServer: {
    hot: true,
    contentBase: path.resolve(__dirname, 'dist'),
    publicPath: '/',
  },
}
