const merge = require('webpack-merge');
const common = require('./webpack.common.js');

module.exports = merge(common, {
    mode: 'production',
    // Development Tools (Map Errors To Source File)
    devtool: 'source-map',
});
