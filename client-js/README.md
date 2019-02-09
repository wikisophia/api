# Arguments API Javascript Client

A javascript client library for the Wikisophia Arguments API.

# Usage

```javascript
import newArgumentsClient from 'api-arguments-client';

const arguments = newArgumentsClient({
    url: "www.wikisophia.net",
    fetch: fetch,
});

arguments.save({
    premises: [
        "Wood floats in water.",
        "Ducks float in water.",
        "Objects which are the same weight will either both sink or both float in water.",
        "Objects with the same weight are made of the same material.",
        "She weighs the same as a duck.",
    ],
    conclusion: "She's made of wood.",
});
```

This library works from both the client and server, provided you can send it a
[fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API) function.
In environments that don't support this natively, you may need to
[Polyfill](https://developer.mozilla.org/en-US/docs/Glossary/Polyfill) it.

See:
- [Node.js fetch polyfill](https://github.com/bitinn/node-fetch)
- [Browser fetch polyfill](https://github.com/github/fetch)
