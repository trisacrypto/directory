import { render, act } from 'utils/test-utils';
import BasicDetailsReview from './BasicDetailsReview';
import { i18n } from '@lingui/core';
import { dynamicActivate } from 'utils/i18nLoaderHelper';

const defaultValues = {
  entity: {
    country_of_registration: 'AF'
  },
  website: 'http://kyc.com',
  business_category: 'GOVERNMENT_ENTITY',
  vasp_categories: ['P2P', 'Kiosk'],
  established_on: '2022-04-22',
  organization_name: 'KYC'
};

describe('<BasicDetailsReview />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should match snapshot', () => {
    // localStorage.setItem('certificateForm', JSON.stringify(defaultValues));

    const { container } = render(<BasicDetailsReview data={defaultValues} />);

    // expect(localStorage.getItem).toHaveBeenCalledWith('certificateForm');
    // expect(Object.keys(localStorage.__STORE__).length).toBe(1);

    expect(container).toMatchSnapshot();
  });

  afterAll(() => {
    jest.clearAllMocks();
  });
});
