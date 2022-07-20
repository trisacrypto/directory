import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render, screen } from 'utils/test-utils';
import BasicDetailsForm from '.';

describe('<BasicDetailsForm />', () => {
  beforeEach(() => {
    dynamicActivate('en');
  });
  it('should render correctly', () => {
    render(<BasicDetailsForm />);

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
