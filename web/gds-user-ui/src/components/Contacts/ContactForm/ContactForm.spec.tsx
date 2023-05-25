import { faker } from '@faker-js/faker';
import userEvent from '@testing-library/user-event';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { fireEvent, render, screen, waitFor, within } from 'utils/test-utils';
import ContactForm from '.';

describe('<ContactForm />', () => {
  beforeEach(() => {
    dynamicActivate('en');
  });

  it('should render correctly', () => {
    const title = faker.random.words();
    const desc = faker.lorem.paragraph();
    const name = 'test';
    render(<ContactForm name={name} title={title} description={desc} />);

    expect(screen.getByTestId('title').textContent).toBe(title);
    expect(screen.getByTestId('description').textContent).toBe(desc);

    // fullname
    const fullNameEl = screen.getByTestId('fullName') as HTMLInputElement;
    expect(fullNameEl).toBeVisible();
    expect(fullNameEl.name).toBe(`${name}.name`);

    // fullname
    const emailAddressEl = screen.getByTestId('email') as HTMLInputElement;
    expect(emailAddressEl).toBeVisible();
    expect(emailAddressEl.name).toBe(`${name}.email`);
    expect(emailAddressEl.type).toBe(`email`);

    // phoneNumber
    const phoneNumberEl = screen.getByTestId('phoneNumber') as HTMLInputElement;
    userEvent.type(phoneNumberEl, `123445566`);
    expect(phoneNumberEl).toHaveValue(`+1 234 455 66`);
  });

  it('should render the correct phone message hint when name is contacts.legal', () => {
    const title = faker.random.words();
    const desc = faker.lorem.paragraph();
    const name = 'contacts.legal';
    render(<ContactForm name={name} title={title} description={desc} />);

    const phoneNumberEl = screen.getByTestId(/legal-contact-phone-number-hint/i);
    expect(phoneNumberEl.textContent).toBe(
      'A business phone number is required to complete physical verification for MainNet registration. Please provide a phone number where the Legal/ Compliance contact can be contacted.'
    );
  });

  it('should render the correct phone message hint when name is different to contacts.legal', () => {
    const title = faker.random.words();
    const desc = faker.lorem.paragraph();
    const name = 'contacts.technical';
    render(<ContactForm name={name} title={title} description={desc} />);

    const phoneNumberEl = screen.getByTestId(/legal-contact-phone-number-hint/i);
    expect(phoneNumberEl.textContent).toBe('If supplied, use full phone number with country code.');
  });

  describe('FullName', () => {
    it('should render an error message', async () => {
      const title = faker.random.words();
      const desc = faker.lorem.paragraph();
      const name = 'contacts.technical';
      render(<ContactForm name={name} title={title} description={desc} />);
      const fullNameEl = screen.getByTestId('fullName') as HTMLInputElement;

      fireEvent.change(fullNameEl, { target: { value: 'test' } });
      fireEvent.change(fullNameEl, { target: { value: '' } });

      expect(screen.getByText(/Preferred name for email communication./i)).toBeVisible();
    });
  });
});
