module.exports = {
    extends: "airbnb-base",
    rules: {
        "no-console": ['error', {'allow': ['warn', 'error']}],
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
