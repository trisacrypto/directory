import { render } from 'utils/test-utils';
import VaspVerification from '../VaspVerification';

describe('<VaspVerification />', () => {
  it('should', () => {
    const { container } = render(<VaspVerification />);

    expect(container).toMatchSnapshot();
  });
});
