const client = require('./client');

describe('The save function', () => {
  test('returns 1', () => {
    expect(client.save()).toBe(1);
  });
});

describe('The get function', () => {
  test('returns 2', () => {
    expect(client.get()).toBe(2);
  });
});
