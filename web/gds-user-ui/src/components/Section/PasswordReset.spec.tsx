import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render, screen, fireEvent } from 'utils/test-utils';
import PasswordReset from './PasswordReset';

jest.mock('../../hooks/useCustomAuth0.ts', () => {
  return jest.fn(() => ({
    auth0ResetPassword: (options: any) => {
      return Promise.resolve(options);
    }
  }));
});

describe('<PasswordReset />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should display required email error message', async () => {
    await act(async () => {
      render(<PasswordReset />);
    });

    fireEvent.input(screen.getByRole('textbox'), {
      target: {
        value: ''
      }
    });

    await act(async () => {
      fireEvent.submit(screen.getByRole('button'));
    });

    expect(screen.getByTestId('error-message').textContent).toBe('Email is required');
  });

  it('should display invalid email error message', async () => {
    await act(async () => {
      render(<PasswordReset />);
    });

    fireEvent.input(screen.getByRole('textbox'), {
      target: {
        value: 'email.com'
      }
    });

    await act(async () => {
      fireEvent.submit(screen.getByRole('button'));
    });

    expect(screen.getByRole('alert').textContent).toBe('Email is invalid');
  });

  it('should display success message when email submitted sucessfully', async () => {
    await act(async () => {
      render(<PasswordReset />);
    });

    fireEvent.input(screen.getByRole('textbox'), {
      target: {
        value: 'email@email.com'
      }
    });

    await act(async () => {
      fireEvent.submit(screen.getByRole('button'));
    });

    expect(screen.getByTestId('success__alert').textContent).toBe(
      'Thank you. We have sent instructions to reset your password to email@email.com. The link to reset your password expires in 24 hours.'
    );
  });

  afterAll(() => {
    jest.clearAllMocks();
  });
});
