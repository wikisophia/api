import { terser } from 'rollup-plugin-terser';
import babel from '@rollup/plugin-babel';

const formats = ['iife', 'cjs', 'umd'];
const configs = [];


formats.forEach(format => {
  configs.push({
    plugins: [babel({
      exclude: "node_modules/**"
    }), terser()],
    input: ["src/client.js"],
    output: [
      {
        file: `${__dirname}/prod/${format}.js`,
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
        file: `${__dirname}/dev/${format}.js`,
        format: format,
        name: 'wksphArguments',
        sourcemap: true
      }
    ],
  });
});

export default configs;