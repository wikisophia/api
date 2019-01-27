# Arguments API Javascript Client

A javascript client library for the Arguments API.

# Usage

```javascript
import newArgumentsClient from 'api-arguments-client';

const arguments = newArgumentsClient({
    url: "www.wikisophia.net",
    fetch: fetch,
});

arguments.save({
    premises: [
        "Things which weigh the same are made of the same material.",
        "She weighs the same as a duck.",
        "Wood floats in water.",
        "Ducks float in water.",
        "Things which are the same weight will either both sink or both float in water.",
    ],
    conclusion: "She's made of wood.",
});
```

The [fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API) function is a global on most
modern browsers. In Node.js and older browsers, you may need to
[Polyfill](https://developer.mozilla.org/en-US/docs/Glossary/Polyfill) it.

See also:
- [Node.js fetch polyfill](https://github.com/bitinn/node-fetch)
- [Browser fetch polyfill](https://github.com/github/fetch)
