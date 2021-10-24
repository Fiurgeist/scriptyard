const path = require('path');

const config = {
    entry: './main.tsx',
    resolve: {
        extensions: ['.js', '.jsx', '.ts', '.tsx'],
    },
    output: {
        library: 'wod',
        path: path.resolve(__dirname, '../public/static/js'),
        filename: 'journey.bundle.js',
    },
    module: {
        rules: [
            {
                test: /\.(ts|tsx)$/,
                exclude: [/node_modules/, /tests/],
                use: ['babel-loader', {loader: 'ts-loader', options: {onlyCompileBundledFiles: true}}],
            },
        ],
    },
};

module.exports = config;
