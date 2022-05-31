import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render } from 'utils/test-utils';
import TrisaImplementationReview from './TrisaImplementationReview';

const defaultValues = {
  trisa_endpoint_testnet: {
    trisa_endpoint: '',
    common_name: 'testnet.kyc.com',
    endpoint: 'testnet.kyc.com:443'
  },
  trisa_endpoint_mainnet: {
    trisa_endpoint: '',
    common_name: 'trisa.kyc.com',
    endpoint: 'trisa.kyc.com:443'
  },
  website: 'http://kyc.com',
  business_category: 'GOVERNMENT_ENTITY',
  vasp_categories: ['P2P', 'Kiosk'],
  established_on: '2022-04-22',
  organization_name: 'KYC'
};

describe('<TrisaImplementationReview />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  beforeEach(() => {
    localStorage.clear();
  });

  it('should match snapshot', () => {
    localStorage.setItem('certificateForm', JSON.stringify(defaultValues));

    const { container } = render(<TrisaImplementationReview />);

    console.log('[after all] localStorage', localStorage.__STORE__);
    expect(localStorage.getItem).toHaveBeenCalledWith('certificateForm');
    expect(Object.keys(localStorage.__STORE__).length).toBe(1);

    expect(container).toMatchSnapshot();
  });

  afterAll(() => {
    jest.clearAllMocks();
  });
});
