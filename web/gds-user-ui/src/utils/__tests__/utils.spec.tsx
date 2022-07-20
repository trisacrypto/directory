import {
  currencyFormatter,
  findStepKey,
  getColorScheme,
  getDomain,
  getStepStatus,
  getValueByPathname,
  hasStepError,
  hasValue,
  isValidUuid
} from 'utils/utils';
import { faker } from '@faker-js/faker';

const steps = [
  { key: 1, status: 'complete', data: {} },
  { key: 2, status: 'complete' },
  { key: 3, status: 'complete' },
  { key: 4, status: 'complete' },
  { key: 5, status: 'complete' },
  { key: 6, status: 'complete' }
];

describe('utils', () => {
  describe('isValidUuid', () => {
    it('should be a valid uuid', () => {
      expect(isValidUuid(faker.datatype.uuid())).toBe(true);
    });

    it('should not be a valid uuid', () => {
      expect(isValidUuid(faker.lorem.text())).toBe(false);
    });
  });

  describe('getValueByPathname', () => {
    it('should return the right property', () => {
      const data = {
        users: [
          {
            id: 1,
            name: 'Leanne Graham',
            username: 'Bret',
            email: 'Sincere@april.biz',
            address: {
              street: 'Kulas Light',
              suite: 'Apt. 556',
              city: 'Gwenborough',
              zipcode: '92998-3874',
              geo: {
                lat: '-37.3159',
                lng: '81.1496'
              }
            },
            phone: '1-770-736-8031 x56442',
            website: 'hildegard.org',
            company: {
              name: 'Romaguera-Crona',
              catchPhrase: 'Multi-layered client-server neural-net',
              bs: 'harness real-time e-markets'
            }
          }
        ]
      };
      expect(getValueByPathname(data, 'users[0].address.street')).toBe(
        data.users[0].address.street
      );
      expect(getValueByPathname(data, 'users[0].company.name')).toBe(data.users[0].company.name);
    });

    it('should return undefined when the path is incorrect', () => {
      const data = {
        users: [
          {
            id: 1,
            name: 'Leanne Graham',
            username: 'Bret',
            email: 'Sincere@april.biz',
            address: {
              street: 'Kulas Light',
              suite: 'Apt. 556',
              city: 'Gwenborough',
              zipcode: '92998-3874',
              geo: {
                lat: '-37.3159',
                lng: '81.1496'
              }
            },
            phone: '1-770-736-8031 x56442',
            website: 'hildegard.org',
            company: {
              name: 'Romaguera-Crona',
              catchPhrase: 'Multi-layered client-server neural-net',
              bs: 'harness real-time e-markets'
            }
          }
        ]
      };
      expect(getValueByPathname(data, 'users[1].address.street')).toBeUndefined();
    });
  });

  describe('getDomain', () => {
    it('should return the right domain', () => {
      expect(getDomain('https://optimistic-rabbi.org')).toBe('optimistic-rabbi.org');
    });

    it('should return null when URL is not correct', () => {
      expect(getDomain(faker.random.words())).toBeNull();
    });
  });

  describe('hasValue', () => {
    it('should return false when some of properties value are truthy', () => {
      const data = {
        userId: 1,
        id: 1,
        title: 'delectus aut autem',
        completed: false
      };
      expect(hasValue(data)).toBe(true);
    });

    it('should return false when some of properties value are falsy', () => {
      const data = {
        userId: 0,
        id: 0,
        title: '',
        completed: false
      };
      expect(hasValue(data)).toBe(false);
    });
  });

  describe('getColorScheme', () => {
    it('should return cyan when status is yes', () => {
      expect(getColorScheme('yes')).toBe('cyan');
    });

    it('should return cyan when status is true', () => {
      expect(getColorScheme('true')).toBe('cyan');
    });

    it('should return cyan when status is true', () => {
      expect(getColorScheme('false')).toBe('cyan');
    });
  });

  describe('currentFormater', () => {
    it('should return a formatted amount in dollar', () => {
      expect(currencyFormatter(5_000)).toBe(`$5,000.00`);
      expect(currencyFormatter(5_000_000)).toBe(`$5,000,000.00`);
    });
  });

  describe('findStepKey', () => {
    it('should return the right step key', () => {
      expect(findStepKey(steps, 5)).toEqual([{ key: 5, status: 'complete' }]);
    });

    it('should return an empty array', () => {
      expect(findStepKey(steps, 7)).toEqual([]);
    });
  });

  describe('getStepStatus', () => {
    it('should return the right step status', () => {
      expect(getStepStatus(steps, 5)).toBe('complete');
    });

    it('should return undefined if step doesnt exist', () => {
      expect(getStepStatus(steps, 9)).toBeUndefined();
    });
  });

  describe('hasStepError', () => {
    it('should return true when some of statuses are equal to error', () => {
      const _steps = [
        { key: 1, status: 'error', data: {} },
        { key: 2, status: 'complete' },
        { key: 3, status: 'error' },
        { key: 4, status: 'complete' }
      ];

      expect(hasStepError(_steps)).toBe(true);
    });

    it('should return false when statuses are different to error', () => {
      const _steps = [
        { key: 1, status: 'complete', data: {} },
        { key: 2, status: 'complete' },
        { key: 3, status: 'complete' },
        { key: 4, status: 'complete' }
      ];

      expect(hasStepError(_steps)).toBe(false);
    });
  });
});
