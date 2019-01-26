module.exports = {
    extends: "airbnb-base",
    rules: {
        "space-before-function-paren": 'off'
    },
    overrides: {
        files: [
            '**/*.test.js'
        ],
        env: {
            jest: true
        }
    }
}
