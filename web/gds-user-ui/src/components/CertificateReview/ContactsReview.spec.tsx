import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render } from 'utils/test-utils';
import ContactsReview from './ContactsReview';

const defaultValues = {
  contacts: {
    administrative: { name: '', email: '', phone: '' },
    technical: { name: 'Test', email: 'test@test.com', phone: '+4535557757' },
    billing: { name: '', email: '' },
    legal: { name: 'Info ', email: 'info@gmail.com', phone: '+93400583484' }
  }
};

describe('<ContactsReview />', () => {
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

    const { container } = render(<ContactsReview data={defaultValues.contacts} />);

    // expect(localStorage.getItem).toHaveBeenCalledWith('certificateForm');
    // expect(Object.keys(localStorage.__STORE__).length).toBe(1);

    expect(container).toMatchSnapshot();
  });

  afterAll(() => {
    jest.clearAllMocks();
  });
});
