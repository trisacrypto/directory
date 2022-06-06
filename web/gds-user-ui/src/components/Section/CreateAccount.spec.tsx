import userEvent from '@testing-library/user-event';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render, screen, waitFor } from 'utils/test-utils';
import CreateAccount from './CreateAccount';
import { BrowserRouter as Router } from 'react-router-dom';
const mockSignWithEmail = jest.fn((values) => {
  return Promise.resolve(values);
});

const mockSignWithSocial = jest.fn((values) => {
  return Promise.resolve(values);
});

describe('<CreateAccount />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should submit form when fields are all filled', async () => {
    render(
      <Router>
        <CreateAccount
          handleSocialAuth={mockSignWithEmail}
          handleSignUpWithEmail={mockSignWithSocial}
        />
      </Router>
    );

    const username = screen.getByTestId('username-field');
    userEvent.type(username, 'test@email.com');

    const password = screen.getByTestId('password-field');
    userEvent.type(password, 'AA!45aaa');

    const submitButton = screen.getByRole('button', { name: /create an account/i });

    userEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.queryAllByRole('alert')).toHaveLength(0);
    });
  });

  it('should throw validation error when all field are empty', async () => {
    const mockHandleAuthFn = jest.fn();
    const mockHandleSignUpWithEmail = jest.fn();

    render(
      <CreateAccount
        handleSocialAuth={mockSignWithEmail}
        handleSignUpWithEmail={mockSignWithSocial}
      />
    );

    const submitButton = screen.getByRole('button', { name: /create an account/i });

    userEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getAllByRole('alert')).toHaveLength(2);
    });
  });
});
