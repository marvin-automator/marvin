var webpack = require("webpack");
var CopyWebpackPlugin = require('copy-webpack-plugin');
var ExtractTextPlugin = require("extract-text-webpack-plugin");

process.env["NODE_ENV"] = "development"

module.exports = {
    devtool: "eval-source-map",
    entry: [
        "./assets/js/application.js",
    ],

    output: {
        filename: "application.js",
        path: __dirname + "/public/assets"
    },
    plugins: [
        new ExtractTextPlugin("application.css"),
        new CopyWebpackPlugin([{
            from: "./assets",
            to: ""
        }], {
            ignore: [
                "js/*",
            ]
        })
    ],
    resolve: {
        extensions: [".js", ".jsx", ".json"],
        modules: ["./assets/js", "./node_modules"]
    },
    module: {
        rules: [{
            test: /\.jsx?$/,
            loader: "babel-loader",
            options: {
                presets: ['react-app']
            },
            exclude: /node_modules/
        }, {
            test: /\.scss$/,
            use: ExtractTextPlugin.extract({
                fallback: "style-loader",
                use:
                    [{
                        loader: "css-loader",
                        options: { sourceMap: true }
                    },
                        {
                            loader: "sass-loader",
                            options: { sourceMap: true }
                        }]
            })
        },{
            test: /\.css$/,
            use: [
                { loader: "style-loader" },
                { loader: "css-loader" }
            ]
        }, {
            test: /\.woff(\?v=\d+\.\d+\.\d+)?$/,
            use: "url-loader?limit=10000&mimetype=application/font-woff"
        }, {
            test: /\.woff2(\?v=\d+\.\d+\.\d+)?$/,
            use: "url-loader?limit=10000&mimetype=application/font-woff"
        }, {
            test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/,
            use: "url-loader?limit=10000&mimetype=application/octet-stream"
        }, {
            test: /\.eot(\?v=\d+\.\d+\.\d+)?$/,
            use: "file-loader"
        }, {
            test: /\.svg(\?v=\d+\.\d+\.\d+)?$/,
            use: "url-loader?limit=10000&mimetype=image/svg+xml"
        }]
    }
};
