var path = require('path');
var webpack = require('webpack');

module.exports = {
  entry: [
    "webpack-dev-server/client?http://0.0.0.0:8080",
    "webpack/hot/only-dev-server",
    './app/index.js'
  ],
  devtool: "source-map",
  output: {
    path: path.join(__dirname, "build"),
    publicPath: '/',
    filename: 'app.js',
    crossOriginLoading: "use-credentials"
  },
  resolveLoader: {
    modulesDirectories: ['node_modules']
  },
  resolve: {
    root: path.resolve('./'),
    extensions: ['', '.js']
  },
  node: {
    fs: "empty"
  },
  devServer: {
    contentBase: 'http://localhost:8000',
    headers: { "Access-Control-Allow-Origin": "*" }
  },
  module: {
    loaders: [
      {
        test: /\.js?$/,
        exclude: /(node_modules|bower_components)/,
        loaders: ['react-hot', 'babel'],
        include: path.join(__dirname, 'app')
      }
    ]
  }
};
