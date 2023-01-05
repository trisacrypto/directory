import { render } from '@/utils/test-utils';
import { Tasks } from '../dashboard';

describe('<Tasks />', () => {
  it('Should keep the same UI', () => {
    const { container } = render(<Tasks />);

    //
    expect(container).toMatchSnapshot('pending-registrations');
  });
});
