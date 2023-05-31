import { faker } from '@faker-js/faker';
import { act, render, screen } from 'utils/test-utils';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import ContactsForm from 'components/Contacts/ContactsForm';
import ContactForm from '..';
import { fireEvent } from '@testing-library/react';

function renderComponent() {
  return render(<ContactsForm />);
}
// mock chakra useDisclosure hook
jest.mock('@chakra-ui/react', () => ({
  ...jest.requireActual('@chakra-ui/react'),
  useDisclosure: jest.fn(() => ({
    isOpen: false,
    onClose: jest.fn(),
    onOpen: jest.fn()
  }))
}));

jest.mock('hooks/useFetchCertificateStep', () => ({
  useFetchCertificateStep: () => ({
    certificateStep: {
      form: jest.fn(),
      errors: jest.fn()
    },

    isFetchingCertificateStep: false
  }),
  __esModule: true
}));

describe('ContactsForm', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should render', () => {
    const { container } = renderComponent();
    expect(container).toMatchSnapshot();
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
