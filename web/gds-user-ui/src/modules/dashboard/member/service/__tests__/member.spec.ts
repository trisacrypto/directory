import axios from 'axios';
// import mockedAxios from 'jest-mock-axios';
// // import MockAdapter from "axios-mock-adapter";
// import axiosInstance from 'utils/axios';

import { getMembersService } from '..';
import { mainnetMembersMockValue, testnetMembersMockValue } from '../../__mocks__';
// const mock = new MockAdapter(axios);
// mock.onGet('/members').reply(200, mainnetMembersMockValue);
jest.mock('axios', () => {
  return {
    create: () => {
      return {
        get: jest.fn(),
        post: jest.fn(),
        put: jest.fn(),
        delete: jest.fn(),
        interceptors: {
          request: { eject: jest.fn(), use: jest.fn() },
          response: { eject: jest.fn(), use: jest.fn() }
        },
        defaults: {
          withCredentials: true
        }
      };
    }
  };
});

describe(' get members lists ', () => {
  it('should return default members list with mainnet as default', async () => {
    // membersservice should be called with mainnet url
    const { data } = mainnetMembersMockValue;
    axios.get = jest.fn().mockResolvedValue({ data });
    await expect(getMembersService()).resolves.toEqual(mainnetMembersMockValue.data);
  });

  it('should return members list with testnet', async () => {
    const { data } = testnetMembersMockValue;
    axios.get = jest.fn().mockResolvedValue({ data });
    await expect(getMembersService('testnet')).resolves.toEqual(testnetMembersMockValue.data);
  });
});
