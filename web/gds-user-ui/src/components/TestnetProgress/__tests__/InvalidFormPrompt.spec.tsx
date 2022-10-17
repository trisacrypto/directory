import userEvent from '@testing-library/user-event';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render, act, screen } from 'utils/test-utils';
import InvalidFormPrompt from '../InvalidFormPrompt';

describe('<InvalidFormPrompt />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should', () => {
    const handleContinueClick = jest.fn();
    const handleClose = jest.fn();
    const isOpen = true;

    render(
      <InvalidFormPrompt
        isOpen={isOpen}
        onClose={handleClose}
        handleContinueClick={handleContinueClick}
      />
    );

    const continueButton = screen.getByRole('button', { name: /continue/i });

    userEvent.click(continueButton);

    expect(handleContinueClick).toHaveBeenCalled();
  });
});
