import newClient from './client';

const url = 'http://some-url.com';

describe('getOne()', () => {
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
    return expect(client.getOne(1)).rejects.toThrow('Server responded with a 500: something bad happened');
  });
});

describe('save()', () => {
  test('calls the right API endpoint with the right data', () => {
    const fetch = jest.fn();
    fetch.mockReturnValueOnce(Promise.resolve({
      status: 201,
      headers: {
        get(header) {
          return header === 'Location' ? '/arguments/1' : '';
        },
      },
    }));
    const client = newClient({ url, fetch });
    const argument = {
      premises: ['foo', 'bar'],
      conclusion: 'baz',
    };
    return client.save(argument).then((resolved) => {
      expect(fetch.mock.calls.length).toBe(1);
      expect(fetch.mock.calls[0][0]).toBe(`${url}/arguments`);
      expect(fetch.mock.calls[0][1]).toEqual({
        method: 'POST',
        mode: 'cors',
        body: argument,
      });
      expect(resolved).toEqual({
        location: '/arguments/1',
      });
    });
  });

  test('rejects if the server returns a 500', () => {
    const fetch = jest.fn();
    fetch.mockReturnValueOnce(Promise.resolve({
      status: 500,
      body: 'Something went wrong',
    }));
    const client = newClient({ url, fetch });

    return expect(client.save({
      conclusion: 'baz',
      premises: ['foo', 'bar'],
    })).rejects.toThrow('Server responded with a 500: Something went wrong');
  });

  test('rejects arguments with no conclusion', () => {
    const fetch = jest.fn();
    const client = newClient({ url, fetch });

    return expect(client.save({
      premises: ['foo', 'bar'],
    })).rejects.toThrow('An argument must have a conclusion.');
  });

  test('rejects arguments with duplicate premises', () => {
    const fetch = jest.fn();
    const client = newClient({ url, fetch });

    return expect(client.save({
      conclusion: 'bar',
      premises: ['foo', 'foo'],
    })).rejects.toThrow("Arguments shouldn't use the same premise more than once. Yours repeats: foo");
  });

  test('rejects arguments with too few premises', () => {
    const fetch = jest.fn();
    const client = newClient({ url, fetch });

    return expect(client.save({
      conclusion: 'bar',
      premises: ['foo'],
    })).rejects.toThrow('An argument must have at least two premises.');
  });
});
