import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { act, render, screen } from 'utils/test-utils';
import VaspVerification from '../VaspVerification';

describe('<Line />', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should match snapshot', () => {
    const { container } = render(<VaspVerification />);

    expect(container).toMatchSnapshot();
  });
});
