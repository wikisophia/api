import { uglify } from 'rollup-plugin-uglify';

const formats = ['iife', 'umd'];
const configs = [];

formats.forEach(format => {
  configs.push({
    plugins: [uglify()],
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