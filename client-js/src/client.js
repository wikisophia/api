function handleServerErrors(response) {
  if (response.status >= 500 && response.status < 600) {
    throw new Error(`The server responded with a ${response.status}: ${response.body}`);
  }
  return response;
}

function onNotFound(value) {
  return function responseOrValue(response) {
    if (response.status === 404) {
      return value;
    }
    return response;
  };
}

function parseJSONResponseBody(response) {
  if (response && response.json) {
    return response.json();
  }
  return response;
}

/**
 * Make sure the premises would be valid.
 * If valid, return null. If not, return an error message explaining what's wrong with it.
 *
 * @param {string[]} premises The premises to validate
 * @return {string|null} Null if the premises is valid, or a string error message otherwise.
 */
function validatePremises(premises) {
  if (premises.length < 2) {
    return 'An argument must have at least two premises.';
  }

  // Sets aren't supported in a few semi-modern browsers
  const premiseSet = {};
  let duplicate = null;
  premises.forEach((premise) => {
    if (premiseSet[premise]) {
      duplicate = `Arguments shouldn't use the same premise more than once. Yours repeats: ${premise}`;
    }
    premiseSet[premise] = true;
  });
  return duplicate;
}

/**
 * Make sure the argument has everything it needs.
 * If valid, return null. If not, return an error message explaining what's wrong with it.
 *
 * @param {Argument} argument The argument to validate
 * @return {string|null} Null if the argument is valid, or a string error message otherwise.
 */
function validateArgument(argument) {
  if (!argument.conclusion) {
    return 'An argument must have a conclusion.';
  }

  return validatePremises(argument.premises);
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
     * @return {Promise<OneArgument>} The Argument, if found, or null if not.
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
        .then(onNotFound(null))
        .then(parseJSONResponseBody);
    },

    /**
     * Get all the arguments with a given conclusion.
     *
     * @param {string} conclusion The conclusion you want to fetch all arguments for.
     * @return {Promise<SomeArguments>} A list of arguments with this conclusion.
     *   If none exist, this will be an empty array.
     */
    getAll(conclusion) {
      if (!conclusion) {
        return Promise.reject(new Error("Can't get arguments with an empty conclusion."));
      }

      return fetch(`${url}/arguments?conclusion=${conclusion}`, {
        mode: 'cors',
      }).then(handleServerErrors)
        .then(onNotFound({ arguments: [] }))
        .then(parseJSONResponseBody);
    },

    /**
     * Save a new argument.
     *
     * @param {Argument} argument The argument to be saved
     * @return {Promise<SaveResponse>} A Promise with info describing where to
     *   find the new argument.
     */
    save(argument) {
      const err = validateArgument(argument);
      if (err) {
        return Promise.reject(new Error(err));
      }

      return fetch(`${url}/arguments`, {
        method: 'POST',
        mode: 'cors',
        body: JSON.stringify(argument),
      }).then(handleServerErrors)
        .then(response => ({
          location: response.headers.get('Location'),
        }));
    },

    /**
     * Update the premises of the argument using its ID.
     *
     * @param {int} id The ID of the argument you want to update.
     * @param {Argument} argument The new argument which should have this ID.
     * @return {Promise<SaveResponse>} A Promise with info describing where to
     *   find the new argument.
     */
    update(id, argument) {
      const err = validateArgument(argument);
      if (err) {
        return Promise.reject(new Error(err));
      }

      return fetch(`${url}/arguments/${id}`, {
        method: 'PATCH',
        mode: 'cors',
        body: JSON.stringify(argument),
      }).then(handleServerErrors)
        .then((response) => {
          if (response.status === 404) {
            return new Promise(((resolve, reject) => {
              response.text().then((responseBody) => {
                reject(new Error(`The server returned a 404: ${responseBody}.`));
              }).catch((readErr) => {
                reject(new Error(`The server returned a 404, and an error occurred while reading the response body: ${readErr.message}.`));
              });
            }));
          }
          return response;
        })
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
 * @typedef {Object} Argument
 *
 * @property {string} conclusion The argument's conclusion.
 * @property {string[]} premises The argument's premises.
 *   This must have at least 2 elements for the argument to be valid.
 */

/**
 * @typedef {Object} ArgumentResponse
 *
 * @property {int} id The argument's ID.
 * @property {int} version The argument's version.
 * @property {string} conclusion The argument's conclusion.
 * @property {string[]} premises The argument's premises.
 *   This must have at least 2 elements for the argument to be valid.
 */

/**
 * @typedef {Object} SomeArguments
 *
 * @property {ArgumentResponse[]} arguments The list of arguments.
*/

/**
 * @typedef {Object} OneArgument
 *
 * @property {ArgumentResponse} argument The argument.
 */

/**
 * @typedef {Object} SaveResponse
 *
 * @property {string} location A URL where the saved argument can be found.
 * @property {ArgumentResponse} argument The argument after
 */
