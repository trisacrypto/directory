import { faker } from '@faker-js/faker';
import { act, render, screen } from 'utils/test-utils';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import ContactsForm from 'components/Contacts/ContactsForm';
import ContactForm from '..';
import { fireEvent } from '@testing-library/react';

function renderComponent() {
  return render(<ContactsForm />);
}

jest.mock('hooks/useFetchCertificateStep', () => ({
  useFetchCertificateStep: () => ({
    certificateStep: {
      form: jest.fn(),
      errors: jest.fn()
    },
    isFetchingCertificateStep: false
  })
}));

describe('ContactsForm', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should render correctly', () => {
    const mockData = {
      contacts: {
        technical: {
          name: 'Technical',
          email: 'abc@123.com',
          phone_number: '555-555-5555'
        },
        legal: {
          name: 'Legal',
          email: 'def@456.com',
          phone_number: '555-555-5555'
        }
      }
    };
    render(<ContactsForm data={mockData} />);
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
