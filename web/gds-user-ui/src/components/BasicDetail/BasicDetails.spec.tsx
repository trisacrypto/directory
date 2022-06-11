import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render, screen, waitFor } from 'utils/test-utils';
import BasicDetails from '.';

describe('<BasicDetails />', () => {
  beforeEach(() => {
    dynamicActivate('en');
  });

  describe('<SectionStatus />', () => {
    it('should show not saved when step status is in progress', async () => {
      const initialState = {
        stepper: {
          currentStep: 1,
          steps: [{ key: 1, status: 'progress', data: {} }],
          lastStep: null,
          hasReachSubmitStep: false
        }
      };

      await waitFor(async () => {
        render(<BasicDetails />, { preloadedState: initialState });
      });

      expect(screen.getByText(/(not saved)/i)).toBeInTheDocument();
    });

    it('should show saved when step status is completed', async () => {
      const initialState = {
        stepper: {
          currentStep: 1,
          steps: [{ key: 1, status: 'complete', data: {} }],
          lastStep: null,
          hasReachSubmitStep: false
        }
      };

      await waitFor(async () => {
        render(<BasicDetails />, { preloadedState: initialState });
      });

      expect(screen.getByText(/(saved)/i)).toBeInTheDocument();
    });
  });
});
