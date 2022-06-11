import { render } from 'utils/test-utils';
import ReviewsSummary from './ReviewsSummary';

describe('<ReviewSummary />', () => {
  it('should match snapshot', () => {
    const { container } = render(<ReviewsSummary />);

    expect(container).toMatchSnapshot();
  });
});
