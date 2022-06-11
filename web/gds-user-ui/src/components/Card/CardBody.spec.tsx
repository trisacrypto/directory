import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render, screen } from 'utils/test-utils';
import CardBody from './CardBody';

describe('<CardBody />', () => {
  beforeEach(() => {
    dynamicActivate('en');
  });

  it('should render', () => {
    render(<CardBody>Test</CardBody>);

    expect(screen.getByText(/Test/i)).toBeInTheDocument();
  });
});
