import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render } from 'utils/test-utils';
import mockAxios from '__mocks__/axios';
import CertificateDetails from './CertificateDetails';

const mock = {
  network_error: {},
  testnet: [
    {
      serial_number: 'DFD0460B63FA147F6E671CF45AF6369E',
      issued_at: '2021-08-03T21:59:31Z',
      expires_at: '2021-08-16T13:12:41Z',
      revoked: false,
      details: {
        chain: '',
        data: '',
        issuer: {
          common_name: 'CipherTrace Issuing CA',
          country: ['US'],
          locality: ['Menlo Park'],
          organization: ['CipherTrace Inc'],
          organizational_unit: [],
          postal_code: [],
          province: ['California'],
          serial_number: '',
          street_address: []
        },
        not_after: '2021-08-16T13:12:41Z',
        not_before: '2021-08-03T21:59:31Z',
        public_key_algorithm: 'RSA',
        revoked: false,
        serial_number: 'DFD0460B63FA147F6E671CF45AF6369E',
        signature: '',
        signature_algorithm: 'SHA384-RSA',
        subject: {
          common_name: 'uniform',
          country: ['US'],
          locality: ['Menlo Park'],
          organization: ['CipherTrace Inc'],
          organizational_unit: [],
          postal_code: [],
          province: ['California'],
          serial_number: '',
          street_address: []
        },
        version: '3'
      }
    }
  ],
  mainnet: [
    {
      serial_number: 'FE4348E441E1D40FA41039E5780C75AD',
      issued_at: '2021-10-08T02:54:28Z',
      expires_at: '2021-10-09T17:31:45Z',
      revoked: false,
      details: {
        chain: '',
        data: '',
        issuer: {
          common_name: 'CipherTrace Issuing CA',
          country: ['US'],
          locality: ['Menlo Park'],
          organization: ['CipherTrace Inc'],
          organizational_unit: [],
          postal_code: [],
          province: ['California'],
          serial_number: '',
          street_address: []
        },
        not_after: '2021-10-09T17:31:45Z',
        not_before: '2021-10-08T02:54:28Z',
        public_key_algorithm: 'RSA',
        revoked: false,
        serial_number: 'FE4348E441E1D40FA41039E5780C75AD',
        signature: '',
        signature_algorithm: 'SHA384-RSA',
        subject: {
          common_name: 'victor',
          country: ['US'],
          locality: ['Menlo Park'],
          organization: ['CipherTrace Inc'],
          organizational_unit: [],
          postal_code: [],
          province: ['California'],
          serial_number: '',
          street_address: []
        },
        version: '3'
      }
    },
    {
      serial_number: '2C3B5861D89F2D32829B325AFB39293F',
      issued_at: '2021-08-04T00:28:34Z',
      expires_at: '2021-08-11T13:53:13Z',
      revoked: true,
      details: {
        chain: '',
        data: '',
        issuer: {
          common_name: 'CipherTrace Issuing CA',
          country: ['US'],
          locality: ['Menlo Park'],
          organization: ['CipherTrace Inc'],
          organizational_unit: [],
          postal_code: [],
          province: ['California'],
          serial_number: '',
          street_address: []
        },
        not_after: '2021-08-11T13:53:13Z',
        not_before: '2021-08-04T00:28:34Z',
        public_key_algorithm: 'RSA',
        revoked: true,
        serial_number: '2C3B5861D89F2D32829B325AFB39293F',
        signature: '',
        signature_algorithm: 'SHA384-RSA',
        subject: {
          common_name: 'zulu',
          country: ['US'],
          locality: ['Menlo Park'],
          organization: ['CipherTrace Inc'],
          organizational_unit: [],
          postal_code: [],
          province: ['California'],
          serial_number: '',
          street_address: []
        },
        version: '3'
      }
    }
  ]
};

describe('<CertificateDetails />', () => {
  beforeEach(() => {
    dynamicActivate('en');
  });

  afterEach(() => {
    mockAxios.reset();
  });

  it('should render that', () => {
    render(<CertificateDetails />);

    mockAxios.get.mockResolvedValueOnce({ data: mock });
    expect(mockAxios.get).toHaveBeenCalledWith('/certificates');
  });
});
