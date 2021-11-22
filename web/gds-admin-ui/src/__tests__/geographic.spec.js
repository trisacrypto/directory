import { render } from '@testing-library/react';
import faker from 'faker';
import { AddressTypeHeaders } from '../constants';
import Geographic, { renderField, renderLines } from '../pages/app/details/BasicDetails/components/Geographic';

describe('defaultEndpointPrefix', () => {
    it('should render component', () => {
        render(<Geographic />);
    });

    describe('RenderLine', () => {
        const addresses = [
            {
                address_line: faker.random.objectElement([['215 Alynn Way', '', 'Queenstown, MD 21658']]),
                address_type: faker.random.objectElement(Object.keys(AddressTypeHeaders)),
                building_name: '',
                building_number: '23',
                country: 'US',
                country_sub_division: 'MA',
                department: faker.commerce.department(),
                district_name: '',
                floor: '',
                post_box: '',
                post_code: faker.address.zipCode(),
                room: '',
                street_name: faker.address.streetName(),
                sub_department: faker.commerce.department(),
                town_location_name: '',
                town_name: faker.address.cityName(),
            },
        ];

        it('should render address line header', () => {
            const addressType = AddressTypeHeaders[addresses[0].address_type];

            const { getByTestId } = render(<Geographic data={addresses} />);

            expect(getByTestId(/addressType/i).textContent).toBe(`${addressType} Address:`);
        });

        it('should render correctly address line', () => {
            const { getByTestId, container } = render(renderLines(addresses[0]));

            expect(container).toMatchInlineSnapshot(`
                <div>
                  <address
                    data-testid="addressLine"
                  >
                    <div>
                      215 Alynn Way
                       
                    </div>
                    
                    <div>
                      Queenstown, MD 21658
                       
                    </div>
                    <div>
                      US
                    </div>
                  </address>
                </div>
            `);
            expect(getByTestId(/addressLine/i).textContent).toBe('215 Alynn Way Queenstown, MD 21658 US');
        });
    });

    describe('RenderField', () => {
        const addresses = [
            {
                address_line: [],
                address_type: faker.random.objectElement(Object.keys(AddressTypeHeaders)),
                building_name: '',
                building_number: '23',
                country: 'US',
                country_sub_division: 'MA',
                department: faker.commerce.department(),
                district_name: '',
                floor: '',
                post_box: '',
                post_code: faker.address.zipCode(),
                room: '',
                street_name: faker.address.streetName(),
                sub_department: faker.commerce.department(),
                town_location_name: '',
                town_name: faker.address.cityName(),
            },
        ];

        it('should render', () => {
            const { container } = render(renderField(addresses[0]));

            expect(container).toMatchInlineSnapshot(`
                <div>
                  <address
                    data-testid="addressField"
                  >
                    Shoes
                     
                    <br />
                    Music
                     
                    <br />
                    23
                     
                    Frida Hollow
                    <br />
                    Rancho Cucamonga
                     
                    
                     
                    MA
                     
                    81030-8195
                      
                    <br />
                    US
                  </address>
                </div>
            `);
        });
    });
});
