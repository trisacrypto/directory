import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render, screen } from 'utils/test-utils';
import BasicDetailsForm from '.';

jest.mock('@chakra-ui/react', () => ({
  ...jest.requireActual('@chakra-ui/react'),
  useDisclosure: jest.fn(() => ({
    isOpen: false,
    onClose: jest.fn(),
    onOpen: jest.fn()
  }))
}));

describe('<BasicDetailsForm />', () => {
  beforeEach(() => {
    dynamicActivate('en');
  });
  it('should render correctly', () => {
    const mockData = {
      website: 'https://www.google.com',
      vasp_categories: ['VASP'],
      business_category: 'Crypto Exchange',
      organization_name: 'Google',
      established_on: '2021-01-01'
    };
    const mockHandleSubmit = jest.fn();
    render(<BasicDetailsForm data={mockData} onNextStepClick={mockHandleSubmit} />);

    // organization_name
    const organizationName = screen.getByRole('textbox', {
      name: /organization name/i
    }) as HTMLInputElement;
    expect(organizationName).toBeVisible();
    expect(organizationName.name).toBe('organization_name');

    // website
    const website = screen.getByRole('textbox', { name: /website/i }) as HTMLInputElement;
    expect(website).toBeVisible();
    expect(website.name).toBe('website');

    // established_on
    const dateOfIncorporation = screen.getByLabelText(
      /date of incorporation \/ establishment/i
    ) as HTMLInputElement;
    expect(dateOfIncorporation).toBeVisible();
    expect(dateOfIncorporation.name).toBe('established_on');
  });
});
