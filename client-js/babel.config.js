module.exports = function(api) {
  api.cache(false);

  if (process.env.BABEL_ENV === 'test') {
    return {
      presets: [
        [
          '@babel/env',
          {
            targets: {
              node: 'current',
            },
          }
        ]
      ]
    };
  }

  return {
    presets: [
      [
        '@babel/env'
      ]
    ]
  };
};
