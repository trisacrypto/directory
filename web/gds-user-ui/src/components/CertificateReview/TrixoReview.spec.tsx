import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render } from 'utils/test-utils';
import LegalPersonReview from './LegalPersonReview';
import TrixoReview from './TrixoReview';

const defaultValues = {
  trixo: {
    primary_national_jurisdiction: 'AF',
    primary_regulator: '',
    other_jurisdictions: [],
    financial_transfers_permitted: 'no',
    has_required_regulatory_program: 'no',
    conducts_customer_kyc: false,
    kyc_threshold: 0,
    kyc_threshold_currency: 'USD',
    must_comply_travel_rule: false,
    applicable_regulations: [{ name: 'FATF Recommendation 16' }],
    compliance_threshold: 3000,
    compliance_threshold_currency: 'USD',
    must_safeguard_pii: false,
    safeguards_pii: false
  }
};

describe('<TrixoReview />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  beforeEach(() => {
    localStorage.clear();
  });

  it('should match snapshot', () => {
    const { container } = render(<TrixoReview />);

    // expect(localStorage.getItem).toHaveBeenCalledWith('certificateForm');
    // expect(Object.keys(localStorage.__STORE__).length).toBe(1);

    expect(container).toMatchSnapshot();
  });

  afterAll(() => {
    jest.clearAllMocks();
  });
});
