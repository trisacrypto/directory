import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render, screen } from 'utils/test-utils';
import StatCard from '.';

describe('<StatCard />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should render props correctly', () => {
    render(<StatCard title="Network Status">10</StatCard>);

    expect(screen.getByTestId('start-card__title').textContent).toBe('Network Status');
    expect(screen.getByTestId('start-card__body').textContent).toBe('10');
  });
});
