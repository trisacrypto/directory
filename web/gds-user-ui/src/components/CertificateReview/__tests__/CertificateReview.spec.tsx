import { render } from 'utils/test-utils';
import CertificateReview from '../index';
describe('<CertificateReview />', () => {
  it('should match snapshot', () => {
    const { container } = render(<CertificateReview />);

    expect(container).toMatchSnapshot();
  });
});
