import newClient from './client';
import getOneResponse from '../../server/samples/get-one-response.json';
import getAllResponse from '../../server/samples/get-all-response.json';
import saveRequest from '../../server/samples/save-request.json';
import updateRequest from '../../server/samples/update-request.json';

const url = 'http://some-url.com';

describe('getOne()', () => {
  test('calls the right API endpoint when fetching the latest version of an argument', () => {
    const fetch = jest.fn();
    fetch.mockReturnValueOnce(Promise.resolve({
      status: 200,
      body: getOneResponse,
    }));

    const client = newClient({ url, fetch });
    return client.getOne(1).then((result) => {
      expect(fetch.mock.calls.length).toBe(1);
      expect(fetch.mock.calls[0][0]).toBe(`${url}/arguments/1`);
      expect(result).toEqual(getOneResponse);
    });
  });

  test('calls the right API endpoint when fetching a specific version of an argument', () => {
    const fetch = jest.fn();
    fetch.mockReturnValueOnce(Promise.resolve({
      status: 404,
    }));

    const client = newClient({ url, fetch });
    return client.getOne(1, 2).then((result) => {
      expect(fetch.mock.calls.length).toBe(1);
      expect(fetch.mock.calls[0][0]).toBe(`${url}/arguments/1/version/2`);
      expect(result).toBe(null);
    });
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

describe('getAll()', () => {
  test('calls the right API endpoint', () => {
    const fetch = jest.fn();
    fetch.mockReturnValueOnce(Promise.resolve({
      status: 200,
      body: getAllResponse,
    }));

    const client = newClient({ url, fetch });
    return client.getAll(getAllResponse[0].conclusion).catch((result) => {
      expect(fetch.mock.calls.length).toBe(1);
      expect(fetch.mock.calls[0][0]).toBe(`${url}/arguments?${getAllResponse[0].conclusion}`);
      expect(fetch.mock.calls[0][1]).toEqual({ method: 'cors' });
      expect(result).toEqual(getAllResponse);
    });
  });

  test('resolves to an empty array when the server responds with a 404', () => {
    const fetch = jest.fn();
    fetch.mockReturnValueOnce(Promise.resolve({
      status: 404,
    }));

    const client = newClient({ url, fetch });
    return expect(client.getAll('nothing ventured, nothing earned')).resolves.toEqual([]);
  });

  test('rejects if the server returns a 500', () => {
    const fetch = jest.fn();
    fetch.mockReturnValueOnce(Promise.resolve({
      status: 500,
      body: 'server failure',
    }));

    const client = newClient({ url, fetch });
    return expect(client.getAll("don't mess with texas")).rejects.toThrow('Server responded with a 500: server failure');
  });

  test('rejects if called with an empty conclusion', () => {
    const fetch = jest.fn();
    const client = newClient({ url, fetch });
    return expect(client.getAll()).rejects.toThrow("Can't get arguments with an empty conclusion.");
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
    return client.save(saveRequest).then((resolved) => {
      expect(fetch.mock.calls.length).toBe(1);
      expect(fetch.mock.calls[0][0]).toBe(`${url}/arguments`);
      expect(fetch.mock.calls[0][1]).toEqual({
        method: 'POST',
        mode: 'cors',
        body: saveRequest,
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

    return expect(client.save(saveRequest)).rejects.toThrow('Server responded with a 500: Something went wrong');
  });

  test('rejects arguments with no conclusion', () => {
    const fetch = jest.fn();
    const client = newClient({ url, fetch });
    const mangled = Object.assign({}, saveRequest);
    delete mangled.conclusion;
    return expect(client.save(mangled)).rejects.toThrow('An argument must have a conclusion.');
  });

  test('rejects arguments with duplicate premises', () => {
    const fetch = jest.fn();
    const client = newClient({ url, fetch });
    const mangled = Object.assign({}, saveRequest);
    mangled.premises = Array(saveRequest.premises.length).fill(saveRequest.premises[0]);
    return expect(client.save(mangled)).rejects.toThrow(`Arguments shouldn't use the same premise more than once. Yours repeats: ${saveRequest.premises[0]}`);
  });

  test('rejects arguments with too few premises', () => {
    const fetch = jest.fn();
    const client = newClient({ url, fetch });
    const mangled = Object.assign({}, saveRequest);
    mangled.premises = ['only one'];
    return expect(client.save(mangled)).rejects.toThrow('An argument must have at least two premises.');
  });
});

describe('update()', () => {
  test('calls the right API endpoint with the right data', () => {
    const fetch = jest.fn();
    fetch.mockReturnValueOnce(Promise.resolve({
      status: 204,
      headers: {
        get(header) {
          return header === 'Location' ? '/arguments/1/version/2' : '';
        },
      },
    }));
    const client = newClient({ url, fetch });
    return client.update(1, updateRequest.premises).then((resolved) => {
      expect(fetch.mock.calls.length).toBe(1);
      expect(fetch.mock.calls[0][0]).toBe(`${url}/argument/1`);
      expect(fetch.mock.calls[0][1]).toEqual({
        method: 'PATCH',
        mode: 'cors',
        body: updateRequest,
      });
      expect(resolved).toEqual({
        location: '/arguments/1/version/2',
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
    return expect(client.update(1, updateRequest.premises)).rejects.toThrow('Something went wrong');
  });

  test('rejects updates with duplicate premises', () => {
    const fetch = jest.fn();
    const client = newClient({ url, fetch });
    const update = Array(updateRequest.premises.length).fill(updateRequest.premises[0]);
    return expect(client.update(1, update)).rejects.toThrow(`Arguments shouldn't use the same premise more than once. Yours repeats: ${updateRequest.premises[0]}`);
  });

  test('rejects updates with too few premises', () => {
    const fetch = jest.fn();
    const client = newClient({ url, fetch });
    return expect(client.update(1, ['only one'])).rejects.toThrow('An argument must have at least two premises.');
  });
});
