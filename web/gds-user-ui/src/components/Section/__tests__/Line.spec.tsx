import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render, screen } from 'utils/test-utils';
import { Line } from '../Line';

describe('<Line />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should render title', () => {
    render(<Line title="Option 1. Set Up Your Own TRISA Node">test</Line>);

    expect(screen.getByTestId('title').textContent).toBe('Option 1. Set Up Your Own TRISA Node');
  });
});
