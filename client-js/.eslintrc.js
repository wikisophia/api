module.exports = {
    extends: 'airbnb-base',
    rules: {
        'space-before-function-paren': 'off',
        'no-use-before-define': 'off'
    },
    overrides: [
        {
            files: [
                '**/*.test.js'
            ],
            env: {
                jest: true
            }
        }
    ]
}
