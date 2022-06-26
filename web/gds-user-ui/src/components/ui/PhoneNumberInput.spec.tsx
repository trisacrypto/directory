import { render } from 'utils/test-utils';
import PhoneNumberInput from './PhoneNumberInput';

describe('<PhoneNumberInput />', () => {
  it('should ', () => {
    const handleChange = jest.fn();
    render(<PhoneNumberInput controlId="" onChange={handleChange} />);
  });
});
