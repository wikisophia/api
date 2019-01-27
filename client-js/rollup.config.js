import { uglify } from 'rollup-plugin-uglify';
import babel from 'rollup-plugin-babel';

const formats = ['iife', 'umd'];
const configs = [];


formats.forEach(format => {
  configs.push({
    plugins: [babel({
      exclude: "node_modules/**"
    }), uglify()],
    input: ["src/client.js"],
    output: [
      {
        file: `${__dirname}/dist/${format}.min.js`,
        format: format,
        name: 'wksphArguments'
      }
    ]
  });
  configs.push({
    plugins: [babel({
      exclude: "node_modules/**"
    })],
    input: ["src/client.js"],
    output: [
      {
        file: `${__dirname}/dist/${format}.dev.js`,
        format: format,
        name: 'wksphArguments',
        sourcemap: true
      }
    ],
  });
});

export default configs;