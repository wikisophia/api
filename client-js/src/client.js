function handleServerErrors(response) {
  if (response.status >= 500 && response.status < 600) {
    throw new Error(`Server responded with ${response.status}: ${response.body}`);
  }
  return response;
}

function nullNotFound(response) {
  if (response.status === 404) {
    return null;
  }
  return response;
}

function getResponseBody(response) {
  if (response) {
    return response.body;
  }
  return response;
}

/**
 * Make sure the argument has everything it needs.
 * If valid, return null. If not, return an error message explaining what's wrong with it.
 *
 * @param {Argument} argument The argument to validate
 * @return {string|null} Null if the argument is valid, or a string error message otherwise.
 */
function validate(argument) {
  if (argument.length < 2) {
    return 'An argument must have at least two premises.';
  }

  if (!argument.conclusion) {
    return 'An argument must have a conclusion.';
  }

  // Sets aren't supported in a few semi-modern browsers
  const premiseSet = {};
  let duplicate = null;
  argument.premises.forEach((premise) => {
    if (premiseSet[premise]) {
      duplicate = `Arguments shouldn't use the same premise more than once. Yours repeats: ${premiseSet[premise]}`;
    }
    premiseSet[premise] = true;
  });
  return duplicate;
}

/**
 * Make a new client.
 *
 * @param {ClientArguments} cfg arguments to configure the client
 *
 * @see https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API
 * @see https://www.npmjs.com/package/node-fetch
 */
export default function newClient({ url, fetch }) {
  return {

    /**
     * Get a specific argument.
     *
     * @param {int} id The ID of the argument to get.
     * @param [int] version The version to get. If undefined, the latest version will be fetched.
     *
     * @return {Promise<Argument>} The Argument, if found, or null if not.
     *   The Promise will reject on network or server errors.
     */
    getOne(id, version) {
      let getURL = `${url}/arguments/${id}`;
      if (version > 0) {
        getURL = `${getURL}/version/${version}`;
      }

      return fetch(getURL, {
        mode: 'cors',
      }).then(handleServerErrors)
        .then(nullNotFound)
        .then(getResponseBody);
    },

    /**
     * Save a new argument.
     *
     * @param {Argument} argument The argument to be saved
     * @return {Promise<SaveResponse>} A Promise with info describing where to
     *   find the new argument.
     */
    save(argument) {
      const err = validate(argument);
      if (err) {
        return Promise.reject(new Error(err));
      }

      return fetch(`${url}/arguments`, {
        method: 'POST',
        mode: 'cors',
        body: argument,
      }).then(handleServerErrors)
        .then(response => ({
          location: response.headers.get('Location'),
        }));
    },
  };
}

/**
 * @typedef {Object} ClientArguments
 *
 * @property {string} url The URL to the server hosting the Arguments API.
 *   For example, "https://arguments.wikisophia.net".
 * @property {function} fetch A function implementing the Fetch API.
 *   For browsers, this can just be the global built-in "fetch".
 *   For node, you'll need to use something like node-fetch.
 */

/**
 * @typedef {Object} Arguments
 *
 * @property {string} conclusion The argument's conclusion.
 * @property {string[]} premises The argument's premises.
 *   This must have at least 2 elements for the argument to be valid..
 */

/**
 * @typedef {Object} SaveResponse
 *
 * @property {string} location A URL where the saved argument can be found.
 */
