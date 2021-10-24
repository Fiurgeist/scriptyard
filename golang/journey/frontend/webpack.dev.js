const merge = require('webpack-merge');
const common = require('./webpack.common.js');

module.exports = (env) => merge(common, {
    mode: 'development',
    // Reload On File Change
    watch: env && env.watch,
    // Development Tools (Map Errors To Source File)
    devtool: 'inline-source-map',
});
