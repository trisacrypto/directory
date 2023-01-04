import { defaultEndpointPrefix } from '..';

describe('defaultEndpointPrefix', () => {
  it('should return the correct dev url prefix', () => {
    const devUrlPrefix = 'http://localhost:4434/v2';

    expect(defaultEndpointPrefix()).toBe(devUrlPrefix);
  });
});
