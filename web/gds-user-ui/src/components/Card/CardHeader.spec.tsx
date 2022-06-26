import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { render, screen } from 'utils/test-utils';
import CardHeader from './CardHeader';

describe('<CardHeader />', () => {
  beforeEach(() => {
    dynamicActivate('en');
  });
  it('should render', () => {
    render(<CardHeader>Test</CardHeader>);

    expect(screen.getByText(/Test/i)).toBeInTheDocument();
  });
});
