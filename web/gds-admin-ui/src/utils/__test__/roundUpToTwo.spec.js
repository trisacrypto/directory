import roundUpToTwo from '@/utils/roundUptoTwo';

describe('roundUpToTwo', () => {
  it('should return correct rounded value', () => {
    expect(roundUpToTwo(7.000000006)).toBe(7);
    expect(roundUpToTwo(1.000000000001)).toBe(1);
    expect(roundUpToTwo(0.000000000001)).toBe(0);
    expect(roundUpToTwo(0.01)).toBe(0.01);
    expect(roundUpToTwo(1.01)).toBe(1.01);
    expect(roundUpToTwo(7.688888)).toBe(7.69);
    expect(roundUpToTwo(7.64555)).toBe(7.65);
    expect(roundUpToTwo(7.60555)).toBe(7.61);
    expect(roundUpToTwo(0.64555)).toBe(0.65);
    expect(roundUpToTwo(0.64355)).toBe(0.64);
    expect(roundUpToTwo(1.1)).toBe(1.1);
  });

  it('should throw an error when passed value is not a number', () => {
    expect(roundUpToTwo('test')).toBeUndefined();
  });
});
