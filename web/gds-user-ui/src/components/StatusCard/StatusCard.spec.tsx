import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render, screen } from 'utils/test-utils';
import StatusCard from '.';

describe('<StatusCard />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should display a teal color', () => {
    const isOnline = 'HEALTHY';
    render(<StatusCard isOnline={isOnline} />);

    expect(screen.getByTestId('status__color')).toHaveAttribute('fill', '#60C4CA');
  });

  it('should display a gray color', () => {
    const isOnline = 'UNKNOWN';
    render(<StatusCard isOnline={isOnline} />);

    expect(screen.getByTestId('status__color')).toHaveAttribute('fill', '#C4C4C4');
  });
});
