import { Response } from 'node-fetch';
import newClient from './client';
import getOneResponseSample from '../../server/samples/get-one-response.json';
import getAllResponseSample from '../../server/samples/get-all-response.json';
import saveRequest from '../../server/samples/save-request.json';
import updateRequest from '../../server/samples/update-request.json';

const url = 'http://some-url.com';
const getOneResponse = {
  argument: getOneResponseSample,
};
const getAllResponse = {
  arguments: getAllResponseSample,
};

function mockOneReturn(status, body, headers) {
  const mock = jest.fn();
  mock.mockReturnValueOnce(Promise.resolve(new Response(body, {
    status,
    headers,
  })));

  return mock;
}

function saveRequestToResponse(req) {
  return {
    argument: Object.assign({}, saveRequest, {id: 1, version: 1})
  }
}

describe('getOne()', () => {
  test('calls the right API endpoint when fetching the latest version of an argument', () => {
    const fetch = mockOneReturn(200, JSON.stringify(getOneResponse), {
      'Content-Type': 'application/json; charset=utf-8',
    });

    const client = newClient({ url, fetch });
    return client.getOne(1).then((result) => {
      expect(fetch.mock.calls.length).toBe(1);
      expect(fetch.mock.calls[0][0]).toBe(`${url}/arguments/1`);
      expect(result).toEqual(getOneResponse);
    });
  });

  test('calls the right API endpoint when fetching a specific version of an argument', () => {
    const fetch = mockOneReturn(404, 'not found');

    const client = newClient({ url, fetch });
    return client.getOne(1, 2).then((result) => {
      expect(fetch.mock.calls.length).toBe(1);
      expect(fetch.mock.calls[0][0]).toBe(`${url}/arguments/1/version/2`);
      expect(result).toBe(null);
    });
  });

  test('rejects if the server responds with a 500', () => {
    const fetch = mockOneReturn(500, 'something bad happened');
    const client = newClient({ url, fetch });
    return expect(client.getOne(1)).rejects.toThrow('The server responded with a 500: something bad happened');
  });
});

describe('getAll()', () => {
  test('calls the right API endpoint', () => {
    const fetch = mockOneReturn(200, JSON.stringify(getAllResponse), {
      'Content-Type': 'application/json; charset=utf-8',
    });

    const client = newClient({ url, fetch });
    const { conclusion } = getAllResponseSample.arguments[0];
    return client.getAll(conclusion).then((result) => {
      expect(fetch.mock.calls.length).toBe(1);
      expect(fetch.mock.calls[0][0]).toBe(`${url}/arguments?conclusion=${conclusion}`);
      expect(fetch.mock.calls[0][1]).toEqual({ mode: 'cors' });
      expect(result).toEqual(getAllResponse);
    });
  });

  test('resolves to an empty array when the server responds with a 404', () => {
    const fetch = mockOneReturn(404, '');
    const client = newClient({ url, fetch });
    return expect(client.getAll('nothing ventured, nothing earned')).resolves.toEqual({ arguments: [] });
  });

  test('rejects if the server returns a 500', () => {
    const fetch = mockOneReturn(500, 'server failure');
    const client = newClient({ url, fetch });
    return expect(client.getAll("don't mess with texas")).rejects.toThrow('The server responded with a 500: server failure');
  });

  test('rejects if called with an empty conclusion', () => {
    const fetch = jest.fn();
    const client = newClient({ url, fetch });
    const gotten = client.getAll();
    expect(fetch.mock.calls.length).toBe(0);
    return expect(gotten).rejects.toThrow("Can't get arguments with an empty conclusion.");
  });
});

describe('save()', () => {
  test('calls the right API endpoint with the right data', () => {
    const response = saveRequestToResponse(saveRequest);
    const mockLocation = '/arguments/1/version/1';
    const fetch = mockOneReturn(201, JSON.stringify(response), {
      Location: mockLocation,
    });
    const client = newClient({ url, fetch });
    return client.save(saveRequest).then((resolved) => {
      expect(fetch.mock.calls.length).toBe(1);
      expect(fetch.mock.calls[0][0]).toBe(`${url}/arguments`);
      expect(fetch.mock.calls[0][1].method).toEqual('POST');
      expect(fetch.mock.calls[0][1].mode).toEqual('cors');
      expect(JSON.parse(fetch.mock.calls[0][1].body)).toEqual(saveRequest);
      expect(resolved).toEqual(Object.assign(saveRequestToResponse(saveRequest), {
        location: mockLocation
      }));
    });
  });

  test('rejects if the server returns a 500', () => {
    const fetch = mockOneReturn(500, 'Something went wrong');
    const client = newClient({ url, fetch });
    return expect(client.save(saveRequest)).rejects.toThrow('The server responded with a 500: Something went wrong');
  });

  test('rejects arguments with no conclusion', () => {
    const fetch = jest.fn();
    const client = newClient({ url, fetch });
    const mangled = Object.assign({}, saveRequest);
    delete mangled.conclusion;
    const saved = client.save(mangled);
    expect(fetch.mock.calls.length).toBe(0);
    return expect(saved).rejects.toThrow('An argument must have a conclusion.');
  });

  test('rejects arguments with duplicate premises', () => {
    const fetch = jest.fn();
    const client = newClient({ url, fetch });
    const mangled = Object.assign({}, saveRequest);
    mangled.premises = Array(saveRequest.premises.length).fill(saveRequest.premises[0]);
    const saved = client.save(mangled);
    expect(fetch.mock.calls.length).toBe(0);
    return expect(saved).rejects.toThrow(`Arguments shouldn't use the same premise more than once. Yours repeats: ${saveRequest.premises[0]}`);
  });

  test('rejects arguments with too few premises', () => {
    const fetch = jest.fn();
    const client = newClient({ url, fetch });
    const mangled = Object.assign({}, saveRequest);
    mangled.premises = ['only one'];
    const saved = client.save(mangled);
    expect(fetch.mock.calls.length).toBe(0);
    return expect(saved).rejects.toThrow('An argument must have at least two premises.');
  });
});

describe('update()', () => {
  test('calls the right API endpoint with the right data', () => {
    const response = saveRequestToResponse(saveRequest);
    const mockLocation = '/arguments/1/version/2';
    const fetch = mockOneReturn(200, JSON.stringify(response), {
      Location: mockLocation,
    });
    const client = newClient({ url, fetch });
    return client.update(1, updateRequest).then((resolved) => {
      expect(fetch.mock.calls.length).toBe(1);
      expect(fetch.mock.calls[0][0]).toBe(`${url}/arguments/1`);
      expect(fetch.mock.calls[0][1].method).toEqual('PATCH');
      expect(fetch.mock.calls[0][1].mode).toEqual('cors');
      expect(JSON.parse(fetch.mock.calls[0][1].body)).toEqual(updateRequest);
      expect(resolved).toEqual(Object.assign(saveRequestToResponse(saveRequest), {
        location: mockLocation
      }));
    });
  });

  test('allows updates of conclusions only', () => {
    const response = saveRequestToResponse(saveRequest);
    const mockLocation = '/arguments/1/version/2';
    const fetch = mockOneReturn(200, JSON.stringify(response), {
      Location: mockLocation,
    });
    const client = newClient({ url, fetch });
    const updatePremisesOnly = Object.assign({}, updateRequest, { conclusion: null });
    return expect(client.update(1, updatePremisesOnly)).resolves.toEqual({
      argument: response.argument,
      location: mockLocation
    });
  })

  test('rejects if the server returns a 500', () => {
    const fetch = mockOneReturn(500, 'Something went wrong');
    const client = newClient({ url, fetch });
    return expect(client.update(1, updateRequest)).rejects.toThrow('The server responded with a 500: Something went wrong');
  });

  test('rejects if the server returns a 404', () => {
    const fetch = mockOneReturn(404, "Argument with id=1 doesn't exist");
    const client = newClient({ url, fetch });
    return expect(client.update(1, updateRequest)).rejects.toThrow("The server returned a 404: Argument with id=1 doesn't exist.");
  });

  test('rejects updates with duplicate premises', () => {
    const fetch = jest.fn();
    const client = newClient({ url, fetch });
    const update = Object.assign({}, updateRequest, {
      premises: Array(updateRequest.premises.length).fill(updateRequest.premises[0]),
    });
    const call = client.update(1, update);
    expect(fetch.mock.calls.length).toBe(0);
    return expect(call).rejects.toThrow(`Arguments shouldn't use the same premise more than once. Yours repeats: ${updateRequest.premises[0]}`);
  });

  test('rejects updates with too few premises', () => {
    const fetch = jest.fn();
    const client = newClient({ url, fetch });
    const update = Object.assign({}, updateRequest, {
      premises: ['only one'],
    });
    const call = client.update(1, update);
    expect(fetch.mock.calls.length).toBe(0);
    return expect(call).rejects.toThrow('An argument must have at least two premises.');
  });

  test('rejects empty updates', () => {
    const fetch = jest.fn();
    const client = newClient({ url, fetch });
    const update = {};
    const call = client.update(1, update);
    expect(fetch.mock.calls.length).toBe(0);
    return expect(call).rejects.toThrow('Updates must change premises, a conclusion, or both.');
  });
});
