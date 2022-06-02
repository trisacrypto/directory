import countryCodeEmoji, { getCountryName } from 'utils/country';

describe('country', () => {
  describe('emojiCountryCode', () => {
    it('should return the country name', () => {
      expect(getCountryName('US')).toBe('United States');
      expect(getCountryName('FR')).toBe('France');
    });
  });

  describe('countryCodeEmoji', () => {
    it('should the right country emoji', () => {
      expect(countryCodeEmoji('US')).toBe('🇺🇸');
      expect(countryCodeEmoji('FR')).toBe('🇫🇷');
    });
  });
});
