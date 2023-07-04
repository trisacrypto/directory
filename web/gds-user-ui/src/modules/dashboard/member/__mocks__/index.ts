export const mainnetMembersMockValue: any = {
  vasps: [
    {
      id: 'fe29b582-20b3-4b75-acb2-4d014af25f28',
      registered_directory: 'vaspdirectory.net',
      common_name: 'mainnet.firecoinex.co',
      endpoint: 'mainnet.firecoinex.co:443',
      name: 'FireCoin Exchange',
      website: 'https://www.firecoinex.co/',
      country: 'SG',
      business_category: 'BUSINESS_ENTITY',
      vasp_categories: ['OTC'],
      verified_on: '2022-04-20T10:45:33Z',
      status: 'VERIFIED',
      first_listed: '2022-04-20T11:10:06Z',
      last_updated: '2022-11-19T09:32:45Z'
    },
    {
      id: '52842ac3-7d9d-4520-8bfb-e6083dbdc8aa',
      registered_directory: 'vaspdirectory.net',
      common_name: 'mainnet.testmachine.com',
      endpoint: 'mainnet.testmachine.com:443',
      name: 'Test Machine',
      website: 'https://testmachine.com',
      country: 'US',
      business_category: 'BUSINESS_ENTITY',
      vasp_categories: ['Exchange', 'Other'],
      verified_on: '2023-03-22T03:55:00Z',
      status: 'VERIFIED',
      first_listed: '2023-03-22T04:22:18Z',
      last_updated: '2023-05-10T04:44:00Z'
    },
    {
      id: 'ec6f2056-726b-4e8c-a916-f7359f6f5581',
      registered_directory: 'vaspdirectory.net',
      common_name: 'mainnet.newcoinex.ai',
      endpoint: 'mainnet.newcoinex.ai:8221',
      name: 'New Coin Exchange',
      website: 'https://newcoinex.ai/',
      country: 'US',
      business_category: 'BUSINESS_ENTITY',
      vasp_categories: ['Other'],
      verified_on: '2023-01-20T12:55:03Z',
      status: 'VERIFIED',
      first_listed: '2023-01-20T14:10:26Z',
      last_updated: '2023-06-18T10:30:01Z'
    }
  ],
  next_page_token: ''
};

export const testnetMembersMockValue: any = {
  vasps: [
    {
      id: '1b99f17d-3441-4885-b9df-8f3475f7e1b4',
      registered_directory: 'trisatest.net',
      common_name: 'trisa-travelrule.sendcoin.io',
      endpoint: 'trisa-travelrule.sendcoin.io:443',
      name: 'SendCoin VASP',
      website: 'https://www.sendcoin.io/',
      country: 'SG',
      business_category: 'BUSINESS_ENTITY',
      vasp_categories: ['OTC'],
      verified_on: '2022-02-10T17:12:23Z',
      status: 'VERIFIED',
      first_listed: '2022-02-10T13:43:55Z',
      last_updated: '2023-04-09T04:02:50Z'
    },
    {
      id: 'aa4c6714-49dc-411c-adfa-b8edc4d58cd7',
      registered_directory: 'trisatest.net',
      common_name: 'testing.example.com',
      endpoint: 'testing.example.com:443',
      name: 'Example Crypto',
      website: 'https://example.com',
      country: 'US',
      business_category: 'BUSINESS_ENTITY',
      vasp_categories: ['Exchange', 'Other'],
      verified_on: '2021-12-07T20:22:00Z',
      status: 'VERIFIED',
      first_listed: '2021-12-01T23:22:18Z',
      last_updated: '2023-01-23T19:19:43Z'
    },
    {
      id: '688059d6-9d14-4b49-8435-f641ba1dec3a',
      registered_directory: 'trisatest.net',
      common_name: 'testnet.spudcoin.ai',
      endpoint: 'testnet.spudcoin.ai:8221',
      name: 'SpudCoin',
      website: 'https://spudcoin.ai/',
      country: 'US',
      business_category: 'BUSINESS_ENTITY',
      vasp_categories: ['Other'],
      verified_on: '2021-07-29T19:11:03Z',
      status: 'VERIFIED',
      first_listed: '2021-07-23T17:10:26Z',
      last_updated: '2022-12-27T18:59:01Z'
    },
    {
      id: '62291255-1ea2-4932-8248-22af4abde298',
      registered_directory: 'trisatest.net',
      common_name: 'test.signal.co.fr',
      endpoint: 'test.signal.co.fr:9212',
      name: 'Signal Coin France',
      website: 'https://ciphertrace.com/',
      country: 'US',
      business_category: 'BUSINESS_ENTITY',
      vasp_categories: ['Other'],
      verified_on: '2021-06-23T17:46:10Z',
      status: 'VERIFIED',
      first_listed: '2021-06-23T17:12:00Z',
      last_updated: '2022-12-27T18:59:22Z'
    },
    {
      id: 'dae1555d-e4cf-4bfb-9858-9b86db71ccb6',
      registered_directory: 'trisatest.net',
      common_name: 'testnet.bitfriend.link',
      endpoint: 'testnet.bitfriend.link:443',
      name: 'BitFriendly',
      website: 'https://bitfriend.link/',
      country: 'US',
      business_category: 'BUSINESS_ENTITY',
      vasp_categories: ['Kiosk'],
      verified_on: '2021-09-14T10:40:30Z',
      status: 'VERIFIED',
      first_listed: '2021-09-01T19:46:04Z',
      last_updated: '2023-01-23T19:29:43Z'
    }
  ],
  next_page_token: 'mLB9CU8O8xQj2XEyjAtlfvTj9imoXnLv/1p8fTLchTg='
};

export const getMockValue = (network: string): any => {
  switch (network) {
    case 'mainnet':
      return mainnetMembersMockValue;
    case 'testnet':
      return testnetMembersMockValue;
    default:
      return mainnetMembersMockValue;
  }
};
