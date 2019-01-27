import newClient from './client';

const url = 'http://some-url.com';

describe('The getOne function', () => {
  test('calls the right API endpoint when fetching the latest version of an argument', () => {
    const fetch = jest.fn();
    fetch.mockReturnValueOnce(Promise.resolve({
      status: 404,
    }));

    const client = newClient({ url, fetch });
    return client.getOne(1).then(() => {
      expect(fetch.mock.calls.length).toBe(1);
      expect(fetch.mock.calls[0][0]).toBe(`${url}/arguments/1`);
    });
  });

  test('calls the right API endpoint when fetching a specific version of an argument', () => {
    const fetch = jest.fn();
    fetch.mockReturnValueOnce(Promise.resolve({
      status: 404,
    }));

    const client = newClient({ url, fetch });
    return client.getOne(1, 2).then(() => {
      expect(fetch.mock.calls.length).toBe(1);
      expect(fetch.mock.calls[0][0]).toBe(`${url}/arguments/1/version/2`);
    });
  });

  test('parses the response body successfully', () => {
    const badArgument = {
      conclusion: 'one thing',
      premises: ["don't relate", 'at all'],
    };
    const fetch = jest.fn();
    fetch.mockReturnValueOnce(Promise.resolve({
      status: 200,
      body: badArgument,
    }));

    const client = newClient({ url, fetch });
    return expect(client.getOne(1)).resolves.toEqual(badArgument);
  });

  test('resolves to null null when the server responds with a 404', () => {
    const fetch = jest.fn();
    fetch.mockReturnValueOnce(Promise.resolve({
      status: 404,
    }));
    const client = newClient({ url, fetch });
    return expect(client.getOne(1)).resolves.toBe(null);
  });

  test('rejects if the server responds with a 500', () => {
    const fetch = jest.fn();
    fetch.mockReturnValueOnce(Promise.resolve({
      status: 500,
      body: 'something bad happened',
    }));
    const client = newClient({ url, fetch });
    return expect(client.getOne(1)).rejects.toEqual(new Error('Server responded with 500: something bad happened'));
  });
});
