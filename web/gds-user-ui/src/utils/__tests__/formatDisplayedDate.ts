import formatDisplayedDate from 'utils/formatDisplayedDate';

describe('formatDisplayedDate', () => {
  it('should return N/A', () => {
    expect(formatDisplayedDate('')).toBe('N/A');
    expect(formatDisplayedDate('aaaa')).toBe('N/A');
  });

  it('should return formated date', () => {
    expect(formatDisplayedDate('2022-11-15T16:10:23.363Z')).toBe('15-11-2022');
    expect(formatDisplayedDate('2022-11-15T16:10:23.363Z', 'YYYY-MM-DD')).toBe('2022-11-15');
    expect(formatDisplayedDate('2022-11-15T16:10:23.363Z', 'YYYY/MM/DD')).toBe('2022/11/15');
    expect(formatDisplayedDate('2022-11-15T16:10:23.363Z', 'YYYY MMMM DD')).toBe(
      '2022 November 15'
    );
  });
});
