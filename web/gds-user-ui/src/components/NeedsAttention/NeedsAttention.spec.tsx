import userEvent from '@testing-library/user-event';
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render, screen } from 'utils/test-utils';
import NeedsAttention from '.';

describe('<NeedsAttention />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should display the right button text', () => {
    render(<NeedsAttention buttonText="This is a button text" />);

    expect(screen.getByRole('button').textContent).toBe('This is a button text');
  });

  it('should call button event listener', () => {
    const handleClickMock = jest.fn();

    render(<NeedsAttention onClick={handleClickMock} buttonText="This is a button text" />);

    userEvent.click(screen.getByRole('button'));

    expect(handleClickMock).toHaveBeenCalledTimes(1);
  });
});
