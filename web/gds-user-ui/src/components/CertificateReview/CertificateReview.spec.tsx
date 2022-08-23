import userEvent from '@testing-library/user-event';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render, screen } from 'utils/test-utils';
import CertificateReview from '.';

// mock useformcontext of react-hook-form
jest.mock('react-hook-form', () => ({
  ...jest.requireActual('react-hook-form'),
  useFormContext: () => ({
    handleSubmit: () => jest.fn(),
    getValues: () => ({})
  })
}));

describe('<CertificateReview />', () => {
  beforeEach(() => {
    dynamicActivate('en');
  });

  it('should render ReviewSubmit Component when user reached the last step', () => {
    const initialState = {
      stepper: {
        currentStep: 2,
        steps: [
          { key: 1, status: 'complete', data: {} },
          { key: 2, status: 'progress' }
        ],
        lastStep: null,
        hasReachSubmitStep: true
      }
    };
    render(<CertificateReview />, { preloadedState: initialState });

    const mainnetSubmitButton = screen.getByRole('button', {
      name: /submit mainnet registration/i
    });

    expect(mainnetSubmitButton).toBeInTheDocument();
  });

  it("should render ReviewSubmit Component when user didn't reach the last step", () => {
    const initialState = {
      stepper: {
        currentStep: 2,
        steps: [
          { key: 1, status: 'complete', data: {} },
          { key: 2, status: 'progress' }
        ],
        lastStep: null,
        hasReachSubmitStep: false
      }
    };
    render(<CertificateReview />, { preloadedState: initialState });

    const reviewEl = screen.getByTestId('review');

    expect(reviewEl).toBeInTheDocument();
  });
});
