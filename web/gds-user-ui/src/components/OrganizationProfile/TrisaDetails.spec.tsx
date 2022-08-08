import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render, screen } from 'utils/test-utils';
import TrisaDetail from './TrisaDetail';

describe('<TrisaDetails />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });
  it('should display data correctly', () => {
    const data = {
      organization: {
        vasp_id: 'VASP-ID',
        first_listed: '22/22/22',
        verified_on: '04/11/24',
        last_updated: '04/12/25'
      }
    };

    const { debug } = render(<TrisaDetail data={data} />);

    expect(screen.getByTestId('vasp_id').textContent).toBe('VASP-ID');
    expect(screen.getByTestId('first_listed').textContent).toBe('22/22/22');
    expect(screen.getByTestId('verified_on').textContent).toBe('04/11/24');
    expect(screen.getByTestId('last_updated').textContent).toBe('04/12/25');
  });

  it('should display N/A when there are no data', () => {
    const data = {};

    render(<TrisaDetail data={data} />);

    expect(screen.getByTestId('vasp_id').textContent).toBe('N/A');
    expect(screen.getByTestId('first_listed').textContent).toBe('N/A');
    expect(screen.getByTestId('verified_on').textContent).toBe('N/A');
    expect(screen.getByTestId('last_updated').textContent).toBe('N/A');
  });
});
