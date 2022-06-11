import { render, screen } from 'utils/test-utils';
import AddressForm from '.';
import { faker } from '@faker-js/faker';
import { dynamicActivate } from 'utils/i18nLoaderHelper';

describe('<AddressForm />', () => {
  beforeEach(() => {
    dynamicActivate('en');
  });
  it('should render correctly', () => {
    const name = 'test';
    const rowIndex = faker.datatype.number();
    render(<AddressForm name={name} rowIndex={rowIndex} />);

    // address_line[0]
    const firstAddressLine = screen.getByTestId('address_line[0]') as HTMLInputElement;

    expect(firstAddressLine).toBeVisible();
    expect(firstAddressLine.name).toBe(`${name}[${rowIndex}].address_line[0]`);

    // address_line[1]
    const secondAddressLine = screen.getByTestId('address_line[1]') as HTMLInputElement;
    expect(secondAddressLine).toBeVisible();
    expect(secondAddressLine.name).toBe(`${name}[${rowIndex}].address_line[1]`);

    // address_line[2]
    const thirdAddressLine = screen.getByTestId('address_line[2]') as HTMLInputElement;
    expect(thirdAddressLine).toBeVisible();
    expect(thirdAddressLine.name).toBe(`${name}[${rowIndex}].address_line[2]`);
  });
});
