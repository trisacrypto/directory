import userEvent from '@testing-library/user-event';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render, screen, waitFor } from 'utils/test-utils';
import CreateAccount from './CreateAccount';

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
      <CreateAccount
        handleSocialAuth={mockSignWithEmail}
        handleSignUpWithEmail={mockSignWithSocial}
      />
    );

    const username = screen.getByTestId('username-field');
    userEvent.type(username, 'test@email.com');

    const password = screen.getByTestId('password-field');
    userEvent.type(password, 'AA!45aaa');

    const submitButton = screen.getByRole('button', { name: /create your account/i });

    userEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.queryAllByRole('alert')).toHaveLength(0);
    });
  });

  it('should throw validation error when all field are empty', async () => {
    render(
      <CreateAccount
        handleSocialAuth={mockSignWithEmail}
        handleSignUpWithEmail={mockSignWithSocial}
      />,
      { route: '/auth/login' }
    );

    const submitButton = screen.getByRole('button', { name: /create your account/i });

    userEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getAllByRole('alert')).toHaveLength(2);
    });
  });

  it('should call google login function', async () => {
    render(
      <CreateAccount
        handleSocialAuth={mockSignWithEmail}
        handleSignUpWithEmail={mockSignWithSocial}
      />,
      { route: '/auth/login' }
    );

    const submitButton = screen.getByRole('button', { name: /continue with google/i });

    userEvent.click(submitButton);

    expect(mockSignWithSocial).toHaveBeenCalled();
    expect(mockSignWithSocial).toHaveBeenCalledTimes(1);
  });

  describe('Show button', () => {
    it('should show password when we click on show button', () => {
      render(
        <CreateAccount
          handleSocialAuth={mockSignWithEmail}
          handleSignUpWithEmail={mockSignWithSocial}
        />,
        { route: '/auth/login' }
      );

      const passwordInputEl = screen.getByPlaceholderText(/password/i) as HTMLInputElement;
      userEvent.type(passwordInputEl, 'test password');

      const showBtn = screen.getByRole('button', { name: /show/i });

      userEvent.click(showBtn);
      expect(passwordInputEl.type).toBe('text');

      userEvent.click(showBtn);
      expect(passwordInputEl.type).toBe('password');
    });
  });
});
