import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render } from 'utils/test-utils';
import TrisaImplementationReview from './TrisaImplementationReview';

const defaultValues = {
  mainnet: {
    common_name: 'testnet.kyc.com',
    endpoint: 'testnet.kyc.com:443'
  },
  testnet: {
    common_name: 'trisa.kyc.com',
    endpoint: 'trisa.kyc.com:443'
  }
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
    // localStorage.setItem('certificateForm', JSON.stringify(defaultValues));

    const { container } = render(<TrisaImplementationReview />);

    // expect(localStorage.getItem).toHaveBeenCalledWith('certificateForm');
    // expect(Object.keys(localStorage.__STORE__).length).toBe(1);

    expect(container).toMatchSnapshot();
  });

  afterAll(() => {
    jest.clearAllMocks();
  });
});
