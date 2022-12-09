import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render, waitFor } from 'utils/test-utils';
import userEvent from '@testing-library/user-event';
import ChangeNameForm from '../ChangeNameForm';

function renderComponent() {
  const Props = {
    onCloseModal: jest.fn()
  };
  return render(<ChangeNameForm {...Props} />);
}

// mock useSelector and useDispatch

describe('<ChangeNameForm />  ', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
    const userStateMock = {
      name: 'test name',
      email: 'test@local.test'
    };
    jest.mock('react-redux', () => ({
      ...jest.requireActual('react-redux'),
      useSelector: jest.fn().mockReturnValueOnce({
        user: userStateMock
      }),
      useDispatch: jest.fn()
    }));
  });

  it('should fill the form', () => {
    const { getByTestId } = renderComponent();
    const firstNameInput = getByTestId('first_name');
    const lastNameInput = getByTestId('last_name');

    userEvent.type(firstNameInput, 'fistName');
    userEvent.type(lastNameInput, 'lastName');
    expect(firstNameInput).toHaveValue('fistName');
    expect(lastNameInput).toHaveValue('lastName');
  });

  it('should submit the form', async () => {
    const { getByTestId } = renderComponent();
    const firstNameInput = getByTestId('first_name');
    const lastNameInput = getByTestId('last_name');
    const submitButton = getByTestId('save_button');
    userEvent.type(firstNameInput, 'fistName');
    userEvent.type(lastNameInput, 'lastName');
    userEvent.click(submitButton);
    await waitFor(() => {
      expect(submitButton).toBeDisabled();
    });
  });
  afterEach(() => {
    jest.clearAllMocks();
  });
});
