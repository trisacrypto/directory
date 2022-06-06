import userEvent from '@testing-library/user-event';
import Login from 'components/Section/Login';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render, screen, waitFor } from 'utils/test-utils';
import { BrowserRouter as Router } from 'react-router-dom';
const mockSignWithEmail = jest.fn((values) => {
  return Promise.resolve(values);
});

const mockSignWithSocial = jest.fn((values) => {
  return Promise.resolve(values);
});

describe('<Login />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  beforeEach(() => {
    render(
      <Router>
        <Login handleSignWithEmail={mockSignWithEmail} handleSignWithSocial={mockSignWithSocial} />
      </Router>
    );
  });

  describe('Email', () => {
    it('should throw an invalid email error', async () => {
      const email = screen.getByTestId(/email/i);
      userEvent.type(email, 'test.mail.com');

      const loginButton = screen.getByTestId(/login-btn/i);
      userEvent.click(loginButton);

      await waitFor(() => {
        expect(screen.getByText(/Email Address is not valid/i)).toBeInTheDocument();
      });
    });
  });

  it('should throw validation errors when all fields are empty', async () => {
    const loginButton = screen.getByTestId(/login-btn/i);

    userEvent.click(loginButton);

    await waitFor(() => {
      expect(screen.getAllByRole('alert')).toHaveLength(2);
    });
    expect(mockSignWithEmail).not.toHaveBeenCalled();
  });

  it('should login when all fields are filled', async () => {
    const username = 'test@email.com';
    const password = '@##@test';
    const usernameEl = screen.getByTestId(/email/i);
    const passwordEl = screen.getByTestId(/password/i);

    userEvent.type(usernameEl, username);
    userEvent.type(passwordEl, password);

    const loginButton = screen.getByTestId(/login-btn/i);

    userEvent.click(loginButton);

    await waitFor(() => expect(screen.queryAllByRole('alert')).toHaveLength(0));
    expect(mockSignWithEmail).toHaveBeenCalledTimes(1);
  });

  it('should call mockSignWithGoogle when continue with google button is clicked ', () => {
    const signWithGoogleButton = screen.getByTestId(/signin-with-google/i);

    userEvent.click(signWithGoogleButton);

    expect(mockSignWithSocial).toHaveBeenCalledTimes(1);
  });
});
